// Copyright © 2019 - 2020 Oscar Campos <oscar.campos@thepimpam.com>
// Copyright © 2017 - William Edwards
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gdnative

import (
	"fmt"
	"strings"
	"unicode"
)

// Registrable is the interface external code communicates with registryClass
type Registrable interface {
	GetBase() string
	GetConstructor() string
	GetDestructor() string
	GetMethods() []string
	GetProperties() []string
	SetConstructor(*registryConstructor)
	SetDestructor(*registryDestructor)
	Constructor() string
	Destructor() string
	AddMethod(*registryMethod)
	AddMethods([]*registryMethod)
	Methods() []*registryMethod
}

type registryClass struct {
	base        string
	constructor *registryConstructor
	destructor  *registryDestructor
	methods     []*registryMethod
	properties  []*registryProperty
	signals     []*registrySignal
}

// GetBase returns back the Godot base class for this type as a string
func (rc *registryClass) GetBase() string {
	return rc.base
}

// GetConstructor returns back this type constructor as a string
func (rc *registryClass) GetConstructor() string {

	if rc.constructor != nil {
		return fmt.Sprintf("func %s() *%s", rc.constructor.customFunc, rc.constructor.class)
	}

	return ""
}

// SetConstructor sets this type constructor
func (rc *registryClass) SetConstructor(constructor *registryConstructor) {

	if constructor != nil {
		rc.constructor = constructor
	}
}

// HasConstructor returns true if this type has custom constructor
func (rc *registryClass) HasConstructor() bool {
	return rc.constructor != nil
}

// Constructor returns the constructor custom function
func (rc *registryClass) Constructor() string {

	if rc.constructor != nil {
		return rc.constructor.customFunc
	}

	return ""
}

// GetDestructor returns back this type destructor as a string
func (rc *registryClass) GetDestructor() string {

	if rc.destructor != nil {
		return fmt.Sprintf("func %s()", rc.destructor.customFunc)
	}

	return ""
}

// SetDestructor sets this type destructor
func (rc *registryClass) SetDestructor(destructor *registryDestructor) {

	if destructor != nil {
		rc.destructor = destructor
	}
}

// HasDestructor returns true if this type has custom destructor
func (rc *registryClass) HasDestructor() bool {
	return rc.destructor != nil
}

// Destructor returns the destructor custom function
func (rc *registryClass) Destructor() string {

	if rc.destructor != nil {
		return rc.destructor.customFunc
	}

	return ""
}

// GetMethods returns back this type exported methods as strings
func (rc *registryClass) GetMethods() []string {

	methods := []string{}
	for _, method := range rc.methods {
		methods = append(methods, fmt.Sprintf("func %s(%s) %s", method.name, method.GetParams(), method.GetReturnValues()))
	}

	return methods
}

// Methods returns back this type list of registryMethods
func (rc *registryClass) Methods() []*registryMethod {
	return rc.methods
}

// AddMethods adds a list of methods for this type
func (rc *registryClass) AddMethods(methods []*registryMethod) {

	for i := range methods {
		rc.AddMethod(methods[i])
	}
}

// AddMethod adds a method to this type
func (rc *registryClass) AddMethod(method *registryMethod) {

	if method != nil {
		rc.methods = append(rc.methods, method)
	}
}

// GetProperties returns back this class properties as strings
func (rc *registryClass) GetProperties() []string {

	properties := []string{}
	for _, property := range rc.properties {
		properties = append(properties, fmt.Sprintf("%s %s %s", property.name, property.kind, property.tag))
	}

	return properties
}

type registryConstructor struct {
	class, customFunc string
}

type registryDestructor struct {
	class, customFunc string
}

type registryMethod struct {
	class, name  string
	params       []*registryMethodParam
	returnValues []*registryMethodReturnValue
}

func (rm *registryMethod) GetName() string {
	return rm.name
}

func (rm *registryMethod) GodotName() string {

	if rm.name[0] == 'V' && unicode.IsUpper(rune(rm.name[1])) {
		return strings.ToLower(fmt.Sprintf("_%s", rm.name[1:]))
	}

	return rm.name
}

// GetParams returns this type method params as a string
func (rm *registryMethod) GetParams() string {

	pairs := []string{}
	for _, param := range rm.params {
		pairs = append(pairs, fmt.Sprintf("%s %s", param.name, param.kind))
	}

	return strings.Join(pairs, ", ")
}

// GetReturnValues returns this type method return values as a string
func (rm *registryMethod) GetReturnValues() string {

	values := []string{}
	for _, value := range rm.returnValues {
		values = append(values, value.kind)
	}
	return strings.Join(values, ", ")
}

// Arguments returns a slice of this method arguments structures
func (rm *registryMethod) Arguments() []*registryMethodParam {
	return rm.params
}

// HasReturns returns true if this method has return values, otherwise returns false
func (rm *registryMethod) HasReturns() bool {
	return len(rm.returnValues) > 0
}

// FunctionCallWithParams returns a string representing how this method should be called
func (rm *registryMethod) FunctionCallWithParams() string {

	arguments := make([]string, len(rm.params))
	for i, arg := range rm.params {
		arguments[i] = arg.name
	}
	return fmt.Sprintf("%s(%s)", rm.name, strings.Join(arguments, ", "))
}

