package command

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"reflect"
	"strings"
)

// Command is a data representation of a single command
type Command struct {
	Description string            `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Template    string            `json:"template" yaml:"template" toml:"template"`
	Values      map[string]string `json:"values" yaml:"values" toml:"values" `
}

// Print parses the Template field and executes uses the Values field
func (c *Command) Print(w io.Writer) error {
	var err error
	tmpl := template.New("command")
	tmpl = tmpl.Funcs(template.FuncMap{"mflag": multiFlag})
	tmpl, err = tmpl.Parse(c.Template)
	if err != nil {
		slog.Error("failed to parse template", "err", err)
		return err
	}
	return tmpl.Execute(w, c.Values)
}

// indirect returns the item at the end of indirection, and a bool to indicate
// if it's nil. If the returned bool is true, the returned value's kind will be
// either a pointer or interface.
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}

// multiFlag ...
func multiFlag(item reflect.Value) (string, error) {
	item, isNil := indirect(item)
	if isNil {
		return "", fmt.Errorf("args of nil pointer")
	}
	if item.Kind() != reflect.String {
		return "", fmt.Errorf("args of type %s", item.Type())
	}

	inputString := item.String()
	splitChar := []rune(inputString)[0]

	parts := strings.Split(inputString, string(splitChar))
	if len(parts) < 3 {
		return "", fmt.Errorf("got: %s want format like: --arg \"value 1\" \"value 2\"", inputString)
	}

	flagString := parts[1]
	out := make([]string, 0, len(parts)*2)

	// 0 is Blank since the first character is the split character
	// 1 is the flag arg
	for _, part := range parts[2:] {
		out = append(out, flagString, part)
	}

	return strings.Join(out, " "), nil
}

var ExampleCmd = map[string]*Command{
	"echo": {
		Description: "print a text message",
		Template:    "echo {{- if .newline}} -n {{end -}} \"{{.msg}}\"",
		Values:      map[string]string{"msg": "howdy world", "newline": "true"},
	},
	"docker-build": {
		Description: "build an image with docker",
		Template:    "docker build {{- if .file}} --file {{.file}} {{end}} {{- if .context}} {{.context}} {{else}} . {{- end}}",
		Values:      map[string]string{"file": "Dockerfile.custom"},
	},
}
