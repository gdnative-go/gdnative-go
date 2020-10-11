![alt text][logo]

# GDNative bindings for Go

Golang bindings to the [Godot GDNative](https://github.com/GodotNativeTools/godot_headers) API.

## Intention

This project is a maintained and a smaller fork of [godot-go](https://github.com/ShadowApex/godot-go)
project, its aim is to make simpler the task of invoke Golang code from [godot](https://github.com/godotengine/godot)'s
GDScript not to write Godot video games using the Golang language. If your are looking for a way to
write Godot games using Golang as development language you can take a look at
[godot-go](https://github.com/ShadowApex/godot-go) project that generates full set of binding for Godot
 classes.

**note**: looks like [godot-go](https://github.com/ShadowApex/godot-go) is unmaintained at the moment
of writing this lines **2020-03-10**

**note**: we are currently working into rewriting the  old godot-go completely from scratch so at some point we will be offering a full set of Godot  Golang wrappers as a separated library

## Motivation

This project exists for three main reasons

1. Bind Golang ecosystem applications to GDScript (like making BBolt or CockroachDB available natively to your GDScript code)
2. Delegate CPU bound calculations to the Golang runtime
3. Use Golang concurrency on NativeScript golang modules from GDScript

If what you are looking for is to write your Godot game in a more performant than GDScript language I
recommend you to use the **godot-mono** version, if C# is not performant enough for you or you want to use a
concurrent native language to write your game code then take a look to [gdnative-rust](https://github.com/GodotNativeTools/godot-rust)
project, it is actively maintained and in good health, you can also take a look to
[godot-cpp](https://github.com/GodotNativeTools/godot-cpp) Godot's C++ bindings.

### Does that means I will never provide full set of bindings in this project?

Yes, that is exactly what it means, this project will only provide bindings for NativeScript

### It is all lost?

No, I am currently working into generate full bindings to Godot full API in a separate project that
uses this one as a base interface between Go and Godot NativeScript. I believe both concepts are
divisible, some times one will want to bring that cool Go library into Godot while some others one
would like to code their whole project using Go, tie this two clearly different needs into the same
library just does not feels right to me

### So is that new library available?

Not yet, will update this repository when its done

## Special Thanks

I would like to give special thanks to the individuals and organizations that had made this project's
existence to be possible 

* Juan Linietsky, Ariel Manzur and the amazing Godot team and Community for their fantastic work
* William Edwards (ShadowApex) for writing the original [https://github.com/ShadowApex/godot-go](godot-go) library this is based on
* Rob Pike, Ken Thompson, Robert Griesemer and the whole Go Language team for creating this language 

## Attributions

This project is based on a previous work by [William Edwards](https://github.com/ShadowApex) and big parts of its
code remains under its license (shown below)

The gdnative-go logo and godopher (godot + gopher) logos are based on the Godot image by 
[Andrea Calabró](https://commons.wikimedia.org/wiki/File:Godot_logo.svg) Licensed under the
[CC-BY 3.0](https://creativecommons.org/licenses/by/3.0/legalcode) license 

The logos are also inspired in a vectorial gopher design by [Takuya Ueda](https://github.com/golang-samples/gopher-vector)
Licensed under the [CC-BY 3.0](https://creativecommons.org/licenses/by/3.0/legalcode) license 

# Licenses

This project diverges from a previous work by William Edwards licensed under the MIT License, both
projects licenses are applicable and thus both of them are shown in this section 

#### GDnative-Go License

```
 Copyright © 2019 - 2020 Oscar Campos <oscar.campos@thepimpam.com>
 Copyright © 2017 - William Edwards

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License
```

#### godot-go License (by ShadowApex)

```
MIT License

Copyright (c) 2017 William Edwards

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

[logo]: img/gdnative-go-logo-with-text.png "GDnative-Go"
