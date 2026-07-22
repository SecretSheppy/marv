# cargo-mutants

[cargo-mutants](https://github.com/sourcefrog/cargo-mutants) is a mutation testing framework for rust. It is supported
by all Marv versions `1.2.6+`.

## Contents

* [Getting Started With cargo-mutants In Marv](#getting-started-with-cargo-mutants-in-marv)
* [JSON parsing](#json-parsing)
* [Status Conversions](#status-conversions)

## Getting Started With cargo-mutants In Marv

1. To get started with cargo-mutants in Marv, run the `marv init` command to generate the required `.marv.yml`
   configuration file.

```terminaloutput
marv init -f cargo-mutants
```

2. Edit the fields under the `cargo-mutants` section in the `.marv.yml` file to point at the relevant locations.
   An example is shown below:

```yaml
# Enable the cargo-mutants framework
cargo-mutants:
    # The relative path to the working directory that the tests are run from.
    test-work-dir: .
    
    # The relative path to the mutants.out directory created by cargo-mutants.
    mutants-out-dir: mutants.out
```

3. Run the `marv` command to launch Marv and click the localhost URL to open the Marv interface.

```terminaloutput
marv
```

## JSON parsing

cargo-mutants exports all of its output in the `mutants.out` directory which is created automatically in the directory
where the tool was run, unless specified otherwise. Marv only reads the `outcomes.json` file from within this directory;
however, not all objects contained within this JSON share a consistent schema. Namely, the first `outcome` has a string
`scenario` field whereas all other `outcome` objects have an object `scenario` field. To prevent this from causing
a failure, Marv does a direct string replacement of `"Baseline"`, the string value of the first `scenario`, with
`null`. This allows Marv to parse the results without issue and should not cause issues, but if one is encountering 
issues with the output of cargo-mutants, this inconsistent JSON formatting is a good place to look.

## Status Conversions

In most frameworks the conversions from the native status to the Marv status are fairly obvious. cargo-mutants has a
very different status naming convention to most other frameworks, and so we explicitly detail the conversions below.

| cargo-mutants status | Marv status |
|:--------------------:|:-----------:|
|    `CaughtMutant`    |  `KILLED`   |
|    `MissedMutant`    | `SURVIVED`  |
|      `Timeout`       |  `TIMEOUT`  |
|      `Unviable`      |  `CRASHED`  |