package javascript

import "encoding/json"

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

	return packageJson.Engines.Node
}
