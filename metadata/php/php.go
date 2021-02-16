package php

import (
	"encoding/json"
)

type composerJson struct {
	Config struct {
		Platform struct {
			Php string `json:"php"`
		} `json:"platform"`
	} `json:"config"`
}

func HandleComposerJson(fileContent string) string {
	var composerJson = &composerJson{}

	err := json.Unmarshal([]byte(fileContent), composerJson)
	if err != nil {
		return ""
	}

	return composerJson.Config.Platform.Php
}
