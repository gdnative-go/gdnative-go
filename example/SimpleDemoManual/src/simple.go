package main

import (
	"fmt"

	"gitlab.com/pimpam-games-studio/gdnative-go/gdnative"
)

// SimpleClass is a structure that we can register with Godot.
type SimpleClass struct {
	base gdnative.Object
}

// Instances is a map of our created Godot classes. This will be populated when
// Godot calls the CreateFunc.
var Instances = map[string]*SimpleClass{}

// NativeScriptInit will run on NativeScript initialization. It is responsible
// for registering all our classes with Godot.
func nativeScriptInit() {
	gdnative.Log.Warning("Initializing nativescript from Go!")

	// Define an instance creation function. This will be called when Godot
	// creates a new instance of our class.
	createFunc := gdnative.InstanceCreateFunc{
		CreateFunc: simpleConstructor,		// reference to the constructor function
		MethodData: "SIMPLE",			// name of the class the constructor is attached to
		FreeFunc:   func(methodData string) {}, // function for freeing resources (usually empty)
	}

	// Define an instance destroy function. This will be called when Godot
	// asks our library to destroy our class instance.
	destroyFunc := gdnative.InstanceDestroyFunc{
		DestroyFunc: simpleDestructor,		 // reference to the destructor function
		MethodData:  "SIMPLE",			 // name of the class the destructor is attached to
		FreeFunc:    func(methodData string) {}, // function for freeing resources (usually empty)
	}

	// Register our class with Godot.
	gdnative.Log.Warning("Registering SIMPLE class...")
	gdnative.NativeScript.RegisterClass(
		"SIMPLE",     // the name of the class we are registering
		"Reference",  // class from which this class inherits from
		&createFunc,  // class constructor
		&destroyFunc, // class destructor
	)

	// Register a method with Godot.
	gdnative.Log.Warning("Registering SIMPLE method...")
	gdnative.NativeScript.RegisterMethod(
		"SIMPLE",   // the name of the class we are registering the method within
		"get_data", // the visible name for the method inside Godot
		&gdnative.MethodAttributes{ // Method RPC type, this will typically be Disabled unless RPC is required
			RPCType: gdnative.MethodRpcModeDisabled,
		},
		&gdnative.InstanceMethod{ // method wrapper
			Method:     simpleMethod,		// the simpleMethod function reference that implements our logic
			MethodData: "SIMPLE",			// method name as will be used with in Godot
			FreeFunc:   func(methodData string) {}, // function for freeing resources (usually empty)
		},
	)
}

func simpleConstructor(object gdnative.Object, methodData string) string {
	gdnative.Log.Println("Creating new SimpleClass...")

	// Create a new instance of our struct.
	instance := &SimpleClass{
		base: object,
	}

	// Use the pointer address as the instance ID
	instanceID := fmt.Sprintf("%p", instance)
	Instances[instanceID] = instance

	// Return the instanceID
	return instanceID
}

func simpleDestructor(object gdnative.Object, methodData, userData string) {
	gdnative.Log.Println("Destroying SimpleClass with ID:", userData, "...")
	// Delete the instance from our map of instances
	delete(Instances, userData)
}

func simpleMethod(object gdnative.Object, methodData, userData string, numArgs int, args []gdnative.Variant) gdnative.Variant {
	gdnative.Log.Println("SIMPLE.get_data() called!")

	data := gdnative.NewStringWithWideString("World from godot-go from instance: " + object.ID() + "!")
	ret := gdnative.NewVariantWithString(data)

	return ret
}

// The "init()" function is a special Go function that will be called when this library
// is initialized. Here we can register our Godot classes.
func init() {
	// Set the initialization script that will run upon NativeScript initialization.
	// This function will handle using the NativeScript API to register all of our
	// classes.
	gdnative.SetNativeScriptInit(nativeScriptInit)
}

// This never gets called, but it is necessary to export as a shared library.
func main() {
}
