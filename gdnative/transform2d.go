package gdnative

/*
#include <gdnative/gdnative.h>
#include "gdnative.gen.h"
*/
import "C"

// NewTransform2D godot_transform2d_new [[godot_transform2d * r_dest] [const godot_real p_rot] [const godot_vector2 * p_pos]] void
func NewTransform2D(rot Real, pos Vector2) *Transform2D {
	var dest C.godot_transform2d
	arg1 := rot.getBase()
	arg2 := pos.getBase()
	C.go_godot_transform2d_new(GDNative.api, &dest, arg1, arg2)
	return &Transform2D{base: &dest}
}

// NewTransform2DAxisOrigin godot_transform2d_new_axis_origin [[godot_transform2d * r_dest] [const godot_vector2 * p_x_axis] [const godot_vector2 * p_y_axis] [const godot_vector2 * p_origin]] void
func NewTransform2DAxisOrigin(xAxis Vector2, yAxis Vector2, origin Vector2) *Transform2D {
	var dest C.godot_transform2d
	arg1 := xAxis.getBase()
	arg2 := yAxis.getBase()
	arg3 := origin.getBase()
	C.go_godot_transform2d_new_axis_origin(GDNative.api, &dest, arg1, arg2, arg3)
	return &Transform2D{base: &dest}
}

// NewTransform2DIdentity godot_transform2d_new_identity [[godot_transform2d * r_dest]] void
func NewTransform2DIdentity() *Transform2D {
	var dest C.godot_transform2d
	C.go_godot_transform2d_new_identity(GDNative.api, &dest)
	return &Transform2D{base: &dest}
}
