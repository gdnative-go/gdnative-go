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
// limitations under the License

package gdnative

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"reflect"
	"strings"
)

const (
	godotRegister    string = "godot::register"
	godotConstructor string = "godot::constructor"
	godotDestructor  string = "godot::destructor"
	godotExport      string = "godot::export"
)

// LookupRegistrableTypeDeclarations parses the given package AST and adds any
// relevant data to the registry so we can generate boilerplate registration code
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
						original := strings.TrimSpace(strings.ReplaceAll(line.Text, "/", ""))
						docstring := strings.ToLower(original)
						if strings.HasPrefix(docstring, godotRegister) {

							className := getClassName(tp)
							class := registryClass{
								base: getBaseClassName(sp),
							}
							classes[className] = &class

							// set alias if defined
							if strings.Contains(docstring, " as ") {
								class.alias = strings.TrimSpace(strings.Split(strings.Split(original, " as ")[1], " ")[0])
							}
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
			class.AddSignals(lookupSignals(className, file))
			class.AddProperties(lookupProperties(className, file))
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
	for i := 0; i < sp.Fields.NumFields()-1; i++ {
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
						fmt.Printf("could not parse constructor comment %s, many parenthesis", docstring)
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
						fmt.Printf("could not parse destructor comment %s, many parenthesis", docstring)
						os.Exit(1)
					}

					structName = structName[1 : len(structName)-1]
					if structName != className {
						// this destructor doesn't match with our class, skip it
						continue
					}

					destructor, err := validateDestructor(structName, fd)
					if err != nil {
						fmt.Printf("could not validate destructor: %s", err)
						os.Exit(1)
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
			fmt.Printf(
				"init function present on files, you can not provide your own init function while autoregistering classes",
			)
			os.Exit(1)
		}

		// ignore non methods
		if fd.Recv == nil {
			continue
		}

		alias, exported := extractExportedAndAliasFromDoc(fd.Doc)

		// ignore non exported methods
		if !fd.Name.IsExported() && !exported {
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
			alias:        alias,
			class:        className,
			params:       lookupParams(fd.Type.Params),
			returnValues: lookupReturnValues(fd),
		}
		methods = append(methods, &method)
	}

	return methods
}

