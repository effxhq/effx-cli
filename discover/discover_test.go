package discover_test

import (
	"os"
	"testing"

	"github.com/effxhq/effx-cli/data"
	"github.com/effxhq/effx-cli/discover"
	"github.com/stretchr/testify/require"
)

func Test_Discover_Services(t *testing.T) {
	dir := "fakedir"
	_ = os.Mkdir(dir, 0755)
	_, _ = os.Create(dir + "/package.json")
	defer os.RemoveAll(dir)

	res, err := discover.DetectServicesFromFiles(dir, []data.EffxYaml{}, "effx-cli")

	require.Nil(t, err)
	require.Len(t, res, 1)
	require.Contains(t, res[0].Name, "fakedir")
}

func Test_Nested_DirectoryName(t *testing.T) {
	dir := "apps"
	_ = os.Mkdir(dir, 0755)
	_ = os.Mkdir(dir+"/dooku", 0755)
	defer os.RemoveAll(dir)

	res, err := discover.DetectServicesFromFiles(dir, []data.EffxYaml{}, "effx-cli")

	require.Nil(t, err)
	require.Len(t, res, 1)
	require.Contains(t, res[0].Name, "dooku")
}

func Test_Discover_Services_From_Yaml(t *testing.T) {
	dir := "test"
	_ = os.Mkdir(dir, 0755)
	_ = os.Mkdir(dir+"/watto", 0755)

	defer os.RemoveAll(dir)

	input := []data.EffxYaml{
		{
			FilePath: dir + "/dooku/effx.yaml",
		}, {
			FilePath: dir + "/tedryn/effx.yaml",
		},
	}

	res := discover.DetectServicesFromEffxYamls(input, "key", "effx-cli")

	require.Len(t, res, 1)
	require.Contains(t, res[0].Name, "watto")
}
