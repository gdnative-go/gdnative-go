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
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"gitlab.com/pimpam-games-studio/gdnative-go/gdnative"
)

func (cmd *listCmd) Run(ctx *context) error {

	fset := token.NewFileSet()
	packages, parseErr := parser.ParseDir(fset, ctx.Path, cmd.filter, parser.ParseComments)
	if parseErr != nil {
		return fmt.Errorf("could not parse Go files at %s: %w", ctx.Path, parseErr)
	}

	for pkg, p := range packages {

		fmt.Printf("Analyzing package: %s\n", pkg)
		gdregistrable := gdnative.LookupRegistrableTypeDeclarations(p)
		for key, data := range gdregistrable {
			base := data.GetBase()
			if base != "" {
				base = fmt.Sprintf("(%s)", base)
			}

			fmt.Printf("Found Structure: %s%s\n", key, base)

			properties := data.GetProperties()
			if len(properties) > 0 {
				fmt.Printf("\tProperties:\n")
			}
			for _, property := range properties {
				fmt.Printf("\t\t%s\n", property)
			}

			fmt.Printf("\tConstructor:\n\t\t%s\n", data.GetConstructor())
			fmt.Printf("\tDestructor:\n\t\t%s\n", data.GetDestructor())

			methods := data.GetMethods()
			if len(methods) > 0 {
				fmt.Printf("\tMethods:\n")
			}
			for _, method := range methods {
				fmt.Printf("\t\t%s\n", method)
			}

			signals := data.Signals()
			if len(signals) > 0 {
				fmt.Printf("\tSignals:\n")
			}
			for _, signal := range data.Signals() {
				fmt.Printf("\t\tsignal %s", signal.Name())

			}
			fmt.Println()
		}
		if ctx.Verbose {
			ast.Print(fset, p)
		}
	}

	return nil
}

func (cmd *listCmd) filter(info os.FileInfo) bool {

	if info.IsDir() {
		return false
	}

	length := len(info.Name())
	if length > 7 && info.Name()[length-7:length-2] == ".gen." {
		return false
	}

	return true
}
