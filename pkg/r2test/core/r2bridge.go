package core

import (
	"fmt"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/arturoeanton/go-r2lang/pkg/r2libs"
	"github.com/arturoeanton/go-r2lang/pkg/r2test/assertions"
)

// r2BridgeContext accumulates the TestSuites built while evaluating an
// R2Lang test file's top-level describe()/it()/beforeEach()/... calls.
// currentSuite tracks which describe() block is currently executing (they
// don't nest in the examples this framework supports, but describe() still
// needs to restore the previous value on exit for robustness);
// currentTestName is set only while a test's own body (it()'s fn) is
// running, so assert.* panics can report which test they belong to.
type r2BridgeContext struct {
	suites          []*TestSuite
	currentSuite    *TestSuite
	currentTestName string
}

// newR2Environment builds an *r2core.Environment carrying the full standard
// library (the same modules pkg/r2lang.RunCode and the REPL register — test
// files should be able to use std/json/http/etc exactly like any other .r2
// script) plus the describe/it/beforeEach/afterEach/beforeAll/afterAll test
// DSL and the assert.* module, wired to build real TestSuite/TestCase
// values instead of the old regex-based placeholder. describe/it/... are
// registered as bare globals (like the "go"/"r2" concurrency keywords),
// not under a module namespace, since they're control-flow-shaped test DSL
// keywords rather than library data operations.
func newR2Environment(ctx *r2BridgeContext) *r2core.Environment {
	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)

	r2libs.RegisterLib(env)
	r2libs.RegisterStd(env)
	r2libs.RegisterIO(env)
	r2libs.RegisterHTTPClient(env)
	r2libs.RegisterRequests(env)
	r2libs.RegisterString(env)
	r2libs.RegisterRegex(env)
	r2libs.RegisterMath(env)
	r2libs.RegisterRand(env)
	r2libs.RegisterTest(env)
	r2libs.RegisterHTTP(env)
	r2libs.RegisterPrint(env)
	r2libs.RegisterOS(env)
	r2libs.RegisterHack(env)
	r2libs.RegisterEncoding(env)
	r2libs.RegisterConcurrency(env)
	r2libs.RegisterSync(env)
	r2libs.RegisterCollections(env)
	r2libs.RegisterValidate(env)
	r2libs.RegisterUnicode(env)
	r2libs.RegisterDate(env)
	r2libs.RegisterDB(env)
	r2libs.RegisterSOAP(env)
	r2libs.RegisterGRPC(env)
	r2libs.RegisterJSON(env)
	r2libs.RegisterXML(env)
	r2libs.RegisterCSV(env)
	r2libs.RegisterJWT(env)
	r2libs.RegisterConsole(env)
	r2libs.RegisterWeb(env)
	r2libs.RegisterGoInterOp(env)
	r2libs.RegisterGraph(env)

	registerTestDSL(env, ctx)

	return env
}

// asUserFunction validates that arg is an R2Lang function value, panicking
// with a clear R2Lang-level message (not a Go type-assertion panic) if not
// — matching the rest of this codebase's builtin-argument-validation
// convention.
func asUserFunction(builtinName, argName string, arg interface{}) *r2core.UserFunction {
	fn, ok := arg.(*r2core.UserFunction)
	if !ok {
		panic(fmt.Sprintf("%s: %s must be a function", builtinName, argName))
	}
	return fn
}

// registerTestDSL installs describe/it/beforeEach/afterEach/beforeAll/
// afterAll and the assert module into env, all populating ctx as they run.
func registerTestDSL(env *r2core.Environment, ctx *r2BridgeContext) {
	env.Set("describe", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("describe(name, fn) needs 2 arguments")
		}
		name, ok := args[0].(string)
		if !ok {
			panic("describe: first argument must be a string")
		}
		fn := asUserFunction("describe", "second argument", args[1])

		suite := &TestSuite{
			Name:        name,
			Description: name,
			Tests:       make([]*TestCase, 0),
			Tags:        make([]string, 0),
			Environment: env,
		}

		previous := ctx.currentSuite
		ctx.currentSuite = suite
		fn.Call()
		ctx.currentSuite = previous

		ctx.suites = append(ctx.suites, suite)
		return nil
	}))

	env.Set("it", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("it(name, fn) needs 2 arguments")
		}
		if ctx.currentSuite == nil {
			panic("it() must be called within a describe() block")
		}
		name, ok := args[0].(string)
		if !ok {
			panic("it: first argument must be a string")
		}
		fn := asUserFunction("it", "second argument", args[1])

		suite := ctx.currentSuite
		test := &TestCase{
			Name:        name,
			Description: name,
			Tags:        make([]string, 0),
			Suite:       suite,
			Func: func() {
				previousName := ctx.currentTestName
				ctx.currentTestName = name
				defer func() { ctx.currentTestName = previousName }()
				fn.Call()
			},
		}
		suite.Tests = append(suite.Tests, test)
		return nil
	}))

	registerHook := func(hookName string, assign func(*TestSuite, func())) {
		env.Set(hookName, r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic(hookName + "(fn) needs 1 argument")
			}
			if ctx.currentSuite == nil {
				panic(hookName + "() must be called within a describe() block")
			}
			fn := asUserFunction(hookName, "argument", args[0])
			assign(ctx.currentSuite, func() { fn.Call() })
			return nil
		}))
	}
	registerHook("beforeEach", func(s *TestSuite, f func()) { s.BeforeEach = f })
	registerHook("afterEach", func(s *TestSuite, f func()) { s.AfterEach = f })
	registerHook("beforeAll", func(s *TestSuite, f func()) { s.BeforeAll = f })
	registerHook("afterAll", func(s *TestSuite, f func()) { s.AfterAll = f })

	registerAssertModule(env, ctx)
}

