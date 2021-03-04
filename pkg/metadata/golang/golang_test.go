package golang_test

import (
	"io/ioutil"
	"testing"

	"github.com/effxhq/effx-cli/pkg/metadata/golang"
	"github.com/stretchr/testify/require"
)

func TestHandleGoMod(t *testing.T) {
	// using the one in this repo rather than an
	// example
	content, _ := ioutil.ReadFile("../../go.mod")
	res := golang.HandleGoModFile(string(content))

	require.Equal(t, res, "1.13")
}
