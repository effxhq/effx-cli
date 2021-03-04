package php_test

import (
	"io/ioutil"
	"testing"

	"github.com/effxhq/effx-cli/pkg/metadata/php"
	"github.com/stretchr/testify/require"
)

func TestHandleComposerJSON(t *testing.T) {
	content, _ := ioutil.ReadFile("./composer.json")
	res := php.HandleComposerJson(string(content))

	require.Equal(t, res, "5.6.1")
}
