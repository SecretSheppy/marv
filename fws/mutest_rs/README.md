# mutest-rs

[mutest-rs](https://github.com/zalanlevai/mutest-rs) is a mutation testing framework for Rust. It is supported by all 
Marv versions `1.0.0+`.

## Getting Started With mutest-rs In Marv

1. To get started with mutes-rs in Marv, run the `marv init` command to generate the required `.marv.yml` configuration
file.

```terminaloutput
marv init -f mutest-rs
```

2. Edit the field under the `mutest-rs` section in the `.marv.yml` file to point at the relevant locations.
   The fields are explained with comments in the below YAML file.

```yaml
# Enable the mutest-rs framework
mutest-rs:
    # The relative path to the source files.
    src: src
    
    # The mutest-rs JSON output directory.
    json-dir: project/json
```

3. Run the `marv` command to launch Marv and click the localhost URL to open the Marv interface.

```terminaloutput
marv
```