package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"example/pkg/command"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestWriteExample(t *testing.T) {
	testTable := []struct {
		name     string
		fileType FileType
	}{
		{name: "json", fileType: FileTypeJSON},
		{name: "yaml", fileType: FileTypeYAML},
		{name: "toml", fileType: FileTypeTOML},
		{name: "json-ext-1", fileType: ParseFileType("json")},
		{name: "json-ext-2", fileType: ParseFileType(".json")},
		{name: "yaml-ext-1", fileType: ParseFileType("yaml")},
		{name: "yaml-ext-2", fileType: ParseFileType("yml")},
		{name: "yaml-ext-3", fileType: ParseFileType(".yaml")},
		{name: "yaml-ext-4", fileType: ParseFileType(".yml")},
		{name: "toml-ext-1", fileType: ParseFileType("toml")},
		{name: "toml-ext-2", fileType: ParseFileType(".toml")},
	}

	for _, c := range testTable {
		t.Run(c.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := EncodeExample(NewTypeWriter(buf, c.fileType)); err != nil {
				t.Fatal(err)
			}
			t.Log(buf.String())
			var commandMap map[string]*command.Command
			if err := NewTypedReader(buf, c.fileType).Decode(&commandMap); err != nil {
				t.Fatal(err)
			}
			t.Logf("%+v", commandMap)
			want := command.ExampleCmd
			got := commandMap
			if !reflect.DeepEqual(want, got) {
				t.Fatalf("want: %v got: %v", want, got)
			}

		})
	}

}

func TestFprintCmd(t *testing.T) {
	config := map[string]*command.Command{
		"one": {
			Template: "echo {{- if .newline}} -n {{end -}} {{.msg}}",
			Values:   map[string]string{"msg": "some message", "newline": "true"}},
		"two": {
			Template: "",
			Values:   map[string]string{},
		},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(&config); err != nil {
		t.Fatal(err)
	}
	jsonBytes := bytes.Clone(buf.Bytes())
	t.Log(buf.String())

	t.Run("success", func(t *testing.T) {
		outBuf := new(bytes.Buffer)
		if err := FprintCmd(outBuf, NewTypedReader(bytes.NewBuffer(jsonBytes), FileTypeJSON), "one"); err != nil {
			t.Fatal(err)
		}
		t.Log(outBuf.String())
	})

	t.Run("decode-error", func(t *testing.T) {
		if err := FprintCmd(io.Discard, NewTypedReader(bytes.NewBuffer(jsonBytes), FileTypeJSON), "one"); err == nil {
			t.Fatal("want decoding error, got nil")
		}
	})
	t.Run("missing-target", func(t *testing.T) {
		if err := FprintCmd(io.Discard, NewTypedReader(bytes.NewBuffer(jsonBytes), FileTypeJSON), "blah"); err == nil {
			t.Fatal("want missing target error, got nil")
		}
	})
}

func TestNewSystemCmd(t *testing.T) {
	config := map[string]*command.Command{
		"one": {
			Template: "echo {{- if .newline}} -n {{end -}} {{.msg}}",
			Values:   map[string]string{"msg": "some message", "newline": "true"}},
		"two": {
			Template: "echo {{- if .newline}} -n {{end -}} ${COMMANDER_TEST_MSG}",
			Values:   map[string]string{"newline": "true"},
		},
	}

	inBuf := new(bytes.Buffer)
	if err := json.NewEncoder(inBuf).Encode(&config); err != nil {
		t.Fatal(err)
	}
	jsonBytes := bytes.Clone(inBuf.Bytes())

	t.Run("success", func(t *testing.T) {
		cmd, err := NewSystemCmd(NewTypedReader(bytes.NewBuffer(jsonBytes), FileTypeJSON), "one")
		if err != nil {
			t.Fatal(err)
		}
		got := cmd.String()
		// Should pass visual test, can't be sure that echo will always bee in /bin/echo
		t.Log("cmd:", got)
		want := "echo -n some message"
		if !strings.Contains(cmd.String(), want) {
			t.Fatalf("want: %s got: %s", want, got)
		}
	})

	t.Run("success-env-expansion", func(t *testing.T) {
		os.Setenv("COMMANDER_TEST_MSG", "test message")
		cmd, err := NewSystemCmd(NewTypedReader(bytes.NewBuffer(jsonBytes), FileTypeJSON), "two")
		if err != nil {
			t.Fatal(err)
		}
		got := cmd.String()
		// Should pass visual test, can't be sure that echo will always bee in /bin/echo
		t.Log("cmd:", got)
		want := "echo -n test message"
		if !strings.Contains(cmd.String(), want) {
			t.Fatalf("want: %s got: %s", want, got)
		}
	})
}

type mockReader struct{ readFunc func([]byte) (int, error) }

func (m mockReader) Read(p []byte) (int, error) { return m.readFunc(p) }

func errorReader() mockReader {
	return mockReader{readFunc: func(b []byte) (int, error) { return 0, errors.New("mock reader: error") }}
}
