# cargo-mutants

This document outlines important details one should be aware of when using cargo-mutants with Marv.

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