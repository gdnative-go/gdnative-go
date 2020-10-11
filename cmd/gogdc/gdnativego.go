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

type context struct {
	Path    string
	Verbose bool
}

type generateCmd struct{}

type listCmd struct{}

// cli defines our command line structure using Kong
var cli struct {
	Generate generateCmd `cmd help:"Generates autotoregistration boilerplate Go code for user defined structures"` //nolint:govet
	List     listCmd     `cmd help:"List user defined autoregistrable data structures"`                            //nolint:govet
	Version  versionCmd  `cmd help:"Show version information and exit"`

	Path    string `type:"path" default:"." help:"Path where execute the command"`
	Verbose bool   `help:"Verbose output"`
}
