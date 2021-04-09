package data

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-git/go-git/v5"
)

// parses repo name from version control url
func getRepoName(url string) string {
	url = strings.Replace(url, "https://", "", 1)
	result := strings.Split(url, "/")
	if len(result) < 2 {
		return ""
	}

	return result[2]
}

func getVersionControlLink(absolutePath, relativePath string) string {
	r, err := git.PlainOpenWithOptions(absolutePath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return ""
	}

	config, err := r.Config()

	if err != nil {
		return ""
	}

	head, err := r.Head()
	if err != nil {
		return ""
	}

	currentBranch := strings.Replace(head.Name().String(), "refs/heads/", "", 1)

	for _, remote := range config.Remotes {
		for _, url := range remote.URLs {
			if url != "" {
				gitUrl, err := parseUrl(url)
				if err != nil {
					return ""
				}

				return includePathToFile(gitUrl, currentBranch, relativePath)
			}
		}
	}
	return ""
}

func parseUrl(urlString string) (string, error) {
	gitSSH := strings.HasPrefix(urlString, "git@")
	gitRepo := strings.HasSuffix(urlString, ".git")

	if gitSSH && gitRepo {
		idx := strings.LastIndex(urlString, ":")

		urlString = strings.Replace(urlString, ":", "/", idx)
		// probably safe to assume https
		urlString = strings.ReplaceAll(urlString, "git@", "https://")
		urlString = strings.ReplaceAll(urlString, ".git", "")
	}

	uri, err := url.Parse(urlString)
	if err != nil {
		return "", fmt.Errorf("invalid urlString: %s", uri)
	}

	return uri.String(), nil
}

func includePathToFile(baseUrl, branchName, relativePathToFile string) string {
	if strings.Contains(baseUrl, "github") {
		return baseUrl + "/edit/" + branchName + "/" + relativePathToFile
	}

	if strings.Contains(baseUrl, "gitlab") {
		return baseUrl + "/~/edit" + branchName + "/" + relativePathToFile
	}
	return ""
}
