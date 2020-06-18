package main

import (
	"gitlab.com/pimpam-games-studio/gdnative-go/gdnative"
)

// Ignored is ignored by the gdnativego compiler
type Ignored struct{}

// SimpleClass is a structure that we can register with Godot.
//godot::register as SIMPLE
type SimpleClass struct {
	Hit         gdnative.Signal
	HP          gdnative.Int   `hint:"range" hint_string:"The player's Hit Points" usage:"Default"`
	Mana, Blood gdnative.Int   `hint:"range" hint_string:"The player points to cast spells"`
	IgnoreMe    gdnative.Float `-` // this property will be ignored
}

// New creates a new SimpleClass value and returns a pointer to it
//godot::constructor(SimpleClass)
func New() *SimpleClass {

	sc := SimpleClass{
		// Signals must be defined as literals or they will be ignored by the gdnativego compiler
		Hit: gdnative.Signal{
			Name:           "hit",
			NumArgs:        gdnative.Int(1),
			NumDefaultArgs: gdnative.Int(1),
			Args: []gdnative.SignalArgument{
				{
					Name:         gdnative.String("power"),
					Type:         gdnative.Int(gdnative.VariantTypeInt),
					Hint:         gdnative.PropertyHintRange,
					HintString:   "Hit power value",
					Usage:        gdnative.PropertyUsageDefault,
					DefaultValue: gdnative.NewVariantInt(gdnative.Int64T(0)),
				},
			},
			DefaultArgs: []gdnative.Variant{
				gdnative.NewVariantInt(gdnative.Int64T(0)),
			},
		},
	}
	return &sc
}

// GetData is automatically registered to SimpleClass on Godot
//godot::export as get_data
func (sc *SimpleClass) GetData() gdnative.Variant {

	gdnative.Log.Println("SIMPLE.get_data() called!")

	data := gdnative.NewStringWithWideString("Hello World from gdnative-go instance!")
	return gdnative.NewVariantWithString(data)
}

// This never gets called, but it necessary to export as a shared library.
func main() {
}
