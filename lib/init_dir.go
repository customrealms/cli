package lib

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed initial/main.ts
var mainTs string

//go:embed initial/tsconfig.json
var tsconfigJson string

//go:embed initial/webpack.config.js
var webpackConfigJs string

type PackageJSONForInit struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Scripts         map[string]string `json:"scripts"`
	Keywords        []string          `json:"keywords"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func InitDir(dir, name string) error {

	// Setup the files
	if err := initDirFiles(dir, name); err != nil {
		return err
	}

	// Run npm install
	cmd := exec.Command("npm", "install")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Initialize the git repo
	cmd = exec.Command("git", "init")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func initDirFiles(dir, name string) error {

	// If there are already files in the dir
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if len(entries) > 0 {
		return errors.New("directory already contains files")
	}

	// Make the directory
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	// Write the .gitignore file

	gitignore := []string{
		"/node_modules",
		"/dist",
		".DS_Store",
		"",
	}
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(strings.Join(gitignore, "\n")), 0777); err != nil {
		return err
	}

	// Write the package.json file

	packageJson := PackageJSONForInit{
		Name:        name,
		Version:     "1.0.0",
		Description: "",
		Scripts: map[string]string{
			"build:jar": fmt.Sprintf("crx build -o ./dist/%s.jar", name),
			"build":     "webpack --mode=production",
			"clean":     "rm -rf ./dist",
		},
		Keywords: []string{},
		Author:   "",
		License:  "ISC",
		Dependencies: map[string]string{
			"@customrealms/core": "^0.1.0",
		},
		DevDependencies: map[string]string{
			"ts-loader":   "^9.2.6",
			"tslib":       "^2.3.1",
			"typescript":  "^4.4.4",
			"webpack":     "^5.63.0",
			"webpack-cli": "^4.9.1",
		},
	}
	jsonBytes, err := json.MarshalIndent(packageJson, "", "\t")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "package.json"), jsonBytes, 0777); err != nil {
		return err
	}

	// Write the tsconfig.json

	if err := os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte(tsconfigJson), 0777); err != nil {
		return err
	}

	// Write the webpack config

	if err := os.WriteFile(
		filepath.Join(dir, "webpack.config.js"),
		[]byte(webpackConfigJs),
		0777,
	); err != nil {
		return err
	}

	// Write the main.ts file

	if err := os.MkdirAll(filepath.Join(dir, "src"), 0777); err != nil {
		return err
	}

	if err := os.WriteFile(
		filepath.Join(dir, "src", "main.ts"),
		[]byte(mainTs),
		0777,
	); err != nil {
		return err
	}

	return nil

}
