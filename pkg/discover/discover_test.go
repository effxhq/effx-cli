package discover_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/effxhq/effx-cli/pkg/discover"
	"github.com/stretchr/testify/require"
)

func Test_Discover_Services(t *testing.T) {
	dir, _ := ioutil.TempDir("", "services")
	defer os.RemoveAll(dir)

	dooku, _ := ioutil.TempDir(dir, "dooku")
	tedryn, _ := ioutil.TempDir(dir, "tedryn")
	_, _ = ioutil.TempDir(dir, "watto")

	tedrynFile, _ := ioutil.TempFile(tedryn, "effx.yaml")
	wattoFile, _ := ioutil.TempFile(dooku, "effx.yaml")

	input := []string{tedrynFile.Name(), wattoFile.Name()}

	res := discover.DetectServices("effx-cli", input)

	require.Len(t, res, 1)
	require.Contains(t, res[0].Name, "watto")
}
