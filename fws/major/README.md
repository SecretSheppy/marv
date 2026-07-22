# Major

[Major](https://mutation-testing.org/) is a mutation testing framework for Java. It is supported by Marv versions 
`1.2.6+`.

## Getting Started With Major In Marv

1. To get started with Major in Marv, run the `marv init` command seen below to generate the required `.marv.yml`
configuration file.

```terminaloutput
marv init -f Major
```

2. Edit the fields under the `major` section of the `.marv.yml` file to reflect the layout of your project. The `major`
fields are described in the below YAML file.

```yaml
# Enable the major framework
major:
    # The path to the projects source directory.
    src-dir: java/major/src
    
    # The path to the directory where major outputs its results. In most use cases
    # this would simply be "."
    output-dir: java/major
```

3. Run the `marv` command to launch Marv and click the localhost URL to open the web interface.

```terminaloutput
marv
```