# mflag

Multiple Flags

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
