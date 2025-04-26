package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/theraccoonbear/CarDan/cardan"
)

const taskTemplate = `- {{ .Name }}
  Owner: {{ .Owner }}
  Duration: {{ .Duration }}
  Start: {{ .Start }}
  Completed: {{ if .CompletedOn }}{{ .CompletedOn }}{{ else }}(in progress){{ end }}
  Depends on:{{ if .DependsOn }}
{{- range .DependsOn }}
    - {{ . }}
{{- end }}
{{ else }} none{{ end }}

`

type Schema struct {
	People    map[string]string `yaml:"people"`
	Durations map[string]int    `yaml:"durations"`
	Tasks     []*Task           `yaml:"tasks"`
}

type Task struct {
	Name        string   `yaml:"name"`
	Owner       string   `yaml:"owner"`
	Duration    int      `yaml:"duration"`
	DependsOn   []string `yaml:"depends_on"`
	Start       string   `yaml:"start"`
	CompletedOn *string  `yaml:"completed_on"`
}

func loadSchema(path string) (*Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, err := cardan.Load(file)
	if err != nil {
		return nil, err
	}

	// Resolve references for "depends_on" fields
	if err := doc.ResolveRefs("depends_on"); err != nil {
		return nil, err
	}

	var schema Schema
	if err := doc.Unmarshal(&schema); err != nil {
		return nil, err
	}
	return &schema, nil
}

func main() {
	schema, err := loadSchema("test_data\\sample_schema.yml")
	if err != nil {
		log.Fatalf("failed to load schema: %v", err)
	}

	fmt.Println("Tasks Overview:")
	tpl := template.Must(template.New("task").Parse(taskTemplate))

	for _, task := range schema.Tasks {
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, task); err != nil {
			log.Fatalf("failed to render template: %v", err)
		}
		fmt.Print(buf.String())
	}

}
