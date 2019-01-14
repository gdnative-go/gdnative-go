// Package gdnative is responsible for parsing and generating binding code for
// Go.
package gdnative

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

// View is a structure that holds the api struct, so it can be used inside
// our temaplte.
type View struct {
	API        API
	StructType string
}

// NotLastElement is a function we use inside the template to test whether or
// not the given element is the last in the slice or not. This is so we can
// correctly insert commas for argument lists.
func (v View) NotLastElement(n int, slice [][]string) bool {
	if n != (len(slice) - 1) {
		return true
	}
	return false
}

// NotVoid checks to see if the return string is void or not. This is used inside
// our template so we can determine if we need to use the `return` keyword in
// the function body.
func (v View) NotVoid(ret string) bool {
	if ret != "void" {
		return true
	}
	return false
}

// HasArgs is a function we use inside the template to test whether or not the
// function has arguments. This is so we can determine if we need to place a
// comma.
func (v View) HasArgs(args [][]string) bool {
	if len(args) != 0 {
		return true
	}
	return false
}

// Generate generates the bindings from the JSON definition
func Generate() {

	// Get the API Path so we can localte the godot api JSON.
	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		panic("$API_PATH is not defined.")
	}
	packagePath := apiPath

	// Create a structure for our template view. This will contain all of
	// the data we need to construct our binding methods.
	var view View

	// Unmarshal the JSON into our struct.
	apis := Parse(packagePath)

	// Add the core API to our view first
	view.API = apis.Core
	view.StructType = "core"

	// Generate the C bindings
	log.Println("Generating", view.StructType, "C headers...")
	WriteTemplate(
		packagePath+"/cmd/generate/templates/gdnative.h.tmpl",
		packagePath+"/gdnative/gdnative.gen.h",
		view,
	)

	log.Println("Generating", view.StructType, "C bindings...")
	WriteTemplate(
		packagePath+"/cmd/generate/templates/gdnative.c.tmpl",
		packagePath+"/gdnative/gdnative.gen.c",
		view,
	)

	// Loop through all of our extensions and generate the bindings for those.
	for _, api := range apis.Extensions {
		view.API = api
		view.StructType = "ext_" + api.Name

		log.Println("Generating", view.StructType, "C headers...")
		WriteTemplate(
			packagePath+"/cmd/generate/templates/gdnative.h.tmpl",
			packagePath+"/gdnative/"+api.Name+".gen.h",
			view,
		)

		log.Println("Generating", view.StructType, "C bindings...")
		WriteTemplate(
			packagePath+"/cmd/generate/templates/gdnative.c.tmpl",
			packagePath+"/gdnative/"+api.Name+".gen.c",
			view,
		)
	}
}

func Parse(packagePath string) APIs {
	// Open the gdnative_api.json file that defines the GDNative API.
	body, err := ioutil.ReadFile(packagePath + "/godot_headers/gdnative_api.json")
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON into our struct.
	var apis APIs
	json.Unmarshal(body, &apis)

	return apis
}

func WriteTemplate(templatePath, outputPath string, view View) {
	// Create a template from our template file.
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// Open the output file for writing
	f, err := os.Create(outputPath)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	// Write the template with the given view.
	err = t.Execute(f, view)
	if err != nil {
		panic(err)
	}
}
