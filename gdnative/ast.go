package gdnative

import (
	"fmt"
	"go/ast"
	"strings"
)

const (
	godotRegister    string = "godot::register"
	godotConstructor string = "godot::constructor"
	godotDestructor  string = "godot::destructor"
)

func LookupRegistrableTypeDeclarations(pkg *ast.Package) map[string]Registrable {

	var classes = make(map[string]Registrable)

	// make a first iteration to capture all registrable classes and their properties
	for _, file := range pkg.Files {

		for _, node := range file.Decls {
			gd, ok := node.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, d := range gd.Specs {
				tp, ok := d.(*ast.TypeSpec)
				if !ok {
					continue
				}

				sp, ok := tp.Type.(*ast.StructType)
				if !ok {
					continue
				}

				if gd.Doc != nil {
					for _, line := range gd.Doc.List {
						docstring := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(line.Text, "/", "")))
						if strings.HasPrefix(docstring, godotRegister) {

							className := getClassName(tp)
							class := registryClass{
								base:       getBaseClassName(sp),
								properties: lookupProperties(sp),
							}
							classes[className] = &class
						}
					}
				}
			}
		}
	}

	// make a second iteration to look for class methods and signals
	for className := range classes {

		class := classes[className]
		for _, file := range pkg.Files {
			class.SetConstructor(lookupInstanceCreateFunc(className, file))
			class.SetDestructor(lookupInstanceDestroyFunc(className, file))
			class.AddMethods(lookupMethods(className, file))
			// signals:     lookupSignals(className, file))
		}
	}

	return classes
}

// getClassName extracts and build the right class name for the registry
func getClassName(tp *ast.TypeSpec) string {

	className := tp.Name.String()
	return className
}

// getBaseClassName extracts the base class name for the registry
func getBaseClassName(sp *ast.StructType) string {

	// TODO: need to make this way smarter to look for parent types of this type
	var baseClassName string
	for i := 0; i < sp.Fields.NumFields(); i++ {
		expr, ok := sp.Fields.List[i].Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		ident, ok := expr.X.(*ast.Ident)
		if !ok {
			continue
		}

		if ident.Name != "godot" {
			continue
		}

		baseClassName = fmt.Sprintf("godot.%s", expr.Sel.Name)
		break
	}

	return baseClassName
}

// lookupInstanceCreateFunc extract the "constructor" for the given type or create a default one
func lookupInstanceCreateFunc(className string, file *ast.File) *registryConstructor {

	for _, node := range file.Decls {
		fd, ok := node.(*ast.FuncDecl)
		if !ok || fd.Doc == nil {
			continue
		}

		for _, line := range fd.Doc.List {
			docstring := strings.TrimSpace(strings.ReplaceAll(line.Text, "/", ""))
			if strings.HasPrefix(strings.ToLower(docstring), godotConstructor) {
				structName := strings.TrimSpace(docstring[len(godotConstructor):])

				// get rid of parenthesis
				if strings.HasPrefix(structName, "(") && strings.HasSuffix(structName, ")") {
					// make sure this is the only parenthesis structName
					if strings.Count(structName, "(") > 1 || strings.Count(structName, ")") > 1 {
						// this is a syntax error
						panic(fmt.Errorf("could not parse constructor comment %s, many parenthesis", docstring))
					}

					structName = structName[1 : len(structName)-1]
					if structName != className {
						// this constructor doesn't match with our class, skip it
						continue
					}

					constructor, err := validateConstructor(structName, fd)
					if err != nil {
						panic(err)
					}

					return constructor
				}
			}
		}
	}

	// if we are here it means the user didn't specified a custom constructor
	return nil
}

