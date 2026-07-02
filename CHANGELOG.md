# Changelog

All notable changes to R2Lang are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/) (`MAJOR.MINOR.PATCH`).
R2Lang is pre-1.0, so `MINOR` releases may still include breaking changes; `PATCH`
releases are bug fixes and hardening only.

**Tag naming**: tags up to and including `v0.1.25` used the format
`vMAJOR.MINOR.PATCH_short-description` (the description was informal shorthand
for a fast-moving hardening sprint). Starting from the next release, tags use
plain `vMAJOR.MINOR.PATCH` — the description now lives only in this file and in
the annotated tag message. Existing tags are kept as-is; they are historical
record and are not renamed.

## [Unreleased]

## [0.1.26] - CHANGELOG and release process
### Added
- `CHANGELOG.md` (this file) and `RELEASING.md`, documenting the release
  history and formalizing semver going forward.

## [0.1.25] - CI hardening, example sweep, r2test fixes
### Fixed
- `pkg/r2test` mocking: `Times(0)`/`AtMost(0)` expectations were a silent
  no-op (a mock declared to reject all calls accepted unlimited calls).
- `pkg/r2test` mocking: isolation-context IDs could collide and alias an
  existing context after cleanup; `IsolationContext` had no synchronization
  despite backing concurrent test execution (confirmed with `-race`).
- `pkg/r2test` mocking: `Spy.callOriginal` panicked on `nil` arguments or
  arity mismatches instead of returning a clean error.
- `pkg/r2test` coverage: `GetSortedFiles` produced `NaN`-corrupted sort
  order for files with zero trackable lines.
- `pkg/r2test` fixtures: `FixtureManager.Load`'s ambiguous-extension
  resolution was non-deterministic (iterated a Go map).
- `pkg/r2test` reporters: `GenerateJUnitReportWithProperties` silently
  discarded the properties it was documented to attach.
- `examples/example42-os-advanced.r2`: `std.sleep(100)` slept 100 seconds
  instead of the intended 100ms (`std.sleep` takes seconds, not ms).
### Changed
- CI's `integration` job now actually fails the build on a real regression
  instead of swallowing every failure with `|| echo`.
- Added `scripts/sweep_examples.sh` (+ `scripts/known_flaky_examples.txt`),
  a versioned, documented sweep across every example script, wired into CI.
- Added a minimum test-coverage floor (40%) to CI.

## [0.1.24] - Remove dead bytecode/JIT code
### Removed
- `pkg/r2core/bytecode.go` and `pkg/r2core/jit_loop.go`: confirmed dead
  code with zero references anywhere in the codebase.

## [0.1.23] - Go interop and graph modules, comprehension/template fixes
### Added
- `native.*` module: exposes Go struct/function interop
  (`native.callFunc/new/setField/getField/callMethod`) to R2Lang scripts,
  plus exported `RegisterNativeFunc`/`RegisterNativeStruct` for host Go
  programs embedding the interpreter.
- `graph.*` module: a directed graph with a fluent API (`addEdge`,
  `getAncestors`, `getDescendants`, `getShortestPath`, ...).
### Fixed
- Nested array/object comprehensions couldn't see earlier generators'
  bindings (`[y for row in matrix for y in row]` panicked).
- Comprehending over a map bound the loop variable to the value instead of
  the key, inconsistent with `for-in`.
- `formatPrintf` fed raw `float64` to `%d`/`%x`/`%s`-style verbs, producing
  garbage instead of converting first.
- Thousands-separator comma grouping corrupted non-numeric strings like
  `"+Inf"` and used scientific notation for large decimals.
- Negative/malformed precision in currency/float/percentage format specs
  produced garbled output.
- Object destructuring and match patterns didn't recognize the
  `map[string]*Variable` representation used by aliased module imports.
- `native.callFunc`/`native.setField` panicked with a raw reflect error
  when passing R2Lang's `float64` numbers to Go `int`-typed
  parameters/fields instead of converting.

## [0.1.22] - Content-type detection
### Fixed
- `DetectContentType` misclassified well-formed XML as `text/html` when it
  contained a tag with a generic name (`<a>`, `<p>`, `<div>`, ...) anywhere
  in the document.

## [0.1.21] - JWT/goroutine hardening
### Fixed
- `Monitor.Unlock()`/`Wait()` (goroutine module) could crash the entire
  interpreter process with an unrecoverable Go `fatal error` when misused
  from a script.
- `jwt.refresh` bypassed expiration/`nbf`/algorithm checks entirely,
  reviving expired tokens indefinitely.
- `httpclient.stringifyXML` had no cycle protection (unrecoverable stack
  overflow on self-referential values) and no response-size limit.
- `unicode.utitle` uppercased instead of title-casing;
  `unicode.ugetCategory` was nondeterministic.
- Hardened `native.*`'s (then still `r2go.go`) reflect-based Go interop
  against nil arguments and numeric type mismatches.

## [0.1.20] - deepCopy stack overflow
### Fixed
- `std.deepCopy` crashed on self-referential arrays/maps and corrupted
  pointer-to-struct interpreter values (functions, dates).
- `std.is`, `std.join`, `rand.randChoice/shuffle/sample` didn't recognize
  `r2core.InterfaceSlice` (only `[]interface{}`).
- `std.print` leaked raw Go pointer addresses when printing function
  values.

## [0.1.19] - JSON/XML parsing
### Fixed
- XML namespace declarations leaked in as regular content attributes.
- JSON→XML conversion silently dropped repeated-tag arrays and scalar leaf
  values.
- `JSON.stringify` leaked raw Go pointer addresses for unsupported types
  instead of erroring.

## [0.1.18] - CSV/console/print
### Fixed
- UTF-8 BOM not stripped from CSV input.
- Non-deterministic column order in `csv.stringify`/`writeFile`.
- Torn/interleaved `console.*` output under concurrent goroutines.
- `printHeader`'s separator width used byte length instead of rune count.

## [0.1.17] - REPL fixes
### Fixed
- The REPL was missing the entire standard library (only concurrency
  builtins were registered).
- Parser/lexer errors called `os.Exit(1)` directly, bypassing all
  `recover()`.
- EOF (Ctrl+D / piped stdin) caused an infinite busy loop.
- Naive incomplete-input detection replaced with real lexer-based bracket
  tracking.
- Command history was never actually persisted; `.exit` skipped cleanup.

## [0.1.16] - DB/SOAP core
### Fixed
- SQL `INTEGER`/`DATETIME` values silently produced wrong arithmetic
  (string concatenation) instead of numeric conversion.
- Postgres placeholder syntax (`$1,$2,...`) wasn't translated from `?`.
- SOAP envelope building didn't XML-escape interpolated parameter values.
- Non-200 SOAP responses that were still valid SOAP XML were treated as
  hard transport errors instead of being parsed for Faults.

## [0.1.15] - gRPC client streaming
### Fixed
- Client-streaming and bidi-streaming gRPC calls were fundamentally broken
  (`closeSend()` cancelled the stream instead of closing it properly).
- `repeated`/`map` proto fields weren't converted, causing panics.
### Added
- `date` module: `clone`, `startOfDay`/`endOfDay`, `startOfMonth`/`endOfMonth`,
  `isLeapYear`, `isWeekend`, `daysInMonth`.

## [0.1.14] - validate/deepEqual
### Added
- `validate` module: `isEmail`/`isURL`/`isIP`.
- `collections.deepEqual`/`deepClone` with cycle detection.

## [0.1.13] - import/finally fixes
### Fixed
- A data race in import bookkeeping (`Environment.imported`/`importStack`).
- `finally` blocks were skipped when the `catch` block itself threw.

## [0.1.12] - DSL recursion limits
### Fixed
- A reentrant deadlock in `dsl{}`'s `.use()`.
- Grammar errors from `token()`/`rule()` were silently discarded.
- DSL action calls had no recursion/timeout limit (unlike regular
  functions), so infinite recursion inside a DSL action hung forever.

## [0.1.11] - r2web routing
### Fixed
- Path parameters never worked; `ctx.params/query/headers` panicked on
  every access; registering the same path with different methods crashed
  `.listen()` at startup.

## [0.1.10] - flatten/NaN
### Fixed
- `collections.flatten` crashed the process (unrecoverable stack overflow)
  on self-referential arrays.
- `math.isEven`/`isOdd`/`isPrime`/`roundTo` had undefined behavior on
  `NaN`/`Inf` input.

## [0.1.9] - gRPC deadlock
### Fixed
- A reentrant deadlock when closing gRPC streams.

## [0.1.8] - AfterEach logging
### Fixed
- A panicking `AfterEach` hook was silently discarded instead of logged.

## [0.1.7] - Mutex crash
### Fixed
- `Mutex.unlock()` without a prior `lock()` crashed the process
  unrecoverably instead of raising a catchable panic.

## [0.1.6] - encoding stdlib
### Added
- `encoding` module: `base64`/`hex`/`url` encode/decode; `uuid.v4()`.
### Fixed
- `filter`/`find`/`sort` callbacks used a raw `.(bool)` type assertion
  instead of boolean coercion, panicking on non-boolean predicate results.

## [0.1.5] - super() fix
### Fixed
- `super.method()` calls corrupted environment state across sequential
  calls.
### Added
- Optional chaining with index access: `x?.[i]`.

## [0.1.4] - sync stdlib
### Added
- `sync` module: `Mutex`, `WaitGroup`, `Semaphore`, `Once`.
### Fixed
- Several broken examples (namespacing, non-mutating `.push()`, HTTP
  response handling).

## [0.1.3] - reserved members
### Fixed
- Reserved words couldn't be used as a member name after a dot (broke
  calls like `assert.true(...)`).
### Added
- `toUpperCase`/`toLowerCase` aliases for `toUpper`/`toLower`.

## [0.1.2] - critical fixes
### Fixed
- Fluent-API methods on `io.Path`/`io.FileStream`/`os.Command` were
  completely unreachable (a dispatch gap in `AccessExpression.Eval`).
- A coverage HTML report crash.
- Test hooks (`BeforeEach`/`AfterEach`/etc.) weren't panic-protected.
### Added
- `regex` module.

## [0.1.1] - DSL hardening
### Added
- `.tokens()/.ast()/.check()/.diagnostics()/.completions()` introspection
  on the `dsl{}` object; `keyword()`/`literal()` builtins.
### Fixed
- A `currentExecutionEnv` race in the DSL engine.

## [0.1.0] - go-dsl migration
### Changed
- Migrated the DSL engine from a hand-rolled backtracking parser to
  [`github.com/arturoeanton/go-dsl`](https://github.com/arturoeanton/go-dsl)
  v1.4.0.
### Fixed
- Multiline function calls/parameters and map-literal separators in the
  parser.

## [0.0.1] - Initial tag
- First tagged version of the project.

[Unreleased]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.26...HEAD
[0.1.26]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.25_ci-sweep-r2test-fixes...v0.1.26
[0.1.25]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.24_remove-dead-bytecode...v0.1.25_ci-sweep-r2test-fixes
[0.1.24]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.23_native-graph-comprehensions...v0.1.24_remove-dead-bytecode
[0.1.23]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.22_content-type-detection...v0.1.23_native-graph-comprehensions
[0.1.22]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.21_jwt-goroutine-hardening...v0.1.22_content-type-detection
[0.1.21]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.20_deepcopy-stackoverflow...v0.1.21_jwt-goroutine-hardening
[0.1.20]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.19_json-xml-parsing...v0.1.20_deepcopy-stackoverflow
[0.1.19]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.18_csv-console-print...v0.1.19_json-xml-parsing
[0.1.18]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.17_repl-fixes...v0.1.18_csv-console-print
[0.1.17]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.16_db-soap-core...v0.1.17_repl-fixes
[0.1.16]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.15_grpc-client-streaming...v0.1.16_db-soap-core
[0.1.15]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.14_validate-deepequal...v0.1.15_grpc-client-streaming
[0.1.14]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.13_import-finally...v0.1.14_validate-deepequal
[0.1.13]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.12_dsl-recursion...v0.1.13_import-finally
[0.1.12]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.11_web-routing...v0.1.12_dsl-recursion
[0.1.11]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.10_flatten-nan...v0.1.11_web-routing
[0.1.10]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.9_grpc-deadlock...v0.1.10_flatten-nan
[0.1.9]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.8_aftereach-logging...v0.1.9_grpc-deadlock
[0.1.8]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.7_mutex-crash...v0.1.8_aftereach-logging
[0.1.7]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.6_encoding-stdlib...v0.1.7_mutex-crash
[0.1.6]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.5_super-fix...v0.1.6_encoding-stdlib
[0.1.5]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.4_sync-stdlib...v0.1.5_super-fix
[0.1.4]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.3_reserved-members...v0.1.4_sync-stdlib
[0.1.3]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.2_critical-fixes...v0.1.3_reserved-members
[0.1.2]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.1_hardening-dsl...v0.1.2_critical-fixes
[0.1.1]: https://github.com/arturoeanton/go-r2lang/compare/v0.1.0...v0.1.1_hardening-dsl
[0.1.0]: https://github.com/arturoeanton/go-r2lang/compare/v0.0.1...v0.1.0
