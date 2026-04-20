<p align="center">
  <img src="web/static/branding/marv_logo.png" style="width: 150px;">
</p>

<p align="center">Mutations Analysis, Review and Visualisation</p>

## About Marv

Marv is a visualization and review tool for mutation testing. Marv displays the same quality visualizations for every
framework that it supports with the ultimate goal being to support as many mutation testing frameworks as possible.

Marv allows for the results of multiple frameworks to be displayed simultaneously, creating the opportunity to review
results across many frameworks or even languages in one go.

![marv ui showing a review of mutest-rs mutations on alacritty](docs/marv_mutest_rs_alacritty.png)

**Caption:** The Marv user interface showing mutants generated my the mutest-rs framework for rust on the Alacritty codebase.

## Supported Frameworks

> [!NOTE]
> Marv does not currently support many mutation testing frameworks. The frameworks that will eventually be supported are listed below, hence the large number of currently not supported frameworks.

* 🏆 Supported out of the box
* ✅️ Supported with some configuration
* ⚠️ Experimental support
* 🚧 In development
* 🚫 Not currently supported
* ❌ Totally incompatible (can never be supported in current state)

| Framework                                                     | language   | Support | Marv Version      | Required Libraries                                                                                                                                                                                             | Notes                                        |
|---------------------------------------------------------------|------------|:-------:|-------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------|
| [stryker-net](https://github.com/stryker-mutator/stryker-net) | C#         |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [crytic](https://github.com/hanneskaeufler/crytic)            | Crystal    |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [go-mutesting](https://github.com/zimmski/go-mutesting)       | Go         |   🚧    |                   |                                                                                                                                                                                                                | https://github.com/SecretSheppy/marv/pull/25 |
| [ooze](https://github.com/gtramontina/ooze)                   | Go         |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [germlins](https://github.com/go-gremlins/gremlins)           | Go         |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [MuCheck](https://github.com/vrthra/mucheck)                  | Haskell    |    ❌    |                   |                                                                                                                                                                                                                |                                              |
| [FitSpec](https://github.com/rudymatela/fitspec)              | Haskell    |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [hcoles/pitest](https://github.com/hcoles/pitest)             | Java       |   ✅️    | pre&#8209;release | <ul><li>[vineflower-server](https://github.com/SecretSheppy/vineflower-server)</li><li>[vineflower](https://github.com/Vineflower/vineflower)</li><li>[garlic](https://github.com/neocanable/garlic)</li></ul> | See [Pitest configuration](#pitest)          |
| [LittleDarwin](https://github.com/aliparsai/LittleDarwin)     | Java       |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [Major](https://mutation-testing.org/)                        | Java       |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [stryker-js](https://github.com/stryker-mutator/stryker-js)   | JavaScript |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [infection](https://github.com/infection/infection)           | PHP        |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [PEST](https://pestphp.com/docs/mutation-testing)             | PHP        |    ❌    |                   |                                                                                                                                                                                                                |                                              |
| [Cosmic Ray](https://github.com/sixty-north/cosmic-ray)       | Python     |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [MutPy](https://github.com/mutpy/mutpy)                       | Python     |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [mutmut](https://github.com/boxed/mutmut)                     | Python     |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [mutant](https://github.com/mbj/mutant)                       | Ruby       |   🚫    |                   |                                                                                                                                                                                                                |                                              |
| [mutest-rs](https://github.com/zalanlevai/mutest-rs)          | Rust       |   🏆    | pre&#8209;release | Native                                                                                                                                                                                                         |                                              |
| [muter](https://github.com/muter-mutation-testing/muter)      | Swift      |    ❌    |                   |                                                                                                                                                                                                                |                                              |

### Pitest Configuration

Pitest must be run with the `-Dfeatures="+EXPORT"` flag. This exports the mutants as class files that Marv can then
decompile and extract the mutants from. This process is not 100% reliable, and Marv will occasionally make mistakes,
however the mutants that Marv extracts are usually correct.

> [!CAUTION]
> The `garlic` decompiler is currently unstable and using it could cause some mutants to be skipped due to a segmentation fault that occurs when running `garlic` on some class files.
