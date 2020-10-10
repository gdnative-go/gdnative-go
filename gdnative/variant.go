package gdnative

/*
#include <gdnative/string.h>
#include <gdnative/variant.h>
#include "gdnative.gen.h"
#include "variant.h"
*/
import "C"

// NewVariantWithString creates a new Variant initialized with the given string
func NewVariantWithString(str String) Variant {
	var variant C.godot_variant
	C.go_godot_variant_new_string(GDNative.api, &variant, str.getBase())

	return Variant{base: &variant}
}

// GetType returns back the VaraintType for this Variant
func (gdt *Variant) GetType() VariantType {
	variantType := C.go_godot_variant_get_type(GDNative.api, gdt.getBase())
	return VariantType(variantType)
}

// VariantArray is a wrapper around Godot C **godot_variant
type VariantArray struct {
	array []Variant
}

func (gdt *VariantArray) getBase() **C.godot_variant {
	variantArray := C.go_godot_variant_build_array(C.int(len(gdt.array)))
	for i, variant := range gdt.array {
		C.go_godot_variant_add_element(variantArray, variant.getBase(), C.int(i))
	}

	return variantArray
}