// lookupSignals look up for every signal that is owned by the type and fill
// a registration data structure with it
func lookupSignals(className string, file ast.Node) []*registrySignal {

	signals := []*registrySignal{}
	ast.Inspect(file, func(node ast.Node) (cont bool) {

		cont = true
		_, ok := node.(*ast.StructType)
		if ok {
			return
		}

		cl, ok := node.(*ast.CompositeLit)
		if !ok {
			return
		}

		st, ok := cl.Type.(*ast.SelectorExpr)
		if !ok {
			return
		}

		if st.X.(*ast.Ident).Name == "gdnative" && st.Sel.Name == "Signal" {
			signal := registrySignal{}
			for i := range cl.Elts {
				kv, ok := cl.Elts[i].(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				key := kv.Key.(*ast.Ident).Name
				switch key {
				case "Name":
					signal.name = kv.Value.(*ast.BasicLit).Value
				case "Args":
					arguments, ok := kv.Value.(*ast.CompositeLit)
					if !ok {
						fmt.Printf(
							"WARNING: arguments on signal %s has the wrong type, it should be *ast.CompositeLit, arguments will be ignored\n",
							signal.name,
						)
						signal.args = "[]*gdnative.SignalArgument{}"
						continue
					}
					signal.args = parseSignalArgs(arguments)
				case "DefaultArgs":
					defaultArguments, ok := kv.Value.(*ast.CompositeLit)
					if !ok {
						fmt.Printf(
							"WARNING: default arguments on signal %s has the wrong type, it should be *ast.CompositeLit, arguments will be ignored\n",
							signal.name,
						)
						signal.defaults = "[]*gdnative.Variant{}"
						continue
					}
					signal.defaults = parseSignalArgs(defaultArguments)
				}
			}

			signals = append(signals, &signal)
		}

		return
	})

	return signals
}

func lookupParams(fields *ast.FieldList) []*registryMethodParam {

	params := []*registryMethodParam{}
	if fields.NumFields() > 0 {
		for _, field := range fields.List {

			kind := parseDefault(field.Type, "")
			switch kind {
			case "":
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

func lookupProperties(className string, file *ast.File) []*registryProperty {

	properties := []*registryProperty{}
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

			if getClassName(tp) != className {
				continue
			}

			if sp.Fields.NumFields() > 0 {
				for _, field := range sp.Fields.List {
					if field.Names == nil {
						// this field is a selector expression, skip it
						continue
					}

					gdnativeKind := ""
					kind := parseDefault(field.Type, "")
					switch kind {
					case "":
						// we don't have a type skip iteration
						continue
					case "gdnative.Signal":
						// this is a signal skip it
						continue
					default:
						gdnativeKind = kind
						typeFormat := fmt.Sprintf("VariantType%s", strings.ReplaceAll(kind, "gdnative.", ""))
						_, ok := VariantTypeLookupMap[typeFormat]
						if ok {
							kind = fmt.Sprintf("gdnative.%s", typeFormat)
							break
						}

						kind = "gdnative.VariantTypeObject"
					}

					alias, exported := extractExportedAndAliasFromDoc(field.Doc)

					tmpProperties := []*registryProperty{}
					for _, name := range field.Names {
						if !name.IsExported() && !exported {
							continue
						}

						tmpProperties = append(tmpProperties, &registryProperty{
							name:         name.String(),
							alias:        alias,
							kind:         kind,
							gdnativeKind: gdnativeKind,
						})
					}

					if field.Tag != nil {
						// check if this tag is the ignore tag
						if field.Tag.Value == "`-`" || field.Tag.Value == "`_`" || field.Tag.Value == "`omit`" {
							continue
						}

						// create a fake reflect.StructTag to lookup our keys
						fakeTag := reflect.StructTag(strings.ReplaceAll(field.Tag.Value, "`", ""))
						rset, rsetOk := fakeTag.Lookup("rset_type")
						usage, usageOk := fakeTag.Lookup("usage")
						hint, hintOk := fakeTag.Lookup("hint")
						hintString, hintStringOk := fakeTag.Lookup("hint_string")
						if !hintStringOk {
							hintString = ""
						}
						if !hintOk {
							hint = "None"
						}
						if !rsetOk {
							rset = "Disabled"
						}
						if !usageOk {
							usage = "Default"
						}

						for i := range tmpProperties {
							mustSetPropertyTagHint(tmpProperties[i], hint)
							mustSetPropertyTagRset(tmpProperties[i], rset)
							mustSetPropertyTagUsage(tmpProperties[i], usage)
							tmpProperties[i].hintString = hintString
						}
						properties = append(properties, tmpProperties...)
					}
				}
			}
		}
	}

	return properties
}

func mustSetPropertyTagRset(property *registryProperty, value string) {

	rpcMode := fmt.Sprintf("MethodRpcMode%s", strings.Title(value))
	_, ok := MethodRpcModeLookupMap[rpcMode]
	if !ok {
		valid := []string{}
		for key := range MethodRpcModeLookupMap {
			valid = append(valid, strings.ToLower(key[13:]))
		}
		fmt.Printf(
			"on property %s: unknown rset_type %s, it must be one of:\n\t%s\n",
			property.name, value, strings.Join(valid, "\n\t"),
		)
		os.Exit(1)
	}
	property.rset = strings.Title(value)
}

func mustSetPropertyTagHint(property *registryProperty, value string) {

	hint := fmt.Sprintf("PropertyHint%s", strings.Title(value))
	_, ok := PropertyHintLookupMap[hint]
	if !ok {
		valid := []string{}
		for key := range PropertyHintLookupMap {
			valid = append(valid, strings.ToLower(key[12:]))
		}
		fmt.Printf("on property %s: unknown hint %s, it must be one of %s", property.name, value, strings.Join(valid, ", "))
		os.Exit(1)
	}
	property.hint = strings.Title(value)
}

func mustSetPropertyTagUsage(property *registryProperty, value string) {

	usage := fmt.Sprintf("PropertyUsage%s", strings.Title(value))
	_, ok := PropertyUsageFlagsLookupMap[usage]
	if !ok {
		valid := []string{}
		for key := range PropertyUsageFlagsLookupMap {
			valid = append(valid, strings.ToLower(key[13:]))
		}
		fmt.Printf("on property %s: unknown usage %s, it must be one of %s", property.name, value, strings.Join(valid, ", "))
		os.Exit(1)
	}
	property.usage = strings.Title(value)
}

func extractExportedAndAliasFromDoc(doc *ast.CommentGroup) (string, bool) {

	var alias string
	var exported bool

	if doc != nil {
		for _, line := range doc.List {
			docstring := strings.TrimSpace(strings.ReplaceAll(line.Text, "/", ""))
			if strings.HasPrefix(docstring, godotExport) {
				exported = true
				if strings.Contains(docstring, " as ") {
					alias = strings.TrimSpace(strings.Split(strings.Split(docstring, " as ")[1], " ")[0])
				}
				break
			}
		}
	}

	return alias, exported
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

func parseKeyValueExpr(expr *ast.KeyValueExpr) (string, string) { //nolint:unused

	var value string
	key := expr.Key.(*ast.Ident).Name
	switch t := expr.Value.(type) {
	case *ast.BasicLit:
		value = t.Value
	case *ast.CompositeLit:
		switch t2 := t.Type.(type) {
		case *ast.KeyValueExpr:
			_, value = parseKeyValueExpr(t2)
		default:
			value = parseDefault(t2, "gdnative.Pointer")
		}
	}

	return key, value
}

func parseSignalArgs(composite *ast.CompositeLit) string {

	buffer := []byte{}
	buf := bytes.NewBuffer(buffer)
	fileSet := token.NewFileSet()

	printer.Fprint(buf, fileSet, composite)
	return buf.String()
}
