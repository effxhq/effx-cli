package parser

import (
	"io/ioutil"

	"github.com/effxhq/effx-go/data"
	"gopkg.in/yaml.v2"
)

func YamlFile(filePath string) (*data.Data, error) {
	res := &data.Data{}

	if yamlFile, err := ioutil.ReadFile(filePath); err != nil {
		return res, err
	} else {
		if err := yaml.Unmarshal(yamlFile, res); err != nil {
			return res, err
		}
	}

	return res, nil
}
