package command

import (
	"bytes"
	"testing"
)

func TestCommand_Print(t *testing.T) {

	testTable := []struct {
		input   *Command
		want    string
		wantErr bool
	}{
		{input: &Command{
			Template: `echo {{- if .msg}} "{{.msg}}"{{end}}`,
			Values:   map[string]string{"msg": "howdy world"},
		}, want: "echo \"howdy world\""},

		{input: &Command{
			Template: `echo {{- if .newLine}} -n{{end}} {{- if .msg}} "{{.msg}}"{{end}}`,
			Values:   map[string]string{"msg": "howdy world", "newLine": "true"},
		}, want: "echo -n \"howdy world\""},

		{input: &Command{
			Template: `docker build --file {{- if .buildFile}} {{.buildFile}} {{- else}} Dockerfile {{- end}} .`,
			Values:   map[string]string{},
		}, want: "docker build --file Dockerfile ."},

		{input: &Command{
			Template: `docker build --file {{- if .buildFile}} {{.buildFile}} {{- else}} Dockerfile {{- end}} .`,
			Values:   map[string]string{"buildFile": "Dockerfile-custom"},
		}, want: "docker build --file Dockerfile-custom ."},

		{input: &Command{
			Template:    `echo {{- if .newLine} -n{{end}} {{- if .msg}} "{{.msg}}"{{end}}`,
			Description: "An echo command",
			Values:      map[string]string{"msg": "howdy world", "newLine": "true"},
		}, wantErr: true},
	}

	for i, c := range testTable {
		buf := new(bytes.Buffer)
		err := c.input.Print(buf)
		got := buf.String()
		t.Logf("CASE %d\n%s\n", i, got)
		if c.wantErr && err == nil {
			t.Fatalf("want error. got nil")
		}

		if got != c.want {
			t.Fatalf("want: %v got: %v", c.want, got)
		}
	}
}
