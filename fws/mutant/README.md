# Mutant

[Mutant](https://github.com/mbj/mutant) is a mutation testing framework for Ruby. It is supported by Marv versions `1.2.6+`.

## Getting started

To get started with mutant in Marv, simply run the following command to generate the
required `.marv.yml` file in your current working directory. It is usually better to do this inside the project you
are working on, but this is not an absolute requirement.

```terminaloutput
marv init -f mutant
```

When the file has been created, you will need to edit the fields under the `mutant` section to match the layout of the
project you are working on. The `mutant` fields are described in the below YAML file.

```yaml
# The default Marv configuration options. These can be changed at will.
marv:
    port: 8080
    output:
        path: .marv
        merge: false

# Enable the mutant framework
mutant:
    # The projects root directory name (usually name of the repository). This is
    # used by Marv to convert the absolute paths exported by mutant into relative
    # paths. See "Locating the source files" for more information.
    root-dir: project
    
    # The relative path to the results directory created by mutant.
    results-dir: results
    
    # An optional field to fix which JSON is loaded by Marv. By default, Marv will
    # scan the contents of the mutant results directory and load the most recently
    # created JSON which has a UUID name. See "Locating the correct JSON output"
    # for more information.
    session: 00000000-0000-0000-0000-000000000000
```

Once this file has been edited and saved, simply run the `marv` command in the directory the file was created and Marv
will process and visualize the results for you.

```terminaloutput
marv
```

If something goes wrong at this stage it is likely that one of the provided paths is slightly incorrect. Reading the
stderr output can most often tell you where you have gone wrong.

## Locating the correct JSON output

The JSON data exported by mutant does not have a standard name. Instead, names each file with the session UUID that it
creates. Marv's default behavior is to read the entirety of the specified `results` directory and use the most recently
created JSON file with a UUID name as the results to process and display. One can, however, direct Marv towards a
specific session ID, which will ensure that Marv always opens the same JSON results file.

## Locating the source files

Mutant exports absolute paths for its source code files. This causes issues when running Marv on the mutant results
from a machine that was not responsible for running mutant.

The mutant configuration for Marv requires the location of the files under test in relation to the working directory
that Marv has been run from. Marv uses the name of the local source files root directory to try and strip out the
absolute paths exported by mutant. An example is provided below for clarification.

If mutant exports the path `/home/user1/project/src/lib/file.rb` and you are running Marv in the
`/home/user2/project` directory, one simply has to set `root-dir: project` and Marv will then strip
`/home/user1/project/` out of all the paths provided by mutant, leaving only relative paths in the Marv output.