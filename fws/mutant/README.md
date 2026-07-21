# Mutant

[Mutant](https://github.com/mbj/mutant) is a mutation testing framework for Ruby. It is supported by Marv versions `1.2.6+`.

## Contents

* [Getting Started With Mutant In Marv](#getting-started-with-mutant-in-marv)
* [Locating The Source Files For Marv](#locating-the-source-files-for-marv)
  * [An example `root-dir` value](#an-example-root-dir-value)
    * [1. Inside The `/home/fred/mutant/quick_start` Directory](#1-inside-the-homefredmutantquick_start-directory)
    * [2. Inside The `/home/fred/mutant` Directory](#2-inside-the-homefredmutant-directory)
* [Loading Results From Mutant Into Marv](#loading-results-from-mutant-into-marv)

## Getting Started With Mutant In Marv

1. To get started with mutant in Marv, simply run the following command to generate the
required `.marv.yml` file in your current working directory. It is usually better to do this inside the project you
are working on as a `.marv.yml` file can be committed and shared if correctly configured, but this is not an absolute
requirement.

```terminaloutput
marv init -f mutant
```

2. When the file has been created, you will need to edit the fields under the `mutant` section to match the layout of the
project you are working on. The `mutant` fields are described in the below YAML file.

```yaml
# Enable the mutant framework
mutant:
    # The path to the projects root directory. This is used by Marv to convert the
    # absolute paths exported by mutant into relative paths.
    root-dir: .
    
    # The relative path to the results directory created by mutant.
    results-dir: .mutant/results
    
    # An optional field to fix which JSON is loaded by Marv. By default, Marv will
    # scan the contents of the mutant results directory and load the most recently
    # created JSON which has a UUID name.
    session: 00000000-0000-0000-0000-000000000000
```

> [!TIP]
> For more information on what the values of `root-dir` and `results-dir` should be in your configuration, see 
> [Locating The Source Files For Marv](#locating-the-source-files-for-marv) and 
> [Loading Results From Mutant Into Marv](#loading-results-from-mutant-into-marv).

3. Once this file has been edited and saved, simply run the `marv` command in the directory the file was created and Marv
will process and visualize the results for you.

```terminaloutput
marv
```

> [!TIP]
> If something goes wrong at this stage it is likely that one of the provided paths is slightly incorrect. Reading the
stderr output can most often tell you where you have gone wrong.

## Locating The Source Files For Marv

Mutant uses absolute paths to reference the source files where mutations were made; however, Marv requires relative
paths for both its internal operations, and to correctly fulfil the Marv mutations schema. To rectify this, Marv uses
the name of the last directory in the path that results from joining the working directory Marv was run in with
the value specified in the `root-dir` field of the Mutant configuration in the `.marv.yml` field.

> [!TIP]
> Whilst an example of how to set the value of `root-dir` is detailed below, it may also help to review
> the [source code (mutant.go)](mutant.go:244) of the path conversion process.

### An example `root-dir` value

This example will follow the `quick_start` example included within Mutant's own source code.

We will start with some prerequisites:
1. We are a user on a Linux system known as `fred`
2. We have cloned the mutant repository into `/home/fred` and so creating the directory `/home/fred/mutant`

We now follow the [instructions to run the Mutant example](https://github.com/mbj/mutant/blob/main/quick_start/README.md).
Assuming that everything has worked correctly, a `.mutant` directory will have been created in 
`/home/fred/mutant/quick_start`. We now have two options for where we can create our `.marv.yml` configuration file.

#### 1. Inside The `/home/fred/mutant/quick_start` Directory

Usually it is best practise to create the `.marv.yml` file in the root directory of the repository 
(see [Option 2](#2-inside-the-homefredmutant-directory)); however, in this case `quick_start` is essentially the root
directory of the demonstration project. To configure Marv correctly here we would run ```marv init -f mutant``` and
then set the following values in the `.marv.yml` file:

```yaml
mutant:
    root-dir: .
    results-dir: .mutant/results
```

With this configuration, when `marv` is run, Marv joins the working directory (`/home/fred/mutant/quick_start`) with
our value for `root-dir` (`.`). From this it identifies that `quick_start` is the point in the absolute paths at which
they should be truncated, so every path like `/home/fred/mutant/quick_start/lib/person.rb` becomes `lib/person.rb` which
is correctly relative to the working directory.

#### 2. Inside The `/home/fred/mutant` Directory

The other option is to observe the best practises and run ```marv init -f mutant``` in the `/home/fred/mutant`
directory (the repository root directory). To configure Marv correctly here we would set the following values in the
`.marv.yml` file:

```yaml
mutant:
    root-dir: .
    results-dir: quick_start/.mutant/results
```

In this case, Marv does the same as above, but all the source paths would become relative to `/home/fred/mutant`: i.e.,
`/home/fred/mutant/quick_start/lib/person.rb` becomes `quick_start/lib/person.rb` which is correctly relative to the
working directory (and repository root directory).

## Loading Results From Mutant Into Marv

Mutant exports results as a JSON file in the `.mutant/results` directory. Each file is named with that particular
sessions UUID. Marv's default behavior is to scan the `results` directory, specified in the `.marv.yml` file, for all
JSON files with
UUID names and to display the results from the most recently created JSON. Marv accepts an optional `session` field in
the `.marv.yml` configuration for mutant. Setting `session` will tell Marv to always load the file named with the 
provided UUID. If no file with the provided name exists no results will be loaded.
