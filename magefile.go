// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var ldflags = "-s -w -X main.version=$VERSION_TAG -X main.commit=$COMMIT_HASH -X main.date=$BUILD_TIME -X main.builtBy=$BUILDER"

var Default = Generate

func init() {

	// make sure we use Go 1.11 modules even if the source lives inside GOPATH
	os.Setenv("GO111MODULE", "on")
}

// Generate Go Bindings
func Generate() error {

	mg.SerialDeps(RetrieveGodotDocumentation, Clean)
	apiPath := getCurrentFilePath()
	generateScript := filepath.Join(getCurrentFilePath(), "cmd", "generate", "main.go")
	err := sh.RunWith(map[string]string{"API_PATH": apiPath}, "go", "run", "-v", generateScript)
	if err != nil {
		return fmt.Errorf("could not genearate Go bindings: %w", err)
	}

	return Build()
}

// Clean cleans previous generations
func Clean() error {

	log.Println("Cleaning previous generation...")
	path := filepath.Join(getCurrentFilePath(), "gdnative", "*.gen.*")
	files, globErr := filepath.Glob(path)
	if globErr != nil {
		return globErr
	}
	for _, filename := range files {
		if err := sh.Rm(filename); err != nil {
			return err
		}
	}
	return nil
}

// RetrieveGodotDocumentation retrieves latest Godot documentation to attach docstrings
func RetrieveGodotDocumentation() error {

	localPath := getCurrentFilePath()
	docPath := filepath.Join(localPath, "doc")
	_, found := os.Stat(docPath)
	if found == nil {
		_ = os.Chdir(docPath)
		log.Println("Godot documentation found. Pulling latest changes...")
		if err := sh.Run("git", "pull", "origin", "master"); err != nil {
			return fmt.Errorf("could not pull latest Godot documentation from git: %w", err)
		}
		_ = os.Chdir(localPath)
		return nil
	}

	log.Println("Godot documentation not found. Cloning the repository...")
	if err := os.MkdirAll(docPath, 0766); err != nil {
		return fmt.Errorf("could not create a new directory on the disk: %w", err)
	}
	_ = os.Chdir(docPath)
	if err := sh.Run("git", "init"); err != nil {
		return fmt.Errorf("could not execute git init: %w", err)
	}
	if err := sh.Run("git", "remote", "add", "-f", "origin", "https://github.com/godotengine/godot.git"); err != nil {
		return fmt.Errorf("could not set origin remote for documentation: %w", err)
	}
	if err := sh.Run("git", "config", "core.sparseCheckout", "true"); err != nil {
		return fmt.Errorf("could not activate core.sparseCheckout: %w", err)
	}
	sparseCheckoutsConfigFile := filepath.Join(".", ".git", "info", "sparse-checkout")
	writeErr := ioutil.WriteFile(sparseCheckoutsConfigFile, []byte("doc/classes"), 0600)
	if writeErr != nil {
		return fmt.Errorf("could not write .git/info/sparse-checkout file: %w", writeErr)
	}
	if err := sh.Run("git", "pull", "origin", "master"); err != nil {
		return fmt.Errorf("error while pulling: %w", err)
	}

	return nil
}

// Build builds the gdnative-go compiler gogdc (also builds the library)
func Build() error {
	return sh.RunWith(flagEnv(), "go", "build", "-ldflags", ldflags, "-x", "./cmd/gogdc")
}

// Install builds and installs the gdnative-go compiler gogdc in $GOPATH/bin
func Install() error {
	return sh.RunWith(flagEnv(), "go", "install", "-ldflags", ldflags, "-x", "./cmd/gogdc")
}

// getCurrentFilePath constructs and returns the current file path on the drive
func getCurrentFilePath() string {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("could not get current file path")
	}

	return filepath.Join(filepath.Dir(filename))
}

// fills environment with build data
func flagEnv() map[string]string {

	commitHash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	buildAuthor, _ := sh.Output("git", "log", "-1", "--pretty=format:%ae")
	versionTag, _ := sh.Output("git", "describe", "--tags", "--abbrev=0")
	if versionTag == "" {
		versionTag = "v0.0.0-dev"
	}

	return map[string]string{
		"COMMIT_HASH": commitHash,
		"VERSION_TAG": versionTag,
		"BUILD_TIME":  time.Now().Format("2006-01-02T15:04:05Z0700"),
		"BUILDER":     buildAuthor,
		"CGO_ENABLED": "1",
	}
}
