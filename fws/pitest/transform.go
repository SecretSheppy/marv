package pitest

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/SecretSheppy/marv/decompilers"
	"github.com/SecretSheppy/marv/pkg/mutations"
	"github.com/aymanbagabas/go-udiff"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

// Groups mutations by class name so all mutations of the same file will are stored under the file path.
func (p *Pitest) groupMutants(bar *progressbar.ProgressBar) FileMutations {
	fileMutations := make(FileMutations)
	for _, m := range p.muts {
		fileMutations[m.SourceCodePath()] = append(fileMutations[m.SourceCodePath()], m)
		bar.Add(1)
	}
	return fileMutations
}

// Assigns a value to every Mutation's Mutation.MutationIndex. This is done because the order that Pitest exports the
// mutants in the mutations.xml file differs to that of the order it exports the mutant class files with the
// -Dfeatures="+EXPORT" flag. Pitest also exports a file for each mutant called details.txt which contains a string
// version of the MutantIdentifier class. Given a target file, to index each mutant: a portion of the MutationIdentifier
// string is reproduced and then compared with each exported details.txt file until a match is found.
func (p *Pitest) indexMutants(fileMutations FileMutations, bar *progressbar.ProgressBar) error {
	for _, mutations := range fileMutations {
		for _, mutant := range mutations {
			for i := 0; i < len(mutations); i++ {
				details := fmt.Sprintf("%d/details.txt", i)
				detailsPath := path.Join(p.yml.Cfg.MutClassPath, mutant.MutantExportDir(), details)
				mutationIdentifier, err := os.ReadFile(detailsPath)
				if err != nil {
					return err
				}

				values := fmt.Sprintf("lineNumber=%d, description=%s", mutant.LineNumber, mutant.Description)
				if strings.Contains(string(mutationIdentifier), values) {
					mutant.MutationIndex = i
					break
				}
			}
			bar.Add(1)
		}
	}
	return nil
}

// exists to ensure that the transform method can never modify any of the private Pitest variables.
type transformable interface {
	config() *YamlConfig
	decompiler() decompilers.Decompiler
}

type transformJob struct {
	Mutations []*Mutation
}

type transformResult struct {
	Mutations mutations.Mutations
}

type transformError struct {
	Err                       error
	Message, File, Decompiler string
}

func (e *transformError) Error() string {
	return e.Err.Error()
}

func (e *transformError) log() {
	log.Error().Err(e.Err).Str("file", e.File).Str("decompiler", e.Decompiler).Msgf("%s - %s", FWName, e.Message)
}

func newTransformError(err error, message, file, decompiler string) transformError {
	return transformError{err, message, file, decompiler}
}

// Manages dispatching and collecting workers that perform the transformation operations for each source file.
func transform(pit transformable, fileMutations FileMutations, bar *progressbar.ProgressBar) (mutations.Mutations, []error) {
	var (
		numWorkers = runtime.NumCPU()
		jobs       = make(chan transformJob, len(fileMutations))
		results    = make(chan transformResult, len(fileMutations))
		errs       = make(chan transformError)

		wg sync.WaitGroup
	)

	for range numWorkers {
		wg.Add(1)
		go transformWorker(jobs, results, errs, pit.config(), pit.decompiler(), &wg, bar)
	}

	for _, mutations := range fileMutations {
		jobs <- transformJob{Mutations: mutations}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	mutants := make(mutations.Mutations)
	accErrors := make([]error, 0)
	for err := range errs {
		accErrors = append(accErrors, &err)
	}
	for result := range results { // NOTE: inner loop only runs once per result
		for k, v := range result.Mutations {
			mutants[k] = v
		}
	}
	return mutants, accErrors
}

// Extracts all mutants for a given source file and returns them in the marv mutations.Mutations format. There are no
// errors that terminate the execution of this process, errors are sent through the channels when they occur and the
// process moves onto the next mutant.
func transformWorker(jobs <-chan transformJob, results chan<- transformResult, errs chan<- transformError, config *YamlConfig, decompiler decompilers.Decompiler, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	for job := range jobs {
		sourceCodePath := path.Join(config.SrcCodePath, job.Mutations[0].SourceCodePath())
		sourceCodeLines, err := readLines(sourceCodePath)
		if err != nil {
			errs <- newTransformError(err, "Failed to process mutations for file", sourceCodePath, decompiler.ExePath())
			bar.Add(1)
			continue
		}

		decompiled := make(map[string]string) // NOTE: Cache of decompiled class files where they can be reused.
		ms := make(mutations.Mutations)

		for _, mutant := range job.Mutations {
			// NOTE: Due to use of test builders in Java, some mutants of the same class may have different source
			// class files, so caching the source class cannot be done outside of this loop.
			sourceClassPath := path.Join(config.SrcClassPath, mutant.SourceClassPath())
			if decompiled[sourceClassPath] == "" {
				result, err := decompiler.Decompile(sourceClassPath)
				if err != nil {
					errs <- newTransformError(err, "Failed to decompile source class file", sourceClassPath, decompiler.ExePath())
					continue
				}
				decompiled[sourceClassPath] = string(result)
			}

			mutantClassPath := path.Join(config.MutClassPath, mutant.MutatedClassPath())
			mutatedClass, err := decompiler.Decompile(mutantClassPath)
			if err != nil {
				errs <- newTransformError(err, "Failed to decompile mutant class file", mutantClassPath, decompiler.ExePath())
				continue
			}

			edits := udiff.Strings(decompiled[sourceClassPath], string(mutatedClass))
			d, err := udiff.ToUnifiedDiff("old", "new", decompiled[sourceClassPath], edits, 0)
			if err != nil {
				errs <- newTransformError(err, "Unified diffing failure", mutantClassPath, decompiler.ExePath())
				continue
			}

			// TODO: this can still produce some strange results. could be good to implement some kind of "mutant
			//  smoothing" that can try and adjust for when lines don't end in ; or a mutant ends in a return statement
			removedLineCount := 0
			builder := strings.Builder{}
			for _, h := range d.Hunks {
				if strings.Contains(h.Lines[0].Content, "import") { // NOTE: ignore import diff hunks
					continue
				}
				for _, l := range h.Lines {
					switch l.Kind {
					case udiff.Delete:
						removedLineCount++
					case udiff.Insert, udiff.Equal:
						builder.WriteString(l.Content)
					}
				}
			}

			// NOTE: sl.LineNumber is technically the first line removed in removedLineCount, so hence the -2 here and
			// -1 in the below mutations.Range.
			endLineNumber := mutant.LineNumber + removedLineCount - 2
			if endLineNumber >= len(sourceCodeLines) {
				errs <- newTransformError(err, "Calculated mutant end line number exceeds source line count",
					mutantClassPath, decompiler.ExePath())
				continue
			}

			srcEndLine := sourceCodeLines[endLineNumber]
			sl := streamlineMutation(
				mutant,
				&mutations.Range{Line: mutant.LineNumber - 1},
				&mutations.Range{Line: endLineNumber, Char: len(srcEndLine) - 1})
			sl.Source = builder.String()
			ms.Append(mutant.SourceCodePath(), sl)
		}

		results <- transformResult{Mutations: ms}
		bar.Add(1)
	}
}

// Reads the contents of the specified file and returns its lines as an array of strings.
func readLines(file string) ([]string, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var lines []string
	for line := range strings.Lines(string(raw)) {
		lines = append(lines, line)
	}
	return lines, nil
}
