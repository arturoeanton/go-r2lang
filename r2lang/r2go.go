package r2lang

import (
	"fmt"
	"reflect"
)

// Estructura que podemos usar para simular objetos Go
// Se crea con "goNew('MyStruct')" y se guarda en un *GoObject
type GoObject struct {
	value reflect.Value
}

// Eval => no hace nada especial
func (g *GoObject) Eval(env *Environment) interface{} {
	return g
}

// r2go.go: Interoperabilidad con Go

// Un pequeño registro de funciones Go que deseamos exponer a R2
// map[funcName -> reflect.Value]
var goFuncRegistry = make(map[string]reflect.Value)

// Un pequeño registro de “constructores” (structName -> función constructor)
var goStructRegistry = make(map[string]func() interface{})

// RegisterGoInterOp: expone funciones que permiten a R2 usar el registro
func RegisterGoInterOp(env *Environment) {

	// 1) goRegisterFunc("nombre", GoFuncion)
	//    => En Go:   goRegisterFunc("miSum", reflect.ValueOf(MiSum))
	//    => En R2:   callGoFunc("miSum", 10, 20)
	env.Set("goRegisterFunc", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("goRegisterFunc necesita (name, goFuncValueReflect)")
		}
		_, ok1 := args[0].(string)
		if !ok1 {
			panic("goRegisterFunc: primer arg debe ser string")
		}
		// El segundo arg DEBE ser un reflect.Value dentro de Go,
		// pero como no se puede pasar reflect.Value desde R2, simulemos
		// que lo guardamos en 'goFuncRegistry' manualmente en Go.
		// => “Truco”: generamos un panic pidiendo que se registre desde Go.
		// Este ejemplo asume que la parte en Go hará:
		//    goFuncRegistry["miSum"] = reflect.ValueOf(MiSum)
		// y en R2 se llama "callGoFunc".
		panic("goRegisterFunc: se debe llamar desde Go, no desde R2, para inyectar la reflect.Value. (Truco de ejemplo)")
	}))

	// 2) callGoFunc("nombre", arg1, arg2, ...)
	//    => Llama a la función de Go que se registró en goFuncRegistry
	env.Set("callGoFunc", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("callGoFunc necesita (funcName, ...)")
		}
		funcName, ok := args[0].(string)
		if !ok {
			panic("callGoFunc: primer arg debe ser string (funcName)")
		}
		fnVal, found := goFuncRegistry[funcName]
		if !found {
			panic(fmt.Sprintf("callGoFunc: no se encontró la función '%s' en goFuncRegistry", funcName))
		}
		if fnVal.Kind() != reflect.Func {
			panic(fmt.Sprintf("callGoFunc: '%s' no es una función", funcName))
		}
		// Convertimos los args[1..] a reflect.Value
		callArgs := make([]reflect.Value, len(args)-1)
		for i := 1; i < len(args); i++ {
			callArgs[i-1] = reflect.ValueOf(args[i])
		}
		// Llamamos la función
		results := fnVal.Call(callArgs)
		// Si hay 0 resultados => nil
		if len(results) == 0 {
			return nil
		}
		// si 1 => lo retornamos
		if len(results) == 1 {
			return results[0].Interface()
		}
		// si mas => array
		arr := make([]interface{}, len(results))
		for i, r := range results {
			arr[i] = r.Interface()
		}
		return arr
	}))

	// 3) goRegisterStruct("Nombre", constructor)
	// => para permitir goNew("Nombre") en R2
	env.Set("goRegisterStruct", BuiltinFunction(func(args ...interface{}) interface{} {
		panic("goRegisterStruct: se debe llamar desde Go con la map, no desde R2 (Truco).")
	}))

	// 4) goNew(structName) => crea un objeto Go y retorna un *GoObject
	env.Set("goNew", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("goNew necesita (structName)")
		}
		sName, ok := args[0].(string)
		if !ok {
			panic("goNew: arg debe ser string (structName)")
		}
		constructor, found := goStructRegistry[sName]
		if !found {
			panic(fmt.Sprintf("goNew: no existe struct '%s' en goStructRegistry", sName))
		}
		inst := constructor() // crea una instancia
		return &GoObject{value: reflect.ValueOf(inst)}
	}))

	// 5) goSetField(goObj, "FieldName", value)
	// => setea un campo exportado en la struct
	env.Set("goSetField", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("goSetField(goObj, fieldName, value)")
		}
		obj, ok1 := args[0].(*GoObject)
		fieldName, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("goSetField: (GoObject, string, value)")
		}
		val := args[2]

		// reflexion para setear
		v := obj.value
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		fieldVal := v.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			panic(fmt.Sprintf("goSetField: no existe el field '%s'", fieldName))
		}
		if !fieldVal.CanSet() {
			panic(fmt.Sprintf("goSetField: field '%s' no se puede setear (exportado?)", fieldName))
		}
		fieldVal.Set(reflect.ValueOf(val))
		return nil
	}))

	// 6) goGetField(goObj, "FieldName") => obtiene un campo
	env.Set("goGetField", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("goGetField(goObj, fieldName)")
		}
		obj, ok1 := args[0].(*GoObject)
		fieldName, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("goGetField: (GoObject, string)")
		}
		v := obj.value
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		fieldVal := v.FieldByName(fieldName)
		if !fieldVal.IsValid() {
			panic(fmt.Sprintf("goGetField: no existe field '%s'", fieldName))
		}
		return fieldVal.Interface()
	}))

	// 7) goCallMethod(goObj, "MethodName", ...args)
	// => llama un método exportado en la struct
	env.Set("goCallMethod", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("goCallMethod(goObj, methodName, ...)")
		}
		obj, ok1 := args[0].(*GoObject)
		methodName, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("goCallMethod: (GoObject, string, ...)")
		}
		callArgs := args[2:] // lo que sigue son parámetros
		// reflexion
		v := obj.value
		m := v.MethodByName(methodName)
		if !m.IsValid() {
			panic(fmt.Sprintf("goCallMethod: no existe método '%s'", methodName))
		}
		// convert callArgs => reflect.Value
		inVals := make([]reflect.Value, len(callArgs))
		for i := 0; i < len(callArgs); i++ {
			inVals[i] = reflect.ValueOf(callArgs[i])
		}
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
	}))
}

// En Go, para registrar tus funciones y structs:

// goFuncRegistry["miSuma"] = reflect.ValueOf(func(a,b int) int {return a+b})
// goStructRegistry["Persona"] = func() interface{} { return &Persona{} }
