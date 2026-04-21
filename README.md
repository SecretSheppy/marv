<p align="center">
  <img src="web/static/branding/marv_logo.png" alt="The Marv logo: Stylistic text that reads 'Marv' over a red rectangle containing a '-' and a green rectangle containing a '+'" style="width: 150px;">
</p>

<h2 align="center">Mutations Analysis, Review and Visualisation</h2>

Marv is a visualization and review tool for mutation testing.

Marv allows for the results of multiple frameworks to be displayed simultaneously, creating the opportunity to review
results across many frameworks or even languages in one go.

### About Marv

Marv is closely related to the [mutest-rs](https://github.com/zalanlevai/mutest-rs) inspector. Both Marv and 
mutest-inspector started off as the same project, the initial prototype reporter for mutest-rs, but then moved in
different directions after the prototype was completed.

## Supported Frameworks

A list of mutation testing frameworks that either are currently supported or will be supported in the future.

* 🏆 Supported out of the box
* ✅️ Supported with some configuration
* ⚠️ Experimental support
* 🚧 In development
* 🚫 Not currently supported

| Framework                                                                                                                  | language   | Support | Marv Version | Required Libraries                                                                                                                                                                                             | Notes                                                   |
|----------------------------------------------------------------------------------------------------------------------------|------------|:-------:|:------------:|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------|
| [Mull](https://mull-project.com/)                                                                                          | C/C++      |   🚧    |              |                                                                                                                                                                                                                |                                                         |
| [Dextool Mutate](https://github.com/joakim-brannstrom/dextool/tree/32bdcad647838121160862b61c852f7a8c26f35c/plugin/mutate) | C/C++      |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [stryker-net](https://github.com/stryker-mutator/stryker-net)                                                              | C#         |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [hcoles/pitest](https://github.com/hcoles/pitest)                                                                          | Java       |   ✅️    |    1.0.0     | <ul><li>[vineflower-server](https://github.com/SecretSheppy/vineflower-server)</li><li>[vineflower](https://github.com/Vineflower/vineflower)</li><li>[garlic](https://github.com/neocanable/garlic)</li></ul> | See [Pitest configuration](#pitest)                     |
| [Major](https://mutation-testing.org/)                                                                                     | Java       |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [stryker-js](https://github.com/stryker-mutator/stryker-js)                                                                | JavaScript |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [infection](https://github.com/infection/infection)                                                                        | PHP        |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [Cosmic Ray](https://github.com/sixty-north/cosmic-ray)                                                                    | Python     |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [MutPy](https://github.com/mutpy/mutpy)                                                                                    | Python     |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [mutant](https://github.com/mbj/mutant)                                                                                    | Ruby       |   🚫    |              |                                                                                                                                                                                                                |                                                         |
| [mutest-rs](https://github.com/zalanlevai/mutest-rs)                                                                       | Rust       |   🏆    |    1.0.0     | Native                                                                                                                                                                                                         |                                                         |

### Pitest Configuration

Pitest must be run with the `-Dfeatures="+EXPORT"` flag. This exports the mutants as class files that Marv can then
decompile and extract the mutants from. This process is not 100% reliable, and Marv will occasionally make mistakes,
however the mutants that Marv extracts are usually correct.

#### Decompilers

> [!CAUTION]
> The `garlic` decompiler is currently unstable and using it could cause some mutants to be skipped due to a segmentation fault that occurs when running `garlic` on some class files.

## Gallery

Screenshots of the Marv user interface showing results from:

* [mutest-rs](https://github.com/zalanlevai/mutest-rs) run on [alacritty](https://github.com/alacritty/alacritty)
* [hcoles/pitest](https://github.com/hcoles/pitest) run on [guava](https://github.com/google/guava)

|                                                                                                                                       |                                                                                                                                      |
|---------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| **Marv Results Overview:** Showing results from `mutest-rs` and `Pitest`<br/> ![](docs/marv_results_overview.png)                     | **Marv Pitest Results:** Showing `Pitest` mutants inline with a file from guava<br/>![](docs/marv_pitest_guava.png)                  |
| **Marv mutest-rs Results:** Showing `mutest-rs` mutants inline with a file from alacritty<br/> ![](docs/marv_mutest_rs_alacritty.png) | **Marv Pitest Mutant:** Showing an isolated `Pitest` mutant inline with a file from guava<br/>![](docs/marv_pitest_guava_mutant.png) |
