# init

The `init` command is used to generate a starting example commander file.
It will default to JSON encoding unless another is specified with the `--encode / -e` flag.

```
Print an example configuration file to stdout

Usage:
  cmdr init [flags] [args]

Flags:
  -e, --encode string   set text encoding for output can be [json,yaml,toml] (default "json")
  -h, --help            help for init

Global Flags:
  -v, --verbose   verbose logging output
```

Example of printing a configuration file for each supported format type:
```
cmdr init
cmdr init -e yaml
cmdr init -e toml
```


This is a useful starting place for building more complex commands.
