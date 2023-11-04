# exec

`cmdr exec` is used to execute rendered commands.

You can also use `-n` to do a dry run which will print the rendered command instead of actually executing it.
This is useful for debugging.

```
execute a command after applying templating

Usage:
  cmdr exec [flags] [args]

Examples:

        cmdr exec --file commander.json --target echo
        cmdr exec --file commander.yaml -e yaml --dry-run --target echo
        cat commander.yaml | cmdr exec --stdin -e yaml --target echo


Flags:
  -n, --dry-run         print rendered command but don't execute
  -e, --encode string   set text encoding for input can be [json,yaml,toml] (default "json")
  -f, --file string     the commander configuration file with templated commands and values (default "commander.json")
  -h, --help            help for exec
  -s, --stdin           read commander file from stdin
  -t, --target string   the command to run in the commander file

Global Flags:
  -v, --verbose   verbose logging output
```

