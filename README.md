# Commander

[![Go Reference](https://pkg.go.dev/badge/github.com/BacchusJackson/commander.svg)](https://pkg.go.dev/github.com/BacchusJackson/commander)
[![Go Report Card](https://goreportcard.com/badge/github.com/BacchusJackson/commander)](https://goreportcard.com/report/github.com/BacchusJackson/commander)


![Commander Logo](assets/commander-logo-dark-theme-splash.png)

## Quick Start

Commander lets you template shell commands using [Go's Templating Library](https://pkg.go.dev/text/template) with helpful built-in functions.
Shell commands can be rendered with provided configured values in a single file.

```
cmdr init > commander.json
```
```json
{
  "docker-build": {
    "description": "build an image with docker",
    "template": "docker build {{- if .file}} --file {{.file}} {{end}} {{- if .context}} {{.context}} {{else}} . {{- end}}",
    "values": {
      "file": "Dockerfile.custom"
    }
  },
  "echo": {
    "description": "print a text message",
    "template": "echo {{- if .newline}} -n {{end -}} \"{{.msg}}\"",
    "values": {
      "msg": "howdy world",
      "newline": "true"
    }
  }
}
```

```shell
cmdr exec --file commander.json --target echo
# output: "howdy world"
```

### Support file formats

Commander supports configuration files in JSON, YAML or TOML using the `--encode` (`-e` shorthand) flag.

YAML configuration
```yaml
docker-build:
    description: build an image with docker
    template: docker build {{- if .file}} --file {{.file}} {{end}} {{- if .context}} {{.context}} {{else}} . {{- end}}
    values:
        file: Dockerfile.custom
echo:
    description: print a text message
    template: echo {{- if .newline}} -n {{end -}} "{{.msg}}"
    values:
        msg: howdy world
        newline: "true"
```

TOML configuration
```toml
[docker-build]
description = "build an image with docker"
template = "docker build {{- if .file}} --file {{.file}} {{end}} {{- if .context}} {{.context}} {{else}} . {{- end}}"
[docker-build.values]
file = "Dockerfile.custom"

[echo]
description = "print a text message"
template = "echo {{- if .newline}} -n {{end -}} \"{{.msg}}\""
[echo.values]
msg = "howdy world"
newline = "true"
```

Note: Environment variable expansion happens when the command is rendered, before execution.
