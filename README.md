# GDNative bindings for Go

Golang bindings to the [Godot GDNative](https://github.com/GodotNativeTools/godot_headers) API.

## Intention

This project is a maintained as a smaller fork of [godot-go](https://github.com/ShadowApex/godot-go)
project, its aim is to make simpler the task of invoke Golang code from [godot](https://github.com/godotengine/godot)'s
GDScript not to write Godot video games using the Golang language. If your are looking for a way to
write Godot games using Golang as development language you can take a look at
[godot-go](https://github.com/ShadowApex/godot-go) project that generates full set of binding for Godot
 classes.

**note**: looks like [godot-go](https://github.com/ShadowApex/godot-go) is unmaintained at the moment
of writing this lines **2020-03-10**

## Motivation

This project exists for three main reasons

1. Bind Golang ecosystem applications to GDScript (like making BBolt or CockroachDB available natively to your GDScript code)
2. Delegate CPU bound calculations to the Golang runtime
3. Use Golang concurrency on GDNative golang modules from GDScript

If what you are looking for is to write your Godot game in a more performant than GDScript language I
recommend you to use the **godot-mono** version, if C# is not performant enough for you or you want to use a
concurrent native language to write your game code then take a look to [gdnative-rust](https://github.com/GodotNativeTools/godot-rust)
project, it is actively maintained and in good health, you can also take a look to
[godot-cpp](https://github.com/GodotNativeTools/godot-cpp) Godot's C++ bindings.
