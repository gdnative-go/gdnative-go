package gdnative

import (
	"fmt"
	"log"
	"strings"
)

// instances is a map of our created Godot classes. This will be
// populated when Godot calls the CreateFunc
var instances = map[string]*Class{}

// Method is a NativeScript registrable function
type Method struct {
	name       string
	funcName   string
	attributes *MethodAttributes
	method     *InstanceMethod
}

// Property is a NativeScript registrable class property
type Property struct {
	Value        Variant
	name         string
	propertyName string
	attributes   *PropertyAttributes
	setFunc      *InstancePropertySet
	getFunc      *InstancePropertyGet
}

// Class is a NativeScript registrable class
type Class struct {
	isTool      bool
	name        string
	base        string
	createFunc  *InstanceCreateFunc
	destroyFunc *InstanceDestroyFunc
	methods     []Method
	properties  []Property
}

// RegisterNewGodotClass creates a new ready to go Godot class for us and registers it with in Godot
func RegisterNewGodotClass(isTool bool, name, base string, constructor *InstanceCreateFunc, destructor *InstanceDestroyFunc,
	methods []Method, properties []Property) {

	// create a new Class value and registers it
	godotClass := Class{isTool, name, base, constructor, destructor, methods, properties}
	godotClass.register()
}

// Register registers a Class value with in Godot
func (c *Class) register() {

	// if the class constructor and destructor are not defined create generic ones
	if c.createFunc == nil {
		c.createFunc = c.createGenericConstructor()
	}

	if c.destroyFunc == nil {
		c.destroyFunc = c.createGenericDestructor()
	}

	// we register the class first
	if c.isTool {
		NativeScript.RegisterToolClass(c.name, c.base, c.createFunc, c.destroyFunc)
	} else {
		NativeScript.RegisterClass(c.name, c.base, c.createFunc, c.destroyFunc)
	}

	// then we iterate over every defined method and register them as well
	for _, method := range c.methods {
		method.register()
	}

	// finally we iterate over any defined property and register them
	for _, property := range c.properties {
		property.register()
	}
}

// creates an InstanceCreateFunc value using the given CreateFunc and return it back
func CreateConstructor(className string, fn CreateFunc) InstanceCreateFunc {

	constructor := InstanceCreateFunc{
		CreateFunc: fn,
		MethodData: className,
		FreeFunc:   func(methodData string) {},
	}

	return constructor
}

// creates an InstanceDestroyFunc value using the given DestroyFunc and return it back
func CreateDestructor(className string, fn DestroyFunc) InstanceDestroyFunc {

	destructor := InstanceDestroyFunc{
		DestroyFunc: fn,
		MethodData:  className,
		FreeFunc:    func(methodData string) {},
	}

	return destructor
}

// creates a generic constructor for any given class
func (c *Class) createGenericConstructor() *InstanceCreateFunc {

	constructorFunc := func(object Object, methodData string) string {

		// use the Godot object ID as ID for this class
		id := object.ID()
		Log.Println(fmt.Sprintf("Creating Go generic class %s(%s) constructor with ID %s", c.name, c.base, id))

		// use the class pointer address as the instance ID
		instances[id] = c

		return id
	}

	createFunc := CreateConstructor(c.name, constructorFunc)
	return &createFunc
}

// creates a generic destructor for any given class
func (c *Class) createGenericDestructor() *InstanceDestroyFunc {

	destructorFunc := func(object Object, methodData, userData string) {

		Log.Println(fmt.Sprintf("Destroying %s value with ID: %s", c.name, userData))
		delete(instances, userData)
	}

	destroyFunc := CreateDestructor(c.name, destructorFunc)
	return &destroyFunc
}

// NewGodotMethod creates a new ready to go Godot method for us and return it back
func NewGodotMethod(className, name string, method MethodFunc) Method {

	// create a new Method value
	godotMethod := Method{
		className,
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

	return godotMethod
}

// registers a Method value with in Godot
func (m *Method) register() {

	NativeScript.RegisterMethod(m.name, m.funcName, m.attributes, m.method)
}

// NewGodotProperty creates a new ready to go Godot property, add it to the given class and return it
func NewGodotProperty(class *Class, name, hint, hintString, usage string,
	setFunc *InstancePropertySet, getFunc *InstancePropertyGet) Property {

	// create ok boolean re-usable helper value
	var ok bool

	// create a new PropertyAttributes value and fill it
	var attributes PropertyAttributes
	attributes.HintString = String(hintString)
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
		if attributes.Usage, ok = PropertyUsageFlagsLookupMap[usageKey]; !ok {
			var allowed []string
			for key := range PropertyUsageFlagsLookupMap {
				allowed = append(allowed, strings.Replace(key, "PropertyUsage", "", 1))
			}
			panic(fmt.Sprintf("unknown property usage %s, allowed types: %s", usage, strings.Join(allowed, ", ")))
		}
	} else {
		attributes.Usage = PropertyUsageDefault
	}

	// create a new Property value
	godotProperty := Property{
		NewVariantNil(),
		class.name,
		name,
		&attributes,
		setFunc,
		getFunc,
	}

	return godotProperty
}

// registers a property within the class/godot
func (p *Property) register() error {

	// if set and get functions are not defined generate generic ones
	if p.setFunc == nil {
		p.setFunc = p.createGenericSetter()
	}
	if p.setFunc == nil || p.getFunc == nil {
		return fmt.Errorf("you can not register a property that does not defines both setter and getter functions")
	}

	NativeScript.RegisterProperty(p.name, p.propertyName, p.attributes, p.setFunc, p.getFunc)
	return nil
}

// creates a generic setter method to set property values if none is provided
func (p *Property) createGenericSetter() *InstancePropertySet {

	propertySetter := func(object Object, classProperty, instanceString string, property Variant) {
		Log.Println(fmt.Sprintf("Creating Go generic property setter for %s.%s", p.name, p.propertyName))
		p.Value = property
	}

	instancePropertySet := InstancePropertySet{
		SetFunc:    propertySetter,
		MethodData: fmt.Sprintf("%s::%s", p.name, p.propertyName),
		FreeFunc:   func(methodData string) {},
	}
	return &instancePropertySet
}

// created a generic getter method to get property values if none is provided
func (p *Property) createGenericGetter() *InstancePropertyGet {

	propertyGetter := func(object Object, classProperty, instanceString string) Variant {
		log.Println(fmt.Sprintf("Creating Go generic property getter for %s.%s", p.name, p.propertyName))
		return p.Value
	}

	instancePropertyGet := InstancePropertyGet{
		GetFunc:    propertyGetter,
		MethodData: fmt.Sprintf("%s::%s", p.name, p.propertyName),
		FreeFunc:   func(methodData string) {},
	}
	return &instancePropertyGet
}
