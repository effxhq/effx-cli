package golang

import (
	"regexp"
	"strings"
)

// regex to find go version from go.mod file
const GoVersionRegex = `go \d+(\.\d+)+`

var goModRegex = regexp.MustCompile(GoVersionRegex)

func HandleGoModFile(fileContent string) string {
	res := goModRegex.FindString(fileContent)
	versionStr := strings.Split(res, " ")

	if len(versionStr) < 2 {
		return ""
	}

	return versionStr[1]
}
