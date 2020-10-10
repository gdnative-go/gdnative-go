#include "variant.h"
#include <gdnative/variant.h>
#include <stdlib.h>

godot_variant **go_godot_variant_build_array(int length) {
	godot_variant **arr = malloc(sizeof(godot_variant *) * length);

	return arr;
}

void go_godot_variant_add_element(godot_variant **array, godot_variant *element,
				  int index) {
	godot_variant copy = *element;
	array[index] = &copy;
}

godot_variant *go_godot_new_variant() {
	godot_variant *var = malloc(sizeof(godot_variant));

	return var;
}
