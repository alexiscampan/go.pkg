//go:build mage
// +build mage

package main

import (
	"github.com/fatih/color"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() {
	color.Red("Installing tools ...")
	mg.SerialDeps(Golang.Vendor, Golang.Tools)
}

type Golang mg.Namespace

// Function to vendor dependecies from tools
func (Golang) Vendor() error {
	color.Blue("Vendoring dependencies ...")
	return sh.RunV("go", "mod", "vendor")
}

// Function to generate tools file
func (Golang) Tools() error {
	color.Blue("Installing tools ...")
	return sh.RunV("go", "generate", "./tools.go")
}
