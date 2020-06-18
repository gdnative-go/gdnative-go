// Copyright © 2019 - 2020 Oscar Campos <oscar.campos@thepimpam.com>
// Copyright © 2017 - William Edwards
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"gitlab.com/pimpam-games-studio/gdnative-go/gdnative"
)

// RegistryData structure is a container for registrable classes
// that is passed to the gdnative_wrapper.go.tmpl template for filling
type RegistryData struct {
	Package string
	Classes map[string]gdnative.Registrable
}

// GDNativeInit construct and returns a SetNativeInitScript call
func (rd RegistryData) GDNativeInit() string {

	initFunctions := []string{}
	for className := range rd.Classes {
		initFunctions = append(initFunctions, fmt.Sprintf("nativeScriptInit%s", className))
	}

	return fmt.Sprintf("gdnative.SetNativeScriptInit(%s)", strings.Join(initFunctions, ", "))
}

func (cmd *generateCmd) Run(ctx *context) error {

	fset := token.NewFileSet()
	packages, parseErr := parser.ParseDir(fset, ctx.Path, cmd.filter, parser.ParseComments)
	if parseErr != nil {
		return fmt.Errorf("could not parse Go files at %s: %w", ctx.Path, parseErr)
	}

	tplPath, pathErr := getTemplatePath("gdnative_wrapper.go")
	if pathErr != nil {
		return fmt.Errorf("could not get GDNative template: %w", pathErr)
	}

	for pkg, p := range packages {

		data := RegistryData{Package: pkg, Classes: map[string]gdnative.Registrable{}}
		registrable := gdnative.LookupRegistrableTypeDeclarations(p)
		if len(registrable) == 0 {
			fmt.Printf("not found any registrable sources on %s", ctx.Path)
			return nil
		}

		for className, classData := range registrable {
			data.Classes[className] = classData
		}

		// create a template from the template file
		tpl, tplErr := template.ParseFiles(tplPath)
		if tplErr != nil {
			return tplErr
		}

		outputFileName := fmt.Sprintf("%s_registrable.gen.go", pkg)
		outputFilePath := filepath.Join(ctx.Path, outputFileName)
		file, fileErr := os.Create(outputFilePath)
		if fileErr != nil {
			return fmt.Errorf("can not open output file %s for writing: %w", outputFilePath, fileErr)
		}

		execErr := tpl.Execute(file, data)
		if execErr != nil {
			return execErr
		}

		return format(outputFilePath)
	}

	return nil
}

// get the template path
func getTemplatePath(templateType string) (string, error) {

	currentPath, err := getCurrentPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(currentPath, "..", "..", "generate", "templates", templateType+".tmpl"), nil
}

// get the current in execution file path on disk
func getCurrentPath() (string, error) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not get current file execution path")
	}

	return filename, nil
}

func (cmd *generateCmd) filter(info os.FileInfo) bool {

	if info.IsDir() {
		return false
	}

	length := len(info.Name())
	if length > 7 && info.Name()[length-7:length-2] == ".gen." {
		return false
	}

	return true
}

// formats the given path with gofmt
func format(filepath string) error {

	fmt.Println("gofmt", "-w", filepath)
	cmd := exec.Command("gofmt", "-w", filepath)
	return cmd.Run()
}
