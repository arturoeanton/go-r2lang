package r2libs

import (
	"fmt"
	"reflect"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// isNumericKind reports whether k is one of Go's numeric reflect.Kinds
// (integer or float families). Used to scope automatic type conversion to
// safe numeric widening/narrowing only — Go's general ConvertibleTo/Convert
// also permits numeric<->string conversions (treating an int as a Unicode
// code point), which would silently do the wrong thing for values crossing
// the R2Lang/Go boundary.
func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// Estructura que podemos usar para simular objetos Go
// Se crea con "goNew('MyStruct')" y se guarda en un *GoObject
type GoObject struct {
	value reflect.Value
}

// Eval => no hace nada especial
func (g *GoObject) Eval(env *r2core.Environment) interface{} {
	return g
}

// r2go.go: Interoperabilidad con Go
//
// Este módulo existe para el caso de uso normal de R2Lang: un programa Go
// que embebe el intérprete y quiere exponerle sus propias funciones/structs
// a los scripts .r2 (RegisterNativeFunc/RegisterNativeStruct, llamadas desde
// Go), quedando disponibles del lado del script bajo el namespace "native"
// (native.callFunc/new/setField/getField/callMethod).

// Un pequeño registro de funciones Go que deseamos exponer a R2
// map[funcName -> reflect.Value]
var goFuncRegistry = make(map[string]reflect.Value)

// Un pequeño registro de “constructores” (structName -> función constructor)
var goStructRegistry = make(map[string]func() interface{})

// RegisterNativeFunc expone una función Go bajo funcName para que los
// scripts R2Lang puedan invocarla via native.callFunc(funcName, ...args).
// Debe llamarse desde el programa Go anfitrión (antes o después de
// RegisterGoInterOp, el registro es un map global del paquete), nunca desde
// un script .r2 — reflect.Value no es un tipo representable en R2Lang.
func RegisterNativeFunc(funcName string, fn interface{}) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("RegisterNativeFunc('%s'): fn debe ser una función Go, se recibió %T", funcName, fn))
	}
	goFuncRegistry[funcName] = v
}

// RegisterNativeStruct expone un constructor Go bajo structName para que los
// scripts R2Lang puedan instanciarlo via native.new(structName). Debe
// llamarse desde el programa Go anfitrión, no desde un script .r2.
func RegisterNativeStruct(structName string, constructor func() interface{}) {
	goStructRegistry[structName] = constructor
}

// buildCallArgs convierte argumentos provenientes de R2Lang en []reflect.Value
// aptos para invocar una función/método Go de tipo fnType via reflect.Value.Call.
// reflect.ValueOf(nil) produce el "zero Value" inválido; pasárselo directamente
// a Call() panics con un mensaje críptico de bajo nivel
// ("reflect: Call using zero Value argument") en lugar de un error claro y
// consistente con el resto de las funciones de este módulo. Cuando conocemos
// el tipo del parámetro esperado (incluyendo funciones variádicas) sustituimos
// nil por el valor cero de ese tipo, lo cual permite pasar nil para
// punteros/interfaces/slices/maps como es habitual en Go.
func buildCallArgs(fnType reflect.Type, args []interface{}, label string) []reflect.Value {
	numIn := fnType.NumIn()
	callVals := make([]reflect.Value, len(args))
	for i, a := range args {
		var paramType reflect.Type
		switch {
		case fnType.IsVariadic() && i >= numIn-1:
			paramType = fnType.In(numIn - 1).Elem()
		case i < numIn:
			paramType = fnType.In(i)
		}
		if a == nil {
			if paramType == nil {
				panic(fmt.Sprintf("%s: no se puede pasar nil como argumento %d (no se pudo determinar el tipo esperado)", label, i))
			}
			callVals[i] = reflect.Zero(paramType)
			continue
		}
		argVal := reflect.ValueOf(a)
		// R2Lang numbers are always float64, so calling any registered Go
		// function/method whose parameter is a different numeric type (int,
		// int64, float32, ...) would previously reach reflect.Value.Call
		// directly and panic with the opaque, uncaught-by-design message
		// "reflect: Call using float64 as type int" instead of performing
		// the natural numeric conversion — confirmed via a real embedding
		// repro (native.callFunc("add", 3, 4) against a Go `func(int,int)
		// int`). Convert only for numeric-to-numeric mismatches, exactly
		// like native.setField already does for struct fields. Deliberately
		// NOT using the broader reflect.Type.ConvertibleTo/Convert for every
		// kind: Go's numeric<->string conversion rules treat an int as a
		// Unicode code point (e.g. float64(65) would silently become the
		// string "A" instead of erroring), which would be a surprising,
		// silent behavior change rather than a natural widening/narrowing.
		if paramType != nil && argVal.Type() != paramType && isNumericKind(argVal.Kind()) && isNumericKind(paramType.Kind()) {
			argVal = argVal.Convert(paramType)
		}
		callVals[i] = argVal
	}
	return callVals
}