// registerAssertModule bridges pkg/r2test/assertions' Assert type (a Go API
// that panics an *assertions.AssertionError on failure) to an "assert"
// module callable from R2Lang test scripts, matching the assert.equals(...)
// / assert.true(...) / etc. surface already used throughout
// examples/testing/*.r2. Every call constructs a fresh Assert using
// ctx.currentTestName so failure messages identify the failing test.
func registerAssertModule(env *r2core.Environment, ctx *r2BridgeContext) {
	currentAssert := func() *assertions.Assert {
		return assertions.NewAssert(ctx.currentTestName)
	}

	fn1 := func(name string, call func(a *assertions.Assert, arg interface{})) (string, r2core.BuiltinFunction) {
		return name, r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic(fmt.Sprintf("assert.%s needs 1 argument", name))
			}
			call(currentAssert(), args[0])
			return nil
		})
	}
	fn2 := func(name string, call func(a *assertions.Assert, x, y interface{})) (string, r2core.BuiltinFunction) {
		return name, r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic(fmt.Sprintf("assert.%s needs 2 arguments", name))
			}
			call(currentAssert(), args[0], args[1])
			return nil
		})
	}

	functions := map[string]r2core.BuiltinFunction{}
	set := func(name string, fn r2core.BuiltinFunction) { functions[name] = fn }

	set(fn2("equals", func(a *assertions.Assert, x, y interface{}) { a.Equals(x, y) }))
	set(fn2("notEquals", func(a *assertions.Assert, x, y interface{}) { a.NotEquals(x, y) }))
	set(fn1("true", func(a *assertions.Assert, x interface{}) { a.True(x) }))
	set(fn1("false", func(a *assertions.Assert, x interface{}) { a.False(x) }))
	set(fn1("nil", func(a *assertions.Assert, x interface{}) { a.Nil(x) }))
	set(fn1("notNil", func(a *assertions.Assert, x interface{}) { a.NotNil(x) }))
	set(fn2("greater", func(a *assertions.Assert, x, y interface{}) { a.Greater(x, y) }))
	set(fn2("greaterOrEqual", func(a *assertions.Assert, x, y interface{}) { a.GreaterOrEqual(x, y) }))
	set(fn2("less", func(a *assertions.Assert, x, y interface{}) { a.Less(x, y) }))
	set(fn2("lessOrEqual", func(a *assertions.Assert, x, y interface{}) { a.LessOrEqual(x, y) }))
	set(fn1("empty", func(a *assertions.Assert, x interface{}) { a.Empty(x) }))
	set(fn1("notEmpty", func(a *assertions.Assert, x interface{}) { a.NotEmpty(x) }))

	set("contains", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("assert.contains needs 2 arguments (haystack, needle)")
		}
		haystack, ok1 := args[0].(string)
		needle, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("assert.contains: both arguments must be strings")
		}
		currentAssert().Contains(haystack, needle)
		return nil
	}))
	set("notContains", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("assert.notContains needs 2 arguments (haystack, needle)")
		}
		haystack, ok1 := args[0].(string)
		needle, ok2 := args[1].(string)
		if !ok1 || !ok2 {
			panic("assert.notContains: both arguments must be strings")
		}
		currentAssert().NotContains(haystack, needle)
		return nil
	}))
	set("hasLength", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("assert.hasLength needs 2 arguments (collection, length)")
		}
		length, ok := args[1].(float64)
		if !ok {
			panic("assert.hasLength: second argument must be a number")
		}
		currentAssert().HasLength(args[0], int(length))
		return nil
	}))
	set("panics", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("assert.panics needs 1 argument (function)")
		}
		fn := asUserFunction("assert.panics", "argument", args[0])
		currentAssert().Panics(func() { fn.Call() })
		return nil
	}))
	set("notPanics", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("assert.notPanics needs 1 argument (function)")
		}
		fn := asUserFunction("assert.notPanics", "argument", args[0])
		currentAssert().NotPanics(func() { fn.Call() })
		return nil
	}))

	r2libs.RegisterModule(env, "assert", functions)
}