// NewVariantType returns the right NewVariant<Type> method from gdnative for our return type
func (rm *registryMethod) NewVariantType() string {

	variant := "NewVariant%s(gdnative.%s(value))"
	conversions := map[string]string{
		"bool":    fmt.Sprintf(variant, "Bool", "Bool"),
		"uint":    fmt.Sprintf(variant, "Uint", "Uint64T"),
		"int":     fmt.Sprintf(variant, "Int", "Int64T"),
		"float64": fmt.Sprintf(variant, "Real", "Double"),
		"string":  fmt.Sprintf(variant, "String", "String"),
	}

	retLength := len(rm.returnValues)
	if retLength == 1 {
		result, ok := conversions[rm.returnValues[0].kind]
		if !ok {
			result = "value"
		}
		return result
	}

	if retLength >= 2 && retLength <= 3 {
		valid := make([]bool, retLength)
		for i, val := range rm.returnValues {
			switch val.kind {
			case "float32", "float64", "gdnative.Double", "gdnative.Real":
				valid[i] = true
			}
		}
		allValid := true
		for _, b := range valid {
			if !b {
				allValid = false
				break
			}
		}

		if allValid {
			return fmt.Sprintf(
				variant,
				fmt.Sprintf("Vector%d", retLength),
				fmt.Sprintf("NewVector%d", retLength),
			)
		}
	}

	return "value"
}

type registryProperty struct {
	name, kind, tag string
}

// Name returns the name of the property back
func (rp *registryProperty) Name() string {
	return rp.name
}

type registrySignal struct{}

type registryMethodParam struct {
	name, kind string
}

// Name returns this param name
func (rmp *registryMethodParam) Name() string {
	return rmp.name
}

// Kind returns this param kind
func (rmp *registryMethodParam) Kind() string {
	return rmp.kind
}

// ConvertFunction returns the right GDNative 'As<Type>' function for this param kind as a string
func (rmp *registryMethodParam) ConvertFunction() string {

	conversions := map[string]string{
		"bool":              "AsBool()",
		"uint":              "AsUint()",
		"int":               "AsInt()",
		"float64":           "AsReal()",
		"string":            "AsString()",
		"vector2":           "AsVector2()",
		"vector3":           "AsVector3()",
		"rect2":             "AsRect2()",
		"transform2d":       "AsTransform2D()",
		"plane":             "AsPlane()",
		"quat":              "AsQuat()",
		"aabb":              "AsAabb()",
		"basis":             "AsBasis()",
		"transform":         "AsTransform()",
		"color":             "AsColor()",
		"nodepath":          "AsNodePath()",
		"rid":               "AsRid()",
		"object":            "AsObject()",
		"dictionary":        "AsDictionary()",
		"arraytype":         "AsArray()",
		"arraytype_byte":    "AsPoolByteArray()",
		"arraytype_int":     "AsPoolIntArray()",
		"arraytype_float":   "AsPoolRealArray()",
		"arraytype_string":  "AsPoolStringArray()",
		"arraytype_vector2": "AsPoolVector2Array()",
		"arraytype_vector3": "AsPoolVector3Array()",
		"arraytype_color":   "AsPoolColorArray()",
	}

	value, ok := conversions[rmp.kind]
	if ok {
		return value
	}

	switch rmp.kind {
	case "float32":
		value = conversions["float64"]
	case "int8", "int16", "int32", "int64", "byte":
		value = conversions["int"]
	case "uint8", "uint16", "uint32", "uint64":
		value = conversions["uint"]
	case "gdnative.Float", "gdnative.Double":
		value = conversions["float64"]
	case "gdnative.Int64T", "gdnative.SignedChar":
		value = conversions["int"]
	case "gdnative.Uint", "gdnative.Uint8T", "gdnative.Uint32T", "gdnative.Uint64T":
		value = conversions["uint"]
	case "gdnative.String", "gdnative.Char", "gdnative.WcharT":
		value = conversions["string"]
	default:
		if strings.HasPrefix(rmp.kind, "gdnative.") {
			value = rmp.kind
		}

		if strings.HasPrefix(rmp.kind, "godot.") || strings.HasPrefix(rmp.kind, "*godot.") {
			value = conversions["object"]
		}

		if strings.Contains(rmp.kind, "ArrayType") {
			value = conversions[parseArrayType(rmp.kind)]
		}

		if strings.Contains(rmp.kind, "MapType") {
			value = conversions["dictionary"]
		}
	}

	if value == "" {
		value = conversions["object"]
	}

	return value
}

type registryMethodReturnValue struct {
	kind string
}

func parseArrayType(array string) string {

	var arrayType string
	var openBracket bool
	var openBracketCount int
	for i := 0; i < len(array); i++ {
		if openBracket {
			if array[i] == ']' {
				openBracketCount--
				if openBracketCount == 0 {
					openBracket = false
					continue
				}
			}

			if array[i] == '[' {
				openBracketCount++
			}
			arrayType += string(array[i])
			continue
		}

		if array[i] == '[' {
			openBracket = true
			openBracketCount = 1
		}
	}

	return fmt.Sprintf("arraytype_%s", strings.ReplaceAll(arrayType, "gdnative.", ""))
}
