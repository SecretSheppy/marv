# Cosmic Ray

[Cosmic Ray](https://github.com/sixty-north/cosmic-ray) is a mutation testing framework for Python. It is supported by
all Marv versions `1.2.7+`.

## Contents

* [Getting Started With Cosmic Ray In Marv](#getting-started-with-cosmic-ray-in-marv)
* [Legacy Documentation](#legacy-documentation)

## Getting Started With Cosmic Ray In Marv

1. To get started with Cosmic Ray in Marv, run the `marv init` command to generate the required `.marv.yml`
   configuration file.

```terminaloutput
marv init -f cosmic-ray
```

2. Edit the fields under the `cosmic-ray` section in the `.marv.yml` file to point at the relevant locations.
   An example is shown below:

```yaml
# Enable the cosmic-ray framework
cosmic-ray:
    # The relative path the cosmic ray session sqlite database.
    sqlite-path: session-name.sqlite
    
    # The relative path to the working directory where cosmic ray was run.
    cr-work-dir: .
```

3. Run the `marv` command to launch Marv and click the localhost URL to open the Marv interface.

```terminaloutput
marv
```

## Legacy Documentation

The Cosmic Ray framework implementation in Marv used to produce erroneous data under certain circumstances. This is
irrelevant if using Marv v1.3.0+, but error details for users of Marv v1.2.11 or lower can be seen [here](legacy_v1.2.11_before.md).