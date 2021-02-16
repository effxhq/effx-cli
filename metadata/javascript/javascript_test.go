package javascript_test

import (
	"io/ioutil"
	"testing"

	"github.com/effxhq/effx-cli/metadata/javascript"
	"github.com/stretchr/testify/require"
)

func TestHandlePackageJSON(t *testing.T) {

	content, _ := ioutil.ReadFile("./package.json")
	res := javascript.HandlePackageJson(string(content))

	require.Equal(t, res, "0.10.0")
}
