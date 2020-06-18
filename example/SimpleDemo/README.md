# GDNative Simple Example

This example demonstrates how you can use the `gdnative` Go package to create
your own Godot bindings with Go. On this example, we use class auto registration
so the `gdnativego` compiler utility will run automatically to create all the
needed boilerplate in your behalf.

**note**: if you want to see a manual registration example look at the `SimpleDemoManual` example

## Compiling

Dependencies:
 * golang 1.6+
 * golang 1.14 or higher recommended

### CGO (Cross platform)
You can use golang to compile the library if you have it installed, just run make:

    make

You can also run `go build` manually to compile the GDNative library:

    go build -v -buildmode=c-shared -o libsimple.so ./src/simple.go && mv libsimple.so ./bin

You can

## Usage

Create a new object using `load("res://SIMPLE.gdns").new()`

This object has following methods you can use:
 * `get_data()`

