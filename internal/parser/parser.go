package parser

import (
	"bytes"
	"io/ioutil"

	"github.com/Velocidex/yaml"
	"github.com/effxhq/effx-cli/data"
)

func YamlFile(filePath string) ([]*data.Data, error) {
	res := []*data.Data{}

	if yamlFile, err := ioutil.ReadFile(filePath); err != nil {
		return res, err
	} else {
		r := bytes.NewReader(yamlFile)
		dec := yaml.NewDecoder(r)

		for {
			d := &data.Data{}

			decoded := dec.Decode(d)

			if decoded != nil {
				break
			}

			res = append(res, d)
		}
	}

	return res, nil
}
