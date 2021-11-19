package initialize

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/customrealms/cli/actions/initialize/template"
)

type InitAction struct {
	Name        string
	Dir         string
	Template    template.Template
	CoreVersion string
	CliVersion  string
}

func (a *InitAction) Run(ctx context.Context) error {

	// Check if NPM is installed on the machine
	if _, err := exec.LookPath("npm"); err != nil {
		fmt.Println("Couldn't find 'npm' command on your machine. Make sure NodeJS is installed.")
		fmt.Println("Visit https://nodejs.org and download the most recent version.")
		return nil
	}

	// If the template is nil, use the default template
	if a.Template == nil {
		tmpl, err := template.NewFromGitHub("customrealms", "cli-default-template", "master")
		if err != nil {
			return nil
		}
		a.Template = tmpl
	}

	// Install the template in the directory
	err := template.Install(
		a.Template,
		a.Dir,
		&template.Options{
			Name:        a.Name,
			CoreVersion: "^0.1.0",
			CliVersion:  "^0.3.0",
		},
	)
	if err != nil {
		return err
	}

	// Run npm install
	cmd := exec.Command("npm", "install")
	cmd.Dir = a.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// If Git is installed on the machine
	if _, err := exec.LookPath("git"); err == nil {

		// Initialize the git repo
		cmd = exec.Command("git", "init")
		cmd.Dir = a.Dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}

	}

	return nil
}
