package java_test

import (
	"io/ioutil"
	"testing"

	"github.com/effxhq/effx-cli/pkg/metadata/java"
	"github.com/stretchr/testify/require"
)

func TestHandlePomFile(t *testing.T) {
	content, _ := ioutil.ReadFile("./pom.xml")
	res := java.HandlePomFile(string(content))

	require.Equal(t, res, "1.9")
}
