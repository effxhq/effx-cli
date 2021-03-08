package discover

import "strings"

func New(val string) *Iterator {
	return &Iterator{val: val}
}

type Iterator struct {
	ptr int
	val string
}

func (i Iterator) HasNext() bool {
	return i.ptr < len(i.val)
}

func (i *Iterator) Peek() string {
	if !i.HasNext() {
		return ""
	}
	return i.val[i.ptr : i.ptr+1]
}

func (i *Iterator) Next() string {
	v := i.Peek()
	i.ptr++
	return v
}

func generateIterators(list []string) []*Iterator {
	result := []*Iterator{}

	for _, s := range list {
		result = append(result, New(s))
	}

	return result
}

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
