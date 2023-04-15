//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

var workingDir = func() string {
	name, _ := os.Getwd()
	return name
}()

const rootModule = "github.com/alexiscampan/go.pkg"

var toolsBinDir = path.Join(workingDir, "tools", "bin")

// Function to init tools path
func init() {
	time.Local = time.UTC

	os.Setenv("PATH", fmt.Sprintf("%s:%s", toolsBinDir, os.Getenv("PATH")))
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	banner := figure.NewColorFigure("go.pkg", "rectangles", "green", true)
	banner.Print()
	fmt.Println("")
	color.Red("---------------------Building 🏗 ---------------------")
	for _, module := range listGoModules() {
		if module == "." {
			continue
		}
		os.Chdir(filepath.Join(workingDir, module))
		color.Blue("Module \"%s\" found ✅", module)
		err := Tidy()
		if err != nil {
			return err
		}
		err = Vendor()
		if err != nil {
			return err
		}
		err = Format()
		if err != nil {
			return err
		}
		err = Lint()
		if err != nil {
			return err
		}
	}
	return nil
}

// Function to check the packages licenses
func License() error {
	color.Cyan("---------------------Check license 📜---------------------")
	return sh.RunV("wwhrd", "check")
}

// Function to update dependencies
func Tidy() error {
	color.Cyan("---------------------Updating dependencies 🔄---------------------")
	return sh.RunV("go", "mod", "tidy", "-v")
}

// Function to lint the code
func Lint() error {
	color.Cyan("---------------------Lint code 🚮---------------------")
	return sh.RunV("golangci-lint", "run")
}

// Function to verify dependencies
func Verify() error {
	color.Cyan("---------------------Verifying dependencies 🔐---------------------")
	return sh.RunV("go", "mod", "verify")
}

// Function to vendor the dependencies
func Vendor() error {
	color.Cyan("---------------------Vendoring dependencies 🧙‍♂️---------------------")
	return sh.RunV("go", "mod", "vendor")
}

// Function to format the go files
func Format() error {
	color.Cyan("---------------------Format go files ♻---------------------")
	files := []string{"-w"}
	files = append(files, listGoFiles()...)
	return sh.RunV("gofumpt", files...)
}

// Function to list all go files that need to be format
func listGoFiles() []string {
	var goFiles []string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "vendor/") {
			return filepath.SkipDir
		}
		if strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})
	return goFiles
}

// Function to list go modules
func listGoModules() []string {
	list, _ := sh.Output("go", "list", "-m")
	lines := strings.Split(list, "\n")

	modules := make([]string, len(lines))
	for i, line := range lines {
		modules[i] = path.Clean(strings.TrimPrefix(strings.TrimPrefix(line, rootModule), "/"))
	}
	return modules
}