// RegisterGoInterOp expone bajo el namespace "native" las funciones que
// permiten a un script R2Lang invocar funciones y structs Go previamente
// registrados desde el programa anfitrión via RegisterNativeFunc /
// RegisterNativeStruct. No hay forma de registrar una función/struct Go
// desde el propio script .r2 (reflect.Value no es representable en
// R2Lang) — ese registro siempre ocurre del lado Go, antes de ejecutar el
// script.
func RegisterGoInterOp(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		// native.callFunc("nombre", arg1, arg2, ...)
		// => Llama a la función Go registrada via RegisterNativeFunc.
		"callFunc": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("native.callFunc necesita (funcName, ...)")
			}
			funcName, ok := args[0].(string)
			if !ok {
				panic("native.callFunc: primer arg debe ser string (funcName)")
			}
			fnVal, found := goFuncRegistry[funcName]
			if !found {
				panic(fmt.Sprintf("native.callFunc: no se encontró la función '%s' (¿se llamó RegisterNativeFunc desde Go?)", funcName))
			}
			if fnVal.Kind() != reflect.Func {
				panic(fmt.Sprintf("native.callFunc: '%s' no es una función", funcName))
			}
			callArgs := buildCallArgs(fnVal.Type(), args[1:], fmt.Sprintf("native.callFunc('%s')", funcName))
			results := fnVal.Call(callArgs)
			if len(results) == 0 {
				return nil
			}
			if len(results) == 1 {
				return results[0].Interface()
			}
			arr := make([]interface{}, len(results))
			for i, r := range results {
				arr[i] = r.Interface()
			}
			return arr
		}),

		// native.new(structName) => crea un objeto Go registrado via
		// RegisterNativeStruct y retorna un *GoObject con acceso a sus
		// campos/métodos exportados.
		"new": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("native.new necesita (structName)")
			}
			sName, ok := args[0].(string)
			if !ok {
				panic("native.new: arg debe ser string (structName)")
			}
			constructor, found := goStructRegistry[sName]
			if !found {
				panic(fmt.Sprintf("native.new: no existe struct '%s' (¿se llamó RegisterNativeStruct desde Go?)", sName))
			}
			inst := constructor()
			return &GoObject{value: reflect.ValueOf(inst)}
		}),

		// native.setField(goObj, "FieldName", value) => setea un campo
		// exportado en la struct.
		"setField": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("native.setField(goObj, fieldName, value)")
			}
			obj, ok1 := args[0].(*GoObject)
			fieldName, ok2 := args[1].(string)
			if !(ok1 && ok2) {
				panic("native.setField: (GoObject, string, value)")
			}
			val := args[2]

			v := obj.value
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			fieldVal := v.FieldByName(fieldName)
			if !fieldVal.IsValid() {
				panic(fmt.Sprintf("native.setField: no existe el field '%s'", fieldName))
			}
			if !fieldVal.CanSet() {
				panic(fmt.Sprintf("native.setField: field '%s' no se puede setear (exportado?)", fieldName))
			}
			if val == nil {
				panic(fmt.Sprintf("native.setField: no se puede asignar nil al field '%s' de tipo %s", fieldName, fieldVal.Type()))
			}
			valReflect := reflect.ValueOf(val)
			// R2Lang numbers are always float64 (there is no int/float distinction
			// at the language level), so any Go struct field of another numeric
			// type (int, int64, float32, ...) would previously reach fieldVal.Set
			// directly and panic with an opaque low-level reflect message such as
			// "reflect.Set: value of type float64 is not assignable to type int",
			// instead of a clear, consistent error. Convert only for
			// numeric-to-numeric mismatches (see isNumericKind) and otherwise fail
			// with a descriptive message. Deliberately not using the broader
			// reflect.Type.ConvertibleTo for every kind: Go's numeric<->string
			// conversion rules treat an int as a Unicode code point (e.g.
			// float64(65) would silently become the string "A" instead of
			// erroring), which would be a surprising, silent behavior change
			// rather than a natural widening/narrowing.
			if !valReflect.Type().AssignableTo(fieldVal.Type()) {
				if isNumericKind(valReflect.Kind()) && isNumericKind(fieldVal.Kind()) {
					valReflect = valReflect.Convert(fieldVal.Type())
				} else {
					panic(fmt.Sprintf("native.setField: no se puede asignar un valor de tipo %s al field '%s' de tipo %s", valReflect.Type(), fieldName, fieldVal.Type()))
				}
			}
			fieldVal.Set(valReflect)
			return nil
		}),

		// native.getField(goObj, "FieldName") => obtiene un campo.
		"getField": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("native.getField(goObj, fieldName)")
			}
			obj, ok1 := args[0].(*GoObject)
			fieldName, ok2 := args[1].(string)
			if !(ok1 && ok2) {
				panic("native.getField: (GoObject, string)")
			}
			v := obj.value
			if v.Kind() == reflect.Ptr {
				v = v.Elem()
			}
			fieldVal := v.FieldByName(fieldName)
			if !fieldVal.IsValid() {
				panic(fmt.Sprintf("native.getField: no existe field '%s'", fieldName))
			}
			return fieldVal.Interface()
		}),

		// native.callMethod(goObj, "MethodName", ...args) => llama un
		// método exportado en la struct.
		"callMethod": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("native.callMethod(goObj, methodName, ...)")
			}
			obj, ok1 := args[0].(*GoObject)
			methodName, ok2 := args[1].(string)
			if !(ok1 && ok2) {
				panic("native.callMethod: (GoObject, string, ...)")
			}
			callArgs := args[2:]
			v := obj.value
			m := v.MethodByName(methodName)
			if !m.IsValid() {
				panic(fmt.Sprintf("native.callMethod: no existe método '%s'", methodName))
			}
			inVals := buildCallArgs(m.Type(), callArgs, fmt.Sprintf("native.callMethod('%s')", methodName))
			results := m.Call(inVals)
			if len(results) == 0 {
				return nil
			}
			if len(results) == 1 {
				return results[0].Interface()
			}
			arr := make([]interface{}, len(results))
			for i, r := range results {
				arr[i] = r.Interface()
			}
			return arr
		}),
	}

	RegisterModule(env, "native", functions)
}
