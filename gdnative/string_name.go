package gdnative

/*
#include <gdnative/string.h>
#include <gdnative/string_name.h>
#include "gdnative.gen.h"
*/
import "C"

// NewStringName returns a StringName wrapper over a Godot C godot_string_name
// initialized with the contents of the given name
func NewStringName(name String) *StringName {
	var dest C.godot_string_name
	arg1 := name.getBase()
	C.go_godot_string_name_new(GDNative.api, &dest, arg1)
	return &StringName{base: &dest}
}

// NewStringNameData returns a StringName wrapper over a Godot C godot_string_name
// initialized with the given name
func NewStringNameData(name Char) *StringName {
	var dest C.godot_string_name
	arg1 := name.getBase()
	C.go_godot_string_name_new_data(GDNative.api, &dest, arg1)
	return &StringName{base: &dest}
}
