package mull

import (
	"fmt"
	"regexp"
	"slices"

	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/internal/languages"
	"github.com/SecretSheppy/marv/internal/mtelib"
	"github.com/SecretSheppy/marv/internal/mutations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

var (
	meta = fwlib.Meta{
		Name:     "Mull",
		Language: languages.Cpp,
		URL:      "https://mull-project.com/",
	}
	function = regexp.MustCompile(`\s*([A-Za-z_0-9]*)\s*\(`)

	// cpp keywords according to https://en.cppreference.com/cpp/keyword
	cppKeyWords = []string{
		"alignas", "alignof", "and", "and_eq", "asm", "atomic_cancel", "atomic_commit", "atomic_noexcept", "auto",
		"bitand", "bitor", "bool", "break", "case", "catch", "char", "char8_t", "char16_t", "char32_t", "class",
		"compl", "concept", "const", "consteval", "constexpr", "constinit", "const_cast", "continue", "contract_assert",
		"co_await", "co_return", "co_yield", "decltype", "default", "delete", "do", "double", "dynamic_cast", "else",
		"enum", "explicit", "export", "extern", "false", "float", "for", "friend", "goto", "if", "inline", "int",
		"long", "mutable", "namespace", "new", "noexcept", "not", "not_eq", "nullptr", "operator", "or", "or_eq",
		"private", "protected", "public", "reflexpr", "register", "reinterpret_cast", "requires", "return", "short",
		"signed", "sizeof", "static", "static_assert", "static_cast", "struct", "switch", "synchronized", "template",
		"this", "thread_local", "throw", "true", "try", "typedef", "typeid", "typename", "union", "unsigned", "using",
		"virtual", "void", "volatile", "wchar_t", "while", "xor", "xor_eq",
	}
)

type YamlConfig struct {
	MTEJson string `yaml:"mte-json"`
}

type YamlWrapper struct {
	Cfg *YamlConfig `yaml:"mull"`
}

func (y *YamlWrapper) Init() interface{} {
	return &YamlWrapper{Cfg: &YamlConfig{}}
}

func (y *YamlWrapper) Load(yml []byte) (bool, error) {
	if err := yaml.Unmarshal(yml, y); err != nil {
		return false, err
	}
	if y.Cfg == nil {
		return false, nil
	}
	return y.Cfg.MTEJson != "", nil
}

type Mull struct {
	yml *YamlWrapper
	mte *mtelib.MTE
}

func NewMull() *Mull {
	return &Mull{yml: &YamlWrapper{}}
}

func (m *Mull) Meta() *fwlib.Meta {
	return &meta
}

func (m *Mull) Yaml() fwlib.FWConfig {
	return m.yml
}

func (m *Mull) LoadResults() error {
	log.Info().Msgf("%s - loading results", m.Meta().Name)
	var err error
	m.mte, err = mtelib.NewMTE(m.yml.Cfg.MTEJson)
	return err
}

func (m *Mull) TransformResults() error {
	log.Info().Msgf("%s - transforming results", m.Meta().Name)

	bar := fwlib.NewProgressbar(m.mte.RawMutationsCount(), "transforming")
	m.mte.Transform(bar)
	fwlib.FinishProgressbar(bar)

	fixed := 0
	for file, conflicts := range m.mte.Mutations() {
		lines := m.mte.ReadLines(file)
		for _, conflict := range conflicts {
			for _, mutation := range conflict.Mutations {
				if mutation.IsBroken() {
					fixed += attemptBrokenMutationFix(mutation)
					conflict.ResizeToInclude(mutation)
				}
				m.generateDescription(lines, mutation)
			}
		}
	}
	log.Info().Msgf("%s - fixed %d broken mutations", m.Meta().Name, fixed)

	return nil
}

func (m *Mull) generateDescription(lines []string, mutation *mutations.Mutation) {
	switch mutation.Operation {
	case "cxx_remove_void_call":
		mutation.Description = fmt.Sprintf("Removed call to void function `%s`", getFuncName(lines, mutation))
	case "cxx_replace_scalar_call":
		mutation.Description = fmt.Sprintf("Replaced call to function `%s` with `42`", getFuncName(lines, mutation))
	case "negate_mutator":
		mutation.Description = "Negated conditionals"
	case "scalar_value_mutator":
		mutation.Description = "Replaced zeros with `42` and non-zeros with `0`"
	default:
		line := lines[mutation.Start.Line]
		endChar := len(line) - 1
		if mutation.Start.Line == mutation.End.Line {
			endChar = mutation.End.Char
		}
		original := line[mutation.Start.Char:endChar]
		if mutation.Replacement != "" {
			mutation.Description = fmt.Sprintf("Replaced `%s` with `%s`", original, mutation.Replacement)
		} else {
			mutation.Description = fmt.Sprintf("Removed `%s`", original)
		}
	}
}

func attemptBrokenMutationFix(mutation *mutations.Mutation) int {
	switch mutation.Operation {
	case "cxx_assign_const", "cxx_init_const", "cxx_remove_void_call", "cxx_replace_scalar_void", "negate_mutator":
		// Operators marv does not currently fix.
		return 0
	case "cxx_bitwise_not_to_noop", "cxx_minus_to_noop", "cxx_post_dec_to_post_inc", "cxx_pre_dec_to_pre_inc", "cxx_remove_negation", "cxx_gt_to_ge", "cxx_gt_to_le", "cxx_lt_to_ge", "cxx_lt_to_le":
		// Operators that replace a source string length of 1 (~, !, -)
		mutation.End.Line = mutation.Start.Line
		mutation.End.Char = mutation.Start.Char + 1
	case "cxx_ge_to_gt", "cxx_ge_to_lt", "cxx_le_to_gt", "cxx_le_to_lt", "cxx_post_inc_to_post_dec", "cxx_pre_inc_to_pre_dec":
		// Operators that replace a source string length of 2 (==, <=, /=, ...)
		mutation.End.Line = mutation.Start.Line
		mutation.End.Char = mutation.Start.Char + 2
	default:
		// Operators that replace a source string of length equal to its replacement string
		mutation.End.Line = mutation.Start.Line
		mutation.End.Char = mutation.Start.Char + len(mutation.Replacement)
	}
	return 1
}

func getFuncName(lines []string, mutation *mutations.Mutation) string {
	match := function.FindAllStringSubmatch(lines[mutation.Start.Line], -1)
	funcStr := "??"
	for _, str := range match {
		// NOTE: takes first non keyword in the replacements string that matches the regex as the function name.
		if !slices.Contains(cppKeyWords, str[1]) {
			funcStr = str[1]
			break
		}
	}
	return funcStr
}

func (m *Mull) Mutations() mutations.Mutations {
	return m.mte.Mutations()
}

func (m *Mull) ReadLines(file string) ([]string, error) {
	return m.mte.ReadLines(file), nil
}
