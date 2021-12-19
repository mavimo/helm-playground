package lib

import (
	"bytes"
	"encoding/json"
	"text/template"

	"gopkg.in/yaml.v2"
)

type ValuesObj map[string]interface{}

type TemplateData struct {
	Values ValuesObj
}

type GetYamlReturnValue struct {
	Yaml string `json:"yaml"`
	Err  string `json:"err"`
}

func toJson(returnValue GetYamlReturnValue) string {
	bytes, err := json.Marshal(returnValue)
	if err != nil {
		return `{"err":"conversion to JSON failed"}`
	}
	return string(bytes)
}

func GetYaml(templateYaml string, valuesYaml string) string {
	valuesData := ValuesObj{}
	if err := yaml.Unmarshal([]byte(valuesYaml), &valuesData); err != nil {
		return toJson(GetYamlReturnValue{
			Err: err.Error(),
		})
	}

	templateData := TemplateData{valuesData}

	var output bytes.Buffer

	t, err := template.New("template").Funcs(funcMap()).Parse(templateYaml)
	if err != nil {
		return toJson(GetYamlReturnValue{
			Err: err.Error(),
		})
	}

	if err := t.Execute(&output, templateData); err != nil {
		return toJson(GetYamlReturnValue{
			Err: err.Error(),
		})
	}

	return toJson(GetYamlReturnValue{
		Yaml: output.String(),
	})
}