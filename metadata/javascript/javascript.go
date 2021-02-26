package javascript

import (
	"encoding/json"
	"strings"
)

// should we use this?
// removes spaces and inequalites
var replacer = strings.NewReplacer(" ", "", ">", "", "<", "", "=", "")

type packageJson struct {
	Engines struct {
		Node string `json:"node"`
	} `json:"engines"`
}

func HandlePackageJson(fileContent string) string {
	var packageJson = &packageJson{}

	err := json.Unmarshal([]byte(fileContent), packageJson)
	if err != nil {
		return ""
	}

	versionString := packageJson.Engines.Node

	return replacer.Replace(versionString)
}
