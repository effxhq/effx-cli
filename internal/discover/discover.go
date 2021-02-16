package discover

import (
	"io/ioutil"
	"strings"
)

func findCommonDirectory(effxFileLocations []string) string {
	matchedEffxFiles := generateIterators(effxFileLocations)
	prefixString := ""

	for len(matchedEffxFiles) > 0 {
		count := make(map[string]int)
		for _, matchedFile := range matchedEffxFiles {
			peek := matchedFile.Peek()
			if peek != "" {
				count[peek]++
			}
		}

		maxK := ""
		maxV := 1
		for k, v := range count {
			if v > maxV {
				maxK = k
				maxV = v
			}
		}

		nextRound := make([]*Iterator, 0, maxV)
		for _, matchedFile := range matchedEffxFiles {
			// advance ptr
			if matchedFile.Next() == maxK {
				// put into next
				nextRound = append(nextRound, matchedFile)
			}
		}

		prefixString += maxK
		matchedEffxFiles = nextRound

	}

	if prefixString == "" {
		return ""
	}

	// prefix string should be a directory ending with a slash
	slashIndex := strings.LastIndex(prefixString, "/")

	if slashIndex != len(prefixString) {
		// trim file name, keep last dir slash
		// example:
		// services/dooku -> services/
		prefixString = prefixString[:slashIndex+1]
	}

	return prefixString
}

func DetectServices(effxFileLocations []string) []string {
	detectedServiceNames := []string{}

	commonDir := findCommonDirectory(effxFileLocations)

	files, err := ioutil.ReadDir(commonDir)
	if err != nil {
		return []string{}
	}

	for _, file := range files {
		// looking at directories only for service locations
		if file.IsDir() {
			contains := false
			for _, effxFileLocation := range effxFileLocations {
				if strings.Contains(effxFileLocation, file.Name()) {
					contains = true
				}
			}
			if !contains {
				detectedServiceNames = append(detectedServiceNames, file.Name())
			}
		}
	}

	return detectedServiceNames
}