// validateConstructor returns an error if the given constructor is not valid, it returns nil otherwise
func validateConstructor(structName string, fd *ast.FuncDecl) (*registryConstructor, error) {

	funcName := fd.Name.String()
	if fd.Recv != nil {
		value := "UnknownType"
		switch t := fd.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			value = t.X.(*ast.Ident).Name
		case *ast.Ident:
			value = t.Name
		}

		return nil, fmt.Errorf("%s is a method of %s type it can not be used as constructor", funcName, value)
	}

	if fd.Type.Params.List != nil {
		return nil, fmt.Errorf("constructors of %s values take no params but %s takes %d", structName, funcName, fd.Type.Params.NumFields())
	}

	if fd.Type.Results == nil || fd.Type.Results.List == nil {
		return nil, fmt.Errorf("constructors of %s values have to return a pointer to *%s but %s return nothing", structName, structName, funcName)
	}

	if fd.Type.Results.NumFields() > 1 {
		return nil, fmt.Errorf("constructors of %s values must return exactly one value but %s returns %d", structName, funcName, fd.Type.Results.NumFields())
	}

	switch t := fd.Type.Results.List[0].Type.(type) {
	case *ast.StarExpr:
		if t.X.(*ast.Ident).Name != structName {
			return nil, fmt.Errorf(
				"constructors of %s values must return a pointer to *%s but %s returns a pointer to %s instead",
				structName, structName, funcName, t.X.(*ast.Ident).Name,
			)
		}
		break
	default:
		return nil, fmt.Errorf("constructors of %s values must return a pointer to *%s but %s returns %v", structName, structName, funcName, t)
	}

	constructor := registryConstructor{
		class:      structName,
		customFunc: fd.Name.String(),
	}
	return &constructor, nil
}

// lookupInstanceDestroyFunc
func lookupInstanceDestroyFunc(className string, file *ast.File) *registryDestructor {

	for _, node := range file.Decls {
		fd, ok := node.(*ast.FuncDecl)
		if !ok || fd.Doc == nil {
			continue
		}

		for _, line := range fd.Doc.List {
			docstring := strings.TrimSpace(strings.ReplaceAll(line.Text, "/", ""))
			if strings.HasPrefix(strings.ToLower(docstring), godotDestructor) {
				structName := strings.TrimSpace(docstring[len(godotDestructor):])

				// get rid of parenthesis
				if strings.HasPrefix(structName, "(") && strings.HasSuffix(structName, ")") {
					// make sure this is the only parenthesis structName
					if strings.Count(structName, "(") > 1 || strings.Count(structName, ")") > 1 {
						// this is a syntax error
						panic(fmt.Errorf("could not parse destructor comment %s, many parenthesis", docstring))
					}

					structName = structName[1 : len(structName)-1]
					if structName != className {
						// this destructor doesn't match with our class, skip it
						continue
					}

					destructor, err := validateDestructor(structName, fd)
					if err != nil {
						panic(err)
					}

					return destructor
				}
			}
		}
	}

	// if we are here that means the user didn't provide a custom destructor function
	return nil
}

// validateDestructor returns an error if the given constructor is not valid, it returns nil otherwise
func validateDestructor(structName string, fd *ast.FuncDecl) (*registryDestructor, error) {

	funcName := fd.Name.String()
	if fd.Recv != nil {
		value := "UnknownType"
		switch t := fd.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			value = t.X.(*ast.Ident).Name
		case *ast.Ident:
			value = t.Name
		}

		return nil, fmt.Errorf("%s is a method of %s type it can not be used as destructor", funcName, value)
	}

	if fd.Type.Params.List != nil {
		return nil, fmt.Errorf("destructors of %s values take no params but %s takes %d", structName, funcName, fd.Type.Params.NumFields())
	}

	if fd.Type.Results != nil || fd.Type.Results.List != nil || len(fd.Type.Results.List) > 0 {
		return nil, fmt.Errorf("destructors of %s values have to return a pointer to *%s but %s return nothing", structName, structName, funcName)
	}

	destructor := registryDestructor{
		class:      structName,
		customFunc: funcName,
	}
	return &destructor, nil
}

