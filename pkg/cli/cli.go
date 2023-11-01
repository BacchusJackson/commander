package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"example/pkg/command"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Encoding is a standard identifier for the supported file types
type Encoding string

const (
	EncodingJSON Encoding = "json"
	EncodingYAML Encoding = "yaml"
	EncodingTOML Encoding = "toml"
	EncodingUnsp Encoding = "unsupported"
)

// ParseEncoding will attempt conversion from a string to a FileType, if it fails
// to match based on a string or file extension, it will return FileTypeUnsp
func ParseEncoding(s string) Encoding {
	switch strings.TrimSpace(strings.ToLower(s)) {
	case "json", ".json":
		return EncodingJSON
	case "yaml", "yml", ".yaml", ".yml":
		return EncodingYAML
	case "toml", ".toml":
		return EncodingTOML
	}
	return EncodingUnsp
}

// EncodeExample will encode a JSON example configuration
// Only the first given outputFormat will be used if defined, otherwise
// the default output format is JSON
func EncodeExample(w *TypedWriter) error {
	return w.Encode(command.ExampleCmd)
}

// NewSystemCmd will parse a configuration from `in` and create a system
// command based on the selected target.
//
// Environment variables are already expanded before the command is generated.
// Always assumes JSON format
func NewSystemCmd(in *TypedReader, target string) (*exec.Cmd, error) {
	cmdStr, err := PrintCmd(in, target)
	if err != nil {
		return nil, err
	}

	cmdStr = os.ExpandEnv(cmdStr)

	cmdParts := strings.Split(cmdStr, " ")
	cmd := exec.Command("")

	switch len(cmdParts) {
	case 0:
		slog.Warn("command to execute is blank")
		return cmd, nil
	case 1:
		cmd = exec.Command(cmdParts[0])
	default:
		cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	}

	return cmd, nil

}

// PrintCmd returns the rendered string form of the command
// Always assumes JSON format for in
func PrintCmd(in *TypedReader, target string) (string, error) {
	buf := new(bytes.Buffer)
	err := FprintCmd(buf, in, target)
	return buf.String(), err
}

// FprintCmd writes the rendered command string to the writer
// w based on the config parsed from in and the target
func FprintCmd(w io.Writer, in *TypedReader, target string) error {

	config := make(map[string]*command.Command)

	if err := in.Decode(&config); err != nil {
		return err
	}

	cmd, ok := config[target]

	if !ok {
		return fmt.Errorf("command not found: '%s'", target)
	}

	return cmd.Print(w)
}

// TypedWriter associates a writer with an output encoding which enables
// it to Encode any object in that encoding
type TypedWriter struct {
	w              io.Writer
	OutputEncoding Encoding
}

// NewTypedWriter creates a TypedWriter with the associated encoding
func NewTypedWriter(w io.Writer, inputEnc Encoding) *TypedWriter {
	return &TypedWriter{w: w, OutputEncoding: inputEnc}
}

// Write wraps the internal writer
func (t *TypedWriter) Write(p []byte) (int, error) {
	return t.w.Write(p)
}

// Encode to the configured Encoding type
func (t *TypedWriter) Encode(v any) error {
	switch t.OutputEncoding {
	case EncodingJSON:
		enc := json.NewEncoder(t.w)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	case EncodingYAML:
		return yaml.NewEncoder(t.w).Encode(v)
	case EncodingTOML:
		enc := toml.NewEncoder(t.w)
		enc.Indent = ""
		return enc.Encode(v)
	}
	return errors.New("unsupported input encoding")
}

// TypedReader associates a reader with an input encoding which enables it
// to Decode from a specific encoding
type TypedReader struct {
	r             io.Reader
	InputEncoding Encoding
}

// NewTypedReader creates a reader / decoder with the configured encoding
func NewTypedReader(r io.Reader, inputEnc Encoding) *TypedReader {
	return &TypedReader{r: r, InputEncoding: inputEnc}
}

// Read wraps the internal reader
func (t *TypedReader) Read(p []byte) (int, error) {
	return t.r.Read(p)
}

// Decode from the internal reader using the specified encoding
func (t *TypedReader) Decode(v any) error {
	switch t.InputEncoding {
	case EncodingJSON:
		return json.NewDecoder(t.r).Decode(v)
	case EncodingYAML:
		return yaml.NewDecoder(t.r).Decode(v)
	case EncodingTOML:
		_, err := toml.NewDecoder(t.r).Decode(v)
		return err
	}
	return errors.New("unsupported input encoding")
}
