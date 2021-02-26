package java

import (
	"regexp"
)

// there are three different ways to find the javen version
// from pom files:
// https://www.baeldung.com/maven-java-version

var JavaRegexList []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile(`<maven.compiler.target>\s*(.*?)\s*</maven.compiler.target>`),
	regexp.MustCompile(`<java.version>\s*(.*?)\s*</java.version>`),
	regexp.MustCompile(`<maven.compiler.release>\s*(.*?)\s*</maven.compiler.release>`),
}

// if it cannot find substring within two strings, return an empty string
func FindBetweenTwoStrings(r *regexp.Regexp, input string) string {
	matches := r.FindAllStringSubmatch(input, -1)

	if len(matches) < 1 || len(matches[0]) < 2 {
		return ""
	}

	return matches[0][1]
}

func HandlePomFile(fileContent string) string {
	for _, regex := range JavaRegexList {
		version := FindBetweenTwoStrings(regex, fileContent)
		if version != "" {
			return version
		}
	}

	return ""
}
