
# Análisis de Fallos en Tests de R2Core

## Introducción

Este documento detalla los fallos encontrados al ejecutar los tests del paquete `r2core`, específicamente en los ficheros `lexer.go` y `lexer_test.go`. Se proporciona un análisis de cada test fallido y se propone una solución.

## Resumen de Tests Fallidos

Los siguientes tests han fallado:

1.  `TestLexer_Comments/single_line_comment`
2.  `TestLexer_LineAndColumn`
3.  `TestLexer_WhitespaceHandling`
4.  `TestLexer_CompleteExpression`
5.  `TestLexer_BDDSyntax`

## Análisis Detallado de Fallos y Soluciones

### 1. `TestLexer_Comments/single_line_comment`

*   **Error:** El lexer no está manejando correctamente los comentarios de una sola línea. En lugar de ignorar el comentario y continuar con el siguiente token, parece que está interpretando el contenido del comentario como parte del código.

*   **Causa:** La lógica para saltar los comentarios de una sola línea en `NextToken` es incorrecta. El bucle que ignora los caracteres del comentario no consume el carácter de nueva línea (`
`), lo que provoca que el lexer lo trate como un token separado en la siguiente iteración.

*   **Solución:** Modificar el bucle que maneja los comentarios de una sola línea para que también consuma el carácter de nueva línea.

    ```go
    // en NextToken dentro de Lexer
    if l.pos+1 < l.length && l.input[l.pos+1] == '/' {
        // comentario de línea
        l.pos += 2
        for l.pos < l.length && l.input[l.pos] != '
' {
            l.nextch()
        }
        // Consumir el '
' para que no sea tratado como un token
        if l.pos < l.length && l.input[l.pos] == '
' {
            l.nextch()
        }
    }
    ```

### 2. `TestLexer_LineAndColumn`

*   **Error:** El test espera que el número de línea se incremente después de procesar un salto de línea, pero el lexer no lo hace correctamente.

*   **Causa:** El lexer incrementa el número de línea en la función `nextch()`, pero el token de nueva línea se devuelve *antes* de que se procese el siguiente carácter, por lo que el número de línea todavía no se ha actualizado.

*   **Solución:** Asegurarse de que el número de línea se actualiza *antes* de devolver el token de nueva línea.

    ```go
    // en NextToken dentro de Lexer
    if string(ch) == "
" {
        l.nextch()
        l.currentToken = Token{Type: TOKEN_SYMBOL, Value: "
", Line: l.line -1, Pos: l.pos, Col: l.col}
        return l.currentToken
    }
    ```

### 3. `TestLexer_WhitespaceHandling`

*   **Error:** El lexer no está manejando correctamente los espacios en blanco (espacios, tabuladores, saltos de línea).

*   **Causa:** Similar al problema con los comentarios, el manejo de los espacios en blanco no es consistente. El lexer debería consumir todos los espacios en blanco consecutivos y luego devolver el siguiente token válido.

*   **Solución:** La función `skipWhitespace` al principio de `NextToken` debería manejar todos los tipos de espacios en blanco de forma correcta. El problema parece estar en cómo se tratan los saltos de línea. La solución del punto anterior debería ayudar a resolver este problema también.

### 4. `TestLexer_CompleteExpression` y `TestLexer_BDDSyntax`

*   **Error:** Ambos tests fallan porque el lexer no está "tokenizando" correctamente expresiones complejas y la sintaxis BDD. Se esperan ciertos tipos de tokens y valores, pero se reciben otros.

*   **Causa:** La causa raíz de estos fallos es una combinación de los problemas anteriores (manejo de comentarios, saltos de línea y espacios en blanco) y una lógica de "tokenización" demasiado simplista para algunos casos. Por ejemplo, el lexer no parece diferenciar correctamente entre palabras clave como `if` o `return` y los identificadores cuando están pegados a otros símbolos.

*   **Solución:**

    1.  **Arreglar los problemas fundamentales:** Primero, es crucial arreglar los problemas con los comentarios, saltos de línea y espacios en blanco como se describió anteriormente.
    2.  **Mejorar la lógica de identificación de tokens:**
        *   Asegurarse de que al identificar un `TOKEN_IDENT`, no esté seguido inmediatamente por un carácter que podría ser parte de otro token (por ejemplo, un `(` después de `if`).
        *   Revisar el orden de las comprobaciones en `NextToken`. Las comprobaciones de símbolos de varios caracteres (como `==` o `<=`) deben realizarse antes que las de un solo carácter.
        *   La lógica para las palabras clave (`if`, `return`, etc.) debe ser robusta y manejar correctamente los casos en los que son parte de un identificador más grande (aunque esto no debería pasar si la lógica de `isLetter` y `isDigit` es correcta).

## Conclusión

Los fallos en los tests del lexer se deben principalmente a un manejo incorrecto de los saltos de línea, comentarios y espacios en blanco. Al solucionar estos problemas fundamentales, la mayoría de los tests deberían pasar. Los tests más complejos, como `TestLexer_CompleteExpression` y `TestLexer_BDDSyntax`, pueden requerir ajustes adicionales en la lógica de identificación de tokens para garantizar que todos los casos se manejen correctamente.
