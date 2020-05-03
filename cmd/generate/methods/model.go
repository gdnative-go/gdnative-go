package methods

// Method represents a Go method
type Method struct {
	Arguments      []Argument
	GoName         string
	GoReturns      string
	Name           string
	ParentStruct   string
	Returns        string
	ReturnsPointer bool
}

// Argument represents a Go method argument
type Argument struct {
	GoName    string
	GoType    string
	IsPointer bool
	Name      string
	Type      string
}
