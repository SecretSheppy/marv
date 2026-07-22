# Mewt

[Mewt](https://github.com/trailofbits/mewt) is a multi-language mutation testing framework. It is supported by Marv 
versions `1.2.5+`.

# Getting Started With Mewt In Marv

1. To get started with Mewt in Marv, run the `marv init` command in the root directory of your project to create the 
required `.marv.yml` configuration file.

```terminaloutput
marv init -f mutant
```

2. Edit the `mewt` section of the configuration to point towards the generated `mewt.sqlite` results database.

```yaml
# Enable the mewt framework
mewt:
    # The relative path, from the directory the .marv.yml file is in, to the
    # generated mewt.sqlite database.
    sqlite-path: mewt/mewt.sqlite
```

3. Run the `marv` command in the root directory of your project and click the localhost link to open the Marv interface.

```terminaloutput
marv
```

## Mewt To Marv Status Conversions

| Mewt Status | Marv Status |
|:-----------:|:-----------:|
| `TestFail`  |  `KILLED`   |
| `Uncaught`  | `SURVIVED`  |
|  `Skipped`  |  `IGNORED`  |
|  `Timeout`  |  `TIMEOUT`  |