// lookupMethods look up for every exported method that is owned by the type
// and fill a registration data structure with it
func lookupMethods(className string, file *ast.File) []*registryMethod {

	methods := []*registryMethod{}
	for _, node := range file.Decls {
		fd, ok := node.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// check for init function, if present fail and complain
		if fd.Name.String() == "init" {
			panic(fmt.Errorf(
				"init function present on files, you can not provide your own init function while autoregistering classes",
			))
		}

		// ignore non exported methods
		if !fd.Name.IsExported() {
			continue
		}

		// ignore non methods
		if fd.Recv == nil {
			continue
		}

		// ignore methods from other types
		switch t := fd.Recv.List[0].Type.(type) {
		case *ast.StarExpr:
			if t.X.(*ast.Ident).Name != className {
				continue
			}
		case *ast.Ident:
			if t.Name != className {
				continue
			}
		default:
			// this should not be possible but just in case
			continue
		}

		funcName := fd.Name.String()
		method := registryMethod{
			name:         funcName,
			class:        className,
			params:       lookupParams(fd.Type.Params),
			returnValues: lookupReturnValues(fd),
		}
		methods = append(methods, &method)
	}

	return methods
}

func lookupParams(fields *ast.FieldList) []*registryMethodParam {

	params := []*registryMethodParam{}
	if fields.NumFields() > 0 {
		for _, field := range fields.List {

			kind := parseDefault(field.Type, "")
			if kind == "" {
				// if we don't have a type skip iteration
				continue
			}

			for _, name := range field.Names {
				params = append(params, &registryMethodParam{
					name: name.String(),
					kind: kind,
				})
			}
		}
	}

	return params
}

func lookupReturnValues(fd *ast.FuncDecl) []*registryMethodReturnValue {

	returnValues := []*registryMethodReturnValue{}
	if fd.Type.Results.NumFields() > 0 {
		for _, result := range fd.Type.Results.List {
			value := registryMethodReturnValue{
				kind: parseDefault(result.Type, "gdnative.Variant"),
			}
			returnValues = append(returnValues, &value)
		}
	}

	return returnValues
}

func lookupProperties(sp *ast.StructType) []*registryProperty {

	properties := []*registryProperty{}
	if sp.Fields.NumFields() > 0 {
		for _, field := range sp.Fields.List {
			if field.Names == nil {
				// this field is a selector expression, skip it
				continue
			}

			kind := parseDefault(field.Type, "")
			if kind == "" {
				// if we don't have a type skip iteration
				continue
			}

			if field.Tag != nil {
				if len(field.Names) > 1 {
					// this is a syntactic error
					panic(fmt.Errorf(
						"struct field tag repeated for fields %v, please put each field in a different line if you want to use field tags", field.Names,
					))
				}

				if !field.Names[0].IsExported() {
					continue
				}

				property := registryProperty{
					name: field.Names[0].Name,
					kind: kind,
					tag:  field.Tag.Value,
				}
				properties = append(properties, &property)
				continue
			}

			for _, name := range field.Names {
				if !name.IsExported() {
					continue
				}

				properties = append(properties, &registryProperty{
					name: name.String(),
					kind: kind,
				})
			}
		}
	}

	return properties
}

func parseDefault(expr ast.Expr, def string) string {

	kind := def
	switch t := expr.(type) {
	case *ast.Ident:
		kind = t.Name
	case *ast.StarExpr:
		kind = parseStarExpr(t)
	case *ast.ParenExpr:
		kind = parseParenExpr(t)
	case *ast.ArrayType:
		kind = parseArray(t)
	case *ast.MapType:
		kind = parseMap(t)
	case *ast.SelectorExpr:
		kind = fmt.Sprintf("%s.%s", parseDefault(t.X, def), t.Sel.String())
	}

	return kind
}

func parseStarExpr(expr *ast.StarExpr) string {

	kind := "*%s"
	return fmt.Sprintf(kind, parseDefault(expr.X, "gdnatve.Pointer"))
}

func parseParenExpr(field *ast.ParenExpr) string {

	return parseDefault(field.X, "gdnative.Pointer")
}

func parseArray(field *ast.ArrayType) string {

	result := "ArrayType[%s]"
	return fmt.Sprintf(result, parseDefault(field.Elt, "ArrayType[]"))
}

func parseMap(field *ast.MapType) string {

	result := "MapType[%s]%s"

	key := parseDefault(field.Key, "gdnative.Pointer")
	value := parseDefault(field.Value, "gdnative.Pointer")
	return fmt.Sprintf(result, key, value)
}

func lookupSignals(className string, file *ast.File) registrySignal {

	return registrySignal{}
}
