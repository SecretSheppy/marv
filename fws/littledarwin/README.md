# LittleDarwin

[LittleDarwin](https://github.com/aliparsai/LittleDarwin) is a mutation testing framework for Java written in Python.

## Getting Started With LittleDarwin In Marv

1. To get started with LittleDarwin in Marv, run the `marv init` command to generate the required `.marv.yml`
   configuration file.

```terminaloutput
marv init -f littledarwin
```

2. Edit the field under the `littledarwin` section in the `.marv.yml` file to point at the generated `LittleDarwinResults` directory.
   An example is shown below:

```yaml
# Enable the littledarwin framework
littledarwin:
  # The relative path to the LittleDarwinResults directory
  results-dir: LittleDarwinResults
```

3. Run the `marv` command to launch Marv and click the localhost URL to open the Marv interface.

```terminaloutput
marv
```
