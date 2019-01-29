package gdnative

import (
	"fmt"
	"strings"
)

// Method is a NativeScript registrable function
type Method struct {
	name       string
	funcName   string
	attributes *MethodAttributes
	method     *InstanceMethod
}

// Property is a NativeScript registrable class property
type Property struct {
	name       string
	path       string
	attributes *PropertyAttributes
	setFunc    *InstancePropertySet
	getFunc    *InstancePropertyGet
}

// Class is a NativeScript registrable class
type Class struct {
	isTool      bool
	name        string
	base        string
	createFunc  *InstanceCreateFunc
	destroyFunc *InstanceDestroyFunc
	methods     []*Method
	properties  []*Property
}

// NewGodotClass creates a new ready to go Godot class for us and registers it with in Godot
func NewGodotClass(isTool bool, name, base string, constructor *InstanceCreateFunc, destructor *InstanceDestroyFunc,
	methods []*Method, properties []*Property) *Class {

	// create a new Class value and registers it
	godotClass := Class{isTool, name, base, constructor, destructor, methods, properties}
	godotClass.register()

	// iterate over all the methods on this class and register them
	for _, method := range godotClass.methods {
		method.register()
	}

	return &godotClass
}

// regiusters a Class value with in Godot
func (c *Class) register() {

	if c.isTool {
		NativeScript.RegisterToolClass(c.name, c.base, c.createFunc, c.destroyFunc)
	} else {
		NativeScript.RegisterClass(c.name, c.base, c.createFunc, c.destroyFunc)
	}
}

// NewGodotMethod creates a new ready to go Godot method for us and return it back
func NewGodotMethod(class *Class, name string, method MethodFunc) *Method {

	// create a new Method value
	godotMethod := Method{
		class.name,
		name,
		&MethodAttributes{
			RPCType: MethodRpcModeDisabled,
		},
		&InstanceMethod{
			Method:     method,
			MethodData: name,
			FreeFunc:   func(methodData string) {},
		},
	}

	return &godotMethod
}

// registers a Method value with in Godot
func (m *Method) register() {

	NativeScript.RegisterMethod(m.name, m.funcName, m.attributes, m.method)
}

// NewGodotProperty creates a new ready to go Godot property, add it to the given class and return it
func NewGodotProperty(class *Class, name, path, hint, hintString, usage string,
	setFunc *InstancePropertySet, getFunc *InstancePropertyGet) *Property {

	// create ok boolean re-usable helper value
	var ok bool

	// create a new PropertyAttributes value and fill it
	var attributes PropertyAttributes
	attributes.HintString = hintString
	attributes.DefaultValue = NewVariantNil()
	if hint != "" {
		hintKey := fmt.Sprintf("PropertyHint%s", hint)
		if attributes.Hint, ok = PropertyHintLookupMap[hintKey]; !ok {
			var allowed []string
			for key := range PropertyHintLookupMap {
				allowed = append(allowed, strings.Replace(key, "PropertyHint", "", 1))
			}
			panic(fmt.Sprintf("unknown property hint %s, allowed types: %s", hint, strings.Join(allowed, ", ")))
		}
	} else {
		attributes.Hint = PropertyHintNone
	}

	if usage != "" {
		usageKey := fmt.Sprintf("PropertyUsage%s", usage)
		if attributes.Usage, ok = PropertyUsageLookupMap[usageKey]; !ok {
			var allowed []string
			for key := range PropertyUsageLookupMap {
				allowed = append(allowed, strings.Replace(key, "PropertyUsage", "", 1))
			}
			panic(fmt.Sprintf("unknnown property usage %s, allowed types: %s", usage, strings.Join(allowed, ", ")))
		}
	} else {
		attributes.Usage = PropertyUsageDefault
	}

	// create a new Property value
	godotProperty := Property{
		name,
		path,
		attributes,
		setFunc,
		getFunc,
	}

	return &godotProperty
}
