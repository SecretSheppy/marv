# The Generic Framework

The generic framework is a framework implementation in Marv that will load and visualize any data from the specified
JSON provided it complies with the Marv mutations schema.

## Getting Started With The Generic Framework

1. To get started with cargo-mutants in Marv, run the `marv init` command to generate the required `.marv.yml`
   configuration file.

```terminaloutput
marv init -f generic
```

2. Edit the fields under the `generic` section in the `.marv.yml` file to point at the relevant locations.
   An example is shown below:

```yaml
# Enable the generic framework
generic:
    # Sets the display name of the generic framework inside the Marv interface.
    framework: my-fw-name

    # The relative path to the Marv mutations schema JSON to be shown.
    marv-json: marv.json

    # The relative path to the source directory.
    src-dir: src
```

3. Run the `marv` command to launch Marv and click the localhost URL to open the Marv interface.

```terminaloutput
marv
```