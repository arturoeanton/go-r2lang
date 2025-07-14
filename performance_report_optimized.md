=== REPORTE DE RENDIMIENTO R2LANG ===
Ejecutando benchmarks...
Sistema: darwin
Arquitectura: arm64
CPUs: 14
Versión Go: go1.24.4
Reporte generado: performance_report.md
Ejecuta: go test -bench=. -benchmem performance_test.go
goos: darwin
goarch: arm64
cpu: Apple M4 Max
BenchmarkBasicArithmetic-14      	    6733	    175645 ns/op	   86500 B/op	    8071 allocs/op
BenchmarkStringOperations-14     	   22992	     54982 ns/op	  120147 B/op	    1073 allocs/op
panic: access to property in unsupported type: []interface {}

goroutine 7 [running]:
github.com/arturoeanton/go-r2lang/pkg/r2core.(*AccessExpression).Eval(0x1400051a080, 0x14000200c60)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/access_expression.go:25 +0x198
github.com/arturoeanton/go-r2lang/pkg/r2core.(*CallExpression).Eval(0x140000b2a50, 0x14000200c60)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/call_expression.go:22 +0xe0
github.com/arturoeanton/go-r2lang/pkg/r2core.(*ExprStatement).Eval(0x101227e80?, 0x10125c1c0?)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/expr_statement.go:8 +0x28
github.com/arturoeanton/go-r2lang/pkg/r2core.(*BlockStatement).Eval(0x140000b2600?, 0x14000200c60)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/block_statement.go:10 +0x6c
github.com/arturoeanton/go-r2lang/pkg/r2core.(*ForStatement).executeStandardLoop(0x140000fe150, 0x14000200c60)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/for_statement.go:100 +0x198
github.com/arturoeanton/go-r2lang/pkg/r2core.(*ForStatement).evalStandardFor(0x140000fe150, 0x14000200c60)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/for_statement.go:80 +0xe0
github.com/arturoeanton/go-r2lang/pkg/r2core.(*ForStatement).Eval(0x140000fe150, 0x14000200c00)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/for_statement.go:26 +0xf4
github.com/arturoeanton/go-r2lang/pkg/r2core.(*BlockStatement).Eval(0x14000113ce8?, 0x14000200c00)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/block_statement.go:10 +0x6c
github.com/arturoeanton/go-r2lang/pkg/r2core.(*UserFunction).NativeCall(0x140000c0b80, 0x0?, {0x0, 0x0, 0x9?})
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/user_function.go:44 +0x2c4
github.com/arturoeanton/go-r2lang/pkg/r2core.(*UserFunction).Call(0x140000c0b80, {0x0?, 0x14000113de8?, 0x100fd5fa0?})
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/user_function.go:54 +0x74
github.com/arturoeanton/go-r2lang/pkg/r2core.(*CallExpression).Eval(0x140000b36e0, 0x14000200ba0)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/call_expression.go:35 +0x24c
github.com/arturoeanton/go-r2lang/pkg/r2core.(*ExprStatement).Eval(0x140000b36b0?, 0x14000200ba0?)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/expr_statement.go:8 +0x28
github.com/arturoeanton/go-r2lang/pkg/r2core.(*Program).Eval(0x14000200ba0?, 0x14000200ba0)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/pkg/r2core/program.go:11 +0x6c
command-line-arguments.BenchmarkArrayOperations(0x140000ea588)
	/Users/arturoeliasanton/github.com/arturoeanton/go-r2lang/performance_test.go:90 +0x4c
testing.(*B).runN(0x140000ea588, 0x1)
	/opt/homebrew/Cellar/go/1.24.4/libexec/src/testing/benchmark.go:219 +0x1a4
testing.(*B).run1.func1()
	/opt/homebrew/Cellar/go/1.24.4/libexec/src/testing/benchmark.go:245 +0x4c
created by testing.(*B).run1 in goroutine 1
	/opt/homebrew/Cellar/go/1.24.4/libexec/src/testing/benchmark.go:238 +0x90
exit status 2
FAIL	command-line-arguments	3.970s
FAIL
