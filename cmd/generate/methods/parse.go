// Package methods is a package that parses the GDNative headers for type definitions
// of methods
package methods

import (
	"os"

	"gitlab.com/pimpam-games-studio/gdnative-go/cmd/generate/gdnative"
)

// Parse will parse the GDNative headers. Takes a list of headers/structs to ignore.
// Definitions in the given headers and definitions
// with the given name will not be added to the returned list of type definitions.
// We'll need to manually create these structures.
func Parse() gdnative.APIs {

	// Get the API Path so we can localte the godot api JSON.
	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		panic("$API_PATH is not defined.")
	}
	packagePath := apiPath

	// Parse the GDNative JSON for method data.
	apis := gdnative.Parse(packagePath)
	return apis
}
