# Getting Started

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

Commands have three parts: a name, description, and template.

_Name_: A custom identifier for the command that is used in `exec` as the `--target` since a file can have multiple 
commands.

_Description_: Optional, a human-readable explanation for what the command does.

_Template_: The command with template actions.

Each command has a section for values which is what gets used in the command during execution.
Inside the template, values are prefixed with a dot while in the values array, there is no dot.

## Building a Command

The best way to build a command is to do so incrementally.
Start with the command as you would execute it locally and pick out the parts that could be template'd.

Build this out into a Commander command, here we'll use TOML for readability 
but note that YAML and JSON are also supported.
TOML has support for multi-line strings by using triple quotes

```toml
[count-chars]
description = "Count the number of characters in a string"
template = """
echo "Some Content" | wc -c
"""
```
Using the `--dry-run` or `-n` flag will print the command after executing the template.

```shell
cmdr exec -n --file example.toml --target count-chars
```

Result:

```shell
echo "Some Content" | wc -c
```

If we want to add the option to customize the command, you can use a template action.

```toml
[count-chars]
description = "Count the number of characters in a string"

template = """
echo "{{.content}}" | wc -c
"""

[count-chars.values]
content = "Some Custom Content"
```

Result:

```shell
echo "Some Custom Content" | wc -c
```

Of course something like this is easy enough to be done with an environment variable or a number of other tools 
directly from the command line but what if you had a more complicated command string?

Take `docker build` for example.
There are a lot of ways to run this command which creates a lot of edge cases when you need to support multiple teams 
or projects with different needs.

Note: The `-` character after the start of an action `{{-` or at the end of an action `-}}` gets rid of the whitespace
character before or after respectively.
Most commands don't care about whitespace but if you want to keep the output clean, these can be useful.

## Examples

Providing a custom Dockerfile only when one is defined

```
docker build {{- if .dockerfile}} --file {{.dockerfile}} {{- end}} .
```

Having a custom Dockerfile and Target

```
docker build {{- if .dockerfile}} --file {{.dockerfile}} {{- end}} \
{{- if .target}} --target {{.target}} {{- end}} .
```

Multiple Build Arguments

Uses the `mflag` custom function built into Commander

```
docker build {{- if .dockerfile}} --file {{.dockerfile -}} {{end}} \
{{- if .target}} --target {{.target}} {{- end}} {{- if .buildArgs}} {{mflag .buildArgs -}} {{end}} .
```

Note that `mflag` expected the format: `<seperator><flag  value>[input values...]`

This is so you can be flexible with the separator character for cases when certain characters might conflict with the
splitting operation, producing unexpected results.

For example: `|--build-arg|key1=value1|key2=value2|`

Output: `--build-arg key1=value1 --build-arg key2=value2`

Here's an example of the full configuration file

```toml
[docker-build]
description = "build images using docker with custom options"

template = """
docker build {{- if .dockerfile}} --file {{.dockerfile -}} {{end}} \
{{- if .target}} --target {{.target}} {{- end}} {{- if .buildArgs}} {{mflag .buildArgs -}} {{end}} .
"""

[docker-build.values]
dockerfile = "Dockerfile-custom"
target = "custom-final"
buildArgs = "|--build-arg|key1=value1|key2=value2"
```

```shell
exec -n --file ./bin/example.toml --target "docker-build"
```

result:

```shell
docker build --file Dockerfile-custom --target custom-final --build-arg key1=value1 --build-arg key2=value2 --build-arg .
```
