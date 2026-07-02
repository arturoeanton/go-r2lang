# Embedding R2Lang en un programa Go

Este es el uso "normal" de R2Lang: un programa Go que embebe el interprete,
expone sus propias funciones/structs, y corre scripts `.r2` que las usan.

## Correr el ejemplo

```bash
go run ./examples/embedding
```

Salida esperada:

```
=== Embedding R2Lang: native.* ===
native.callFunc("add", 3, 4) = 7
Hello, R2Lang!
```

## El patron

1. **Registrar del lado Go, antes de correr cualquier script.** `r2libs.RegisterNativeFunc(nombre, fn)` expone una funcion Go; `r2libs.RegisterNativeStruct(nombre, constructor)` expone un struct.

   ```go
   r2libs.RegisterNativeFunc("add", Add)
   r2libs.RegisterNativeStruct("Greeter", func() interface{} { return &Greeter{} })
   ```

2. **Armar el `Environment` y registrar `RegisterGoInterOp`** (expone el modulo `native`) junto con el resto de la stdlib que el script necesite:

   ```go
   env := r2core.NewEnvironment()
   r2libs.RegisterStd(env)
   r2libs.RegisterGoInterOp(env)
   ```

3. **Consumir desde el script** via el namespace `native`:

   ```r2
   let sum = native.callFunc("add", 3, 4)
   let g = native.new("Greeter")
   native.setField(g, "Name", "R2Lang")
   std.print(native.callMethod(g, "Hello"))
   ```

## Funciones disponibles bajo `native`

| Funcion | Uso |
|---|---|
| `native.callFunc(nombre, ...args)` | Llama una funcion registrada con `RegisterNativeFunc` |
| `native.new(nombre)` | Crea una instancia registrada con `RegisterNativeStruct` |
| `native.setField(obj, campo, valor)` | Setea un campo exportado del struct |
| `native.getField(obj, campo)` | Lee un campo exportado del struct |
| `native.callMethod(obj, metodo, ...args)` | Llama un metodo exportado del struct |

## Notas importantes

- **El registro (`RegisterNativeFunc`/`RegisterNativeStruct`) solo puede hacerse desde Go**, nunca desde un script `.r2` — no hay forma de representar un `reflect.Value` en R2Lang. `native.callFunc`/`native.new` fallan con un mensaje claro si el nombre no fue registrado.
- **Los numeros en R2Lang son siempre `float64`.** Si tu funcion/struct Go espera `int`, `int64`, `float32`, etc, `native.callFunc`/`native.setField` convierten automaticamente — no hace falta que tu codigo Go acepte `float64` en todos lados.
- Solo campos/metodos **exportados** (que empiezan con mayuscula) son accesibles.

## Restringir os.Command/exec para scripts no confiables

Por diseno, `os.Command`/`os.execCmd`/`os.runProcess`/`os.execWithTimeout`/
`os.execWithEnv` le dan al script acceso de shell sin restricciones (es una
funcionalidad deliberada, no un bug). Si tu programa anfitrion va a correr
scripts que no controla del todo, llama a `r2libs.SetCommandPolicy(...)`
ANTES de correr el script:

```go
r2libs.SetCommandPolicy(r2libs.CommandPolicy{
    // Restringe os.Command(...).run() a un allowlist de ejecutables (no
    // aplica shell, asi que el allowlist es seguro: no hay forma de
    // encadenar comandos via ";"/"&&"/"|").
    AllowedCommands: map[string]bool{"git": true, "echo": true},

    // Bloquea por completo execCmd/runProcess/execWithTimeout/execWithEnv,
    // que corren via "sh -c" (un shell completo no se puede restringir de
    // forma segura con un allowlist parcial). Para exponer una operacion
    // especifica del host a un script restringido, usa
    // RegisterNativeFunc en vez del shell.
    DisableShell: true,
})
```

El valor por defecto (`CommandPolicy{}`, sin llamar a `SetCommandPolicy`)
es sin restricciones — el comportamiento historico del proyecto no cambia
si no optás explicitamente por esto.
