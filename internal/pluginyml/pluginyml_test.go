package pluginyml_test

import (
	"embed"
	"path"
	"testing"

	"github.com/customrealms/cli/internal/pluginyml"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed testdata/*
	testdataFS embed.FS
)

func readTestFile(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := testdataFS.ReadFile(path.Join("testdata", filename))
	if err != nil {
		t.Fatalf("read test file: %v", err)
	}
	return data
}

func TestUnmarshalPluginYml(t *testing.T) {
	t.Run("plugin-1.yml", func(t *testing.T) {
		// Unmarshal the plugin-1.yml file
		var plugin pluginyml.Plugin
		err := yaml.Unmarshal(readTestFile(t, "plugin-1.yml"), &plugin)
		require.NoError(t, err, "unmarshal plugin-1.yml")

		// Check the plugin details
		require.Equal(t, "ScrapBukkit", plugin.Name)
		require.Equal(t, "1.0.0", plugin.Version)
		require.Equal(t, "com.dinnerbone.bukkit.scrap.ScrapBukkit", plugin.Main)
		require.NotNil(t, plugin.Description)
		require.Equal(t, "Miscellaneous administrative commands for Bukkit. This plugin is one of the default plugins shipped with Bukkit.\n", *plugin.Description)

		// No commands in this plugin
		require.Nil(t, plugin.Commands)
		require.Len(t, plugin.Commands, 0)

		// Check the permissions
		require.NotNil(t, plugin.Permissions)
		rootPerm := plugin.Permissions["scrapbukkit.*"]
		require.NotNil(t, rootPerm)
		require.Equal(t, "Gives all permissions for Scrapbukkit", *rootPerm.Description)
		require.Equal(t, "op", *rootPerm.Default)
		require.Len(t, rootPerm.Children, 8)

		// Check the children permissions (nested permission blocks)
		require.NotNil(t, rootPerm.Children["scrapbukkit.remove"].Permission)
		require.NotNil(t, rootPerm.Children["scrapbukkit.time"].Permission)
		require.NotNil(t, rootPerm.Children["scrapbukkit.tp"].Permission)
		require.NotNil(t, rootPerm.Children["scrapbukkit.give"].Permission)
		require.NotNil(t, rootPerm.Children["scrapbukkit.clear"].Permission)

		// Check the children permissions (boolean permissions)
		require.NotNil(t, rootPerm.Children["scrapbukkit.some.standard.perm"].Bool)
		require.NotNil(t, rootPerm.Children["scrapbukkit.some.other.perm"].Bool)
		require.NotNil(t, rootPerm.Children["scrapbukkit.some.bad.perm"].Bool)
		require.Equal(t, true, *rootPerm.Children["scrapbukkit.some.standard.perm"].Bool)
		require.Equal(t, true, *rootPerm.Children["scrapbukkit.some.other.perm"].Bool)
		require.Equal(t, false, *rootPerm.Children["scrapbukkit.some.bad.perm"].Bool)

		// Check a nested child permission
		child := rootPerm.Children["scrapbukkit.time"].Permission
		require.Equal(t, "Allows the player to view and change the time", *child.Description)
		require.Len(t, child.Children, 2)
		require.Equal(t, "Allows the player to view the time", *child.Children["scrapbukkit.time.view"].Permission.Description)
		require.Equal(t, "true", *child.Children["scrapbukkit.time.view"].Permission.Default)
		require.Equal(t, "Allows the player to change the time", *child.Children["scrapbukkit.time.change"].Permission.Description)
		require.Nil(t, child.Children["scrapbukkit.time.change"].Permission.Default)

		// require.NotNil(t, plugin.Permissions["scrapbukkit.*"].Children["scrapbukkit.remove"].Permission)
		// require.Nil(t, plugin.Permissions["scrapbukkit.*"].Children["scrapbukkit.remove"].Bool)
	})
}
