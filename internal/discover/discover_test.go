package discover_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/effxhq/effx-cli/internal/discover"
	"github.com/stretchr/testify/require"
)

func TestConsumer_SetupFS(t *testing.T) {
	dir, _ := ioutil.TempDir("", "services")
	defer os.RemoveAll(dir)

	dooku, _ := ioutil.TempDir(dir, "dooku")
	tedryn, _ := ioutil.TempDir(dir, "tedryn")
	_, _ = ioutil.TempDir(dir, "watto")

	tedrynFile, _ := ioutil.TempFile(tedryn, "effx.yaml")
	wattoFile, _ := ioutil.TempFile(dooku, "effx.yaml")

	input := []string{tedrynFile.Name(), wattoFile.Name()}

	res := discover.DetectServices(input)

	require.Len(t, res, 1)
	require.Contains(t, res[0], "watto")
}
