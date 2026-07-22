<img src="docs/marv_banner.png" alt="Marv" style="width: 100%;">

<h2 align="center">Marv: Mutations Analysis, Review and Visualisation</h2>

Marv is a data processing, visualization and review tool for mutation testing. Marv provides both a standardized
mutations format (through the Marv mutations schema) and visualization for results of mutation analysis across all
[supported frameworks](#supported-frameworks).

**What Makes Marv Great?:**
* High quality visualizations of source code mutations
* Calculates a range of metrics by which to evaluate the results
* Generates textual descriptions of mutations where none are provided
* Filter mutants by mutation status
* Work with results from multiple frameworks simultaneously
* Resolves issues with output from mutation testing tools
* Produces a standardized output (the Marv mutations schema)
* Can store and export textual reviews for each mutation
* Fast, responsive and intuitive user interface ([gallery](#gallery))
* Themable interface through JSON theme files
* Scales to hundreds of thousands of mutations without issue
* A large number of supported frameworks
* Framework support is built into Marv (tools don't need to add and maintain support for Marv)

## Table of Contents

* [Supported Frameworks](#supported-frameworks)
  * [LLM Based Mutation Testing Frameworks](#llm-based-mutation-testing-frameworks)
* [Installation](#installation)
  * [Build from source](#build-from-source)
  * [Libraries](#libraries)
* [Usage](#usage)
* [Gallery](#gallery)
* [Export Formats](#export-formats)
  * [Marv Mutations Schema](#marv-mutations-schema)
  * [Marv Reviews Schema](#marv-reviews-schema)
* [Other](#other)

## Supported Frameworks

The following table lists the mutation testing frameworks that are supported by Marv, and breaks down the quality of
support for each framework against certain key features.

* 🏆 Supported to a high standard
* 🥈 Supported to an acceptable standard
* ⚠️ Supported but experimental
* 🚧 In development
* 🚫 Not supported

| Framework                                                                       | Language | Source Code Replacements |  Descriptions  | Operators | Framework IDs | Marv Version | Documentation                                                                  |
|---------------------------------------------------------------------------------|----------|:------------------------:|:--------------:|:---------:|:-------------:|:------------:|--------------------------------------------------------------------------------|
| [Mewt](https://github.com/trailofbits/mewt)                                     | —        |            🏆            | 🥈<sup>1</sup> |    🏆     |      🏆       |    1.2.5     | [mewt](fws/mewt/README.md)                                                     |
| [Mull](https://mull-project.com/)                                               | C/C++    |            🏆            | 🏆<sup>1</sup> |    🏆     |      🏆       |    1.2.0     | [stryker-mte](fws/stryker_net/README.md)                                       |
| [stryker&#8288;&#45;&#8288;net](https://github.com/stryker-mutator/stryker-net) | C#       |            🏆            | 🏆<sup>2</sup> |    🏆     |      🏆       |    1.2.0     | [stryker-mte](fws/stryker_net/README.md)                                       |
| [go&#8288;&#45;&#8288;mutesting](https://github.com/zimmski/go-mutesting)       | Go       |            🥈            |       🚫       |    🏆     |      🚫       |    1.2.1     | [go-mutesting](fws/go_mutesting/README.md)                                     |
| [pitest](https://github.com/hcoles/pitest)                                      | Java     |            ⚠️            |       🏆       |    🏆     |      🚫       |    1.0.0     | [pitest](fws/pitest/README.md)                                                 |
| [Major](https://mutation-testing.org/)                                          | Java     |            🏆            | 🏆<sup>1</sup> |    🏆     |      🏆       |    1.2.6     | [major](fws/major/README.md)                                                   |
| [stryker&#8288;&#45;&#8288;js](https://github.com/stryker-mutator/stryker-js)   | JS/TS    |            🏆            | 🏆<sup>2</sup> |    🏆     |      🏆       |    1.2.0     | [stryker-mte](fws/stryker_net/README.md)                                       |
| [infection](https://github.com/infection/infection)                             | PHP      |            🥈            |       🏆       |    🏆     |      🏆       |    1.2.0     | [stryker-mte](fws/stryker_net/README.md), [infection](fws/infection/README.md) |
| [mutant](https://github.com/mbj/mutant)                                         | Ruby     |            🥈            | 🥈<sup>1</sup> |    🚫     |      🏆       |    1.2.6     | [mutant](fws/mutant/README.md)                                                 |
| [cargo&#8288;&#45;&#8288;mutants](https://github.com/sourcefrog/cargo-mutants)  | Rust     |            🏆            |       🏆       |    🏆     |      🚫       |    1.2.6     | [cargo-mutants](fws/cargo_mutants/README.md)                                   |
| [mutest&#8288;&#45;&#8288;rs](https://github.com/zalanlevai/mutest-rs)          | Rust     |            🏆            |       🏆       |    🏆     |      🏆       |    1.0.0     | [mutest-rs](fws/mutest_rs/README.md)                                           |
| [strkyer4s](https://github.com/stryker-mutator/stryker4s)                       | Scala    |            🏆            | 🏆<sup>2</sup> |    🏆     |      🏆       |    1.2.0     | [stryker-mte](fws/stryker_net/README.md)                                       |

**Notes:**
1. Generated by Marv
2. If provided by framework

### LLM Based Mutation Testing Frameworks

Marv does not provide first party support for LLM based mutation testing frameworks. Instead, for LLM based tools to
be compatible with Marv, they must output in the [Marv mutations schema](api/marv-mutations-schema.json) which is supported through the use of the
`generic` framework implementation.

## Installation

Marv can be installed with the `go` tool. To use the installed executable, add the `GOPATH` environment variable 
to your system path. For more information run `go help install`.
```
go install github.com/SecretSheppy/marv@latest
```

### Build from source

Builds exactly as a normal go project would. See the [go.dev](https://go.dev/doc/tutorial/compile-install) tutorial
for more information. The target file to build is [`cmd/marv/main.go`](cmd/marv/main.go).
```cli
go build main.go -o marv
```

### Libraries

If using a framework that requires external libraries, they will need to be set with the `MARV_LIB_PATH` environment
variable. The alternative to this is to put the library into the local directory where the Marv tool is being run.

## Usage

A simple guide of how to run Marv on a project for the first time. If at any point you need more information about
one of the Marv commands, try using the help command.
```cli
marv help [command]
```

1. The first step is to ensure that Marv is correctly installed. If the Marv executable is correctly installed, running
the Marv version command will output a version number. If an error is printed, then it likely means you need to add
the Marv executable install location to your system path.
```cli
marv --version
```

2. Run the list command to see a `list` of all the frameworks that your installed version of Marv supports.
```cli
marv list
```

3. Then navigate to your project location and run the Marv `init` command with the list of frameworks you are using
(Marv framework names are case-sensitive, so make sure to copy them correctly from the output of the `list` command).
This will create a `.marv.yml` file in the directory that Marv was run in. The file will contain the default Marv
configuration as well as a blank configuration for each framework that was listed.
```cli
marv init -f [framework] -f [framework] ...
```

4. Now fill in the configurations for each framework. Where frameworks require paths, using paths relative to a repository
will allow you to safely commit the `.marv.yml` file for others to use. When finished with the configurations, simply
run `marv`.
```cli
marv
```

5. If you have correctly configured the frameworks then that is it! Provided you keep the `.marv.yml` configuration
file then all you have to do in future is simply run `marv`.

## Gallery

Screenshots of the Marv user interface showing results from various frameworks.

|                                                                                                                                                                                                                          |                                                                                                                                                                                                   |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Marv Results Overview:** Showing results from [stryker-net](https://github.com/stryker-mutator/stryker-net) run on itself<br/> ![](docs/marv_results_overview.png)                                                     | **Marv Pitest Results:** Showing [hcoles/pitest](https://github.com/hcoles/pitest) mutants inline with a file from [guava](https://github.com/google/guava)<br/> ![](docs/marv_pitest_guava.png)  |
| **Marv mutest-rs Results:** Showing [mutest-rs](https://github.com/zalanlevai/mutest-rs) mutants inline with a file from [alacritty](https://github.com/alacritty/alacritty)<br/> ![](docs/marv_mutest_rs_alacritty.png) | **Marv Infection PHP Mutant:** Showing an isolated [Infection](https://github.com/infection/infection) mutant inline with a file from its own source<br/> ![](docs/marv_infection_php_mutant.png) |

## Export Formats

Marv exports both the mutations and reviews as a `.json` marshal of its internal mutations format for all frameworks. 
By using the `-m` or `--merge` flags, the results from all frameworks are merged into one large `.json` file. Brief
textual descriptions of both formats are given below, and schemas for each can be found in the [`api`](api) directory:

* [Mutations Format Schema](api/marv-mutations-schema.json)
* [Reviews Format Schema](api/marv-reviews-schema.json)

### Marv Mutations Schema

The mutations schema follows the internal structures defined in [`internal/mutations`](internal/mutations/mutations.go). 
The basic structure is `file path` > `conflict region` > `mutation`. Marv uses `conflict regions` or internally called 
`mutations.Conflict` to wrap all mutations that would conflict with each other if just rendered inline due to overlaps.

Any `ID` field is a UUID created by Marv. Where frameworks create mutant identifiers, they are stored against the mutant
as `FrameworkMutantID`.

### Marv Reviews Schema

Reviews are exported against the corresponding mutations Marv Mutation ID and their Framework Mutation 
ID (if applicable). The review structure is defined in [`internal/review`](internal/review/review.go).

## Other

* [Icon Licenses (icons.md)](icons.md)
* [Contributing to Marv](CONTRIBUTING.md)
