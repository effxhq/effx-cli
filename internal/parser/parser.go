package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	"github.com/effxhq/effx-cli/data"
	gyaml "github.com/ghodss/yaml"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type EffxYaml struct {
	FilePath string
}

func (y EffxYaml) getFilePattern() string {
	return "(.+\\.)?effx\\.ya?ml$"
}

func (y EffxYaml) isEffxYaml() (bool, error) {
	pattern := y.getFilePattern()
	matched, err := regexp.MatchString(pattern, y.FilePath)
	return matched, err
}

func (y EffxYaml) ToApiResources() ([]data.ApiResource, error) {
	ok, err := y.isEffxYaml()
	if !ok {
		pattern := y.getFilePattern()
		errString := fmt.Sprintf("Not an Effx Yaml. %s must match pattern: %s", y.FilePath, pattern)
		return nil, errors.New(errString)
	}
	if err != nil {
		return nil, err
	}

	var docs []data.ApiResource

	if yamlFile, err := ioutil.ReadFile(y.FilePath); err != nil {
		return docs, err
	} else {
		buf := bytes.NewBuffer(yamlFile)
		reader := utilyaml.NewYAMLReader(bufio.NewReader(buf))
		for {
			b, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			if len(b) == 0 {
				break
			}
			jsonBytes, _ := gyaml.YAMLToJSON(b)
			resource := data.ApiResourceContent{Content: jsonBytes}
			docs = append(docs, resource)
		}
	}
	return docs, nil
}
