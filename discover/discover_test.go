package discover_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/effxhq/effx-cli/data"
	"github.com/effxhq/effx-cli/discover"
	"github.com/stretchr/testify/require"
)

// func Test_Discover_Services_From_Yaml(t *testing.T) {
// 	dir, _ := ioutil.TempDir("", "services")
// 	defer os.RemoveAll(dir)

// 	dooku, _ := ioutil.TempDir(dir, "dooku")
// 	tedryn, _ := ioutil.TempDir(dir, "tedryn")
// 	_, _ = ioutil.TempDir(dir, "watto")

// 	tedrynFile, _ := ioutil.TempFile(tedryn, "effx.yaml")
// 	wattoFile, _ := ioutil.TempFile(dooku, "effx.yaml")

// 	input := []data.EffxYaml{
// 		{
// 			FilePath: tedrynFile.Name(),
// 		}, {
// 			FilePath: wattoFile.Name(),
// 		},
// 	}

// 	res := discover.DetectServicesFromEffxYamls(input, "key", "effx-cli")

// 	require.Len(t, res, 1)
// 	require.Contains(t, res[0].Name, "watto")
// }

func Test_Discover_Services(t *testing.T) {
	dir, _ := ioutil.TempDir("", "fakedir")
	defer os.RemoveAll(dir)

	_, _ = os.Create(dir + "/package.json")

	res, err := discover.DetectServicesFromFiles(dir, []data.EffxYaml{}, "effx-cli")

	require.Nil(t, err)
	require.Len(t, res, 1)
	require.Contains(t, res[0].Name, "fakedir")
}

func Test_Nested_DirectoryName(t *testing.T) {
	dir, _ := ioutil.TempDir("", "apps")
	defer os.RemoveAll(dir)

	_, _ = ioutil.TempDir(dir, "dooku")

	res, err := discover.DetectServicesFromFiles(dir, []data.EffxYaml{}, "effx-cli")

	require.Nil(t, err)
	require.Len(t, res, 1)
	require.Contains(t, res[0].Name, "dooku")
}
