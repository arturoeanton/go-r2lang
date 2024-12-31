package r2lang

import (
	"fmt"
	"strings"
	"time"
	// Podrías importar alguna librería de colores (ej. "github.com/fatih/color") si quieres
)

// r2print.go: Funciones de impresión avanzadas para R2

func RegisterPrint(env *Environment) {
	// 1) printRepeat(str, count)
	env.Set("printRepeat", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("printRepeat necesita (str, count)")
		}
		s, ok1 := args[0].(string)
		count := int(toFloat(args[1]))
		if !ok1 {
			panic("printRepeat: primer arg debe ser string")
		}
		for i := 0; i < count; i++ {
			fmt.Print(s)
		}
		fmt.Println() // salto de línea
		return nil
	}))

	// 2) printBox(str, width)
	//   Crea una "caja" ASCII rodeando el texto
	env.Set("printBox", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("printBox necesita (str, width)")
		}
		text, ok1 := args[0].(string)
		width := int(toFloat(args[1]))
		if !ok1 {
			panic("printBox: primer arg debe ser string")
		}
		if width < len(text)+2 {
			width = len(text) + 2
		}
		// Cabecera
		fmt.Println("+" + strings.Repeat("-", width) + "+")
		// Texto centrado
		space := width - len(text)
		leftPad := space / 2
		rightPad := space - leftPad
		fmt.Printf("|%s%s%s|\n", strings.Repeat(" ", leftPad), text, strings.Repeat(" ", rightPad))
		// Pie
		fmt.Println("+" + strings.Repeat("-", width) + "+")
		return nil
	}))

	// 3) debugInspect(value)
	//    Imprime detalles sobre un valor (Go).
	env.Set("debugInspect", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("debugInspect necesita (value)")
		}
		val := args[0]
		fmt.Printf("[debugInspect] Valor = %v (tipo=%T)\n", val, val)
		return nil
	}))

	// 4) printColor(str, colorName)
	//    Ejemplo simple de colores ANSI (sin librería externa).
	//    Soporta: "red", "green", "yellow", "blue", "reset"
	env.Set("printColor", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("printColor(str, colorName)")
		}
		txt, ok1 := args[0].(string)
		colorName, ok2 := args[1].(string)
		if !(ok1 && ok2) {
			panic("printColor: (string, string)")
		}
		colorCode := ""
		switch strings.ToLower(colorName) {
		case "red":
			colorCode = "\033[31m"
		case "green":
			colorCode = "\033[32m"
		case "yellow":
			colorCode = "\033[33m"
		case "blue":
			colorCode = "\033[34m"
		case "reset":
			colorCode = "\033[0m"
		default:
			colorCode = "\033[0m" // reset
		}
		// Imprimir con color, y reset al final
		fmt.Print(colorCode, txt, "\033[0m\n")
		return nil
	}))

	// 5) printProgress(label, totalSteps, stepDelayMs)
	//    Muestra un “progress bar” en la terminal (muy simple).
	//    Ej: printProgress("Loading", 10, 200)
	env.Set("printProgress", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("printProgress(label, totalSteps, stepDelayMs)")
		}
		label, ok1 := args[0].(string)
		total := int(toFloat(args[1]))
		delayMs := int(toFloat(args[2]))
		if !ok1 {
			panic("printProgress: label debe ser string")
		}
		// Cada step, imprime algo
		for i := 0; i <= total; i++ {
			pct := float64(i) / float64(total) * 100.0
			bar := strings.Repeat("#", i) + strings.Repeat(" ", total-i)
			fmt.Printf("\r%s [%s] %.0f%%", label, bar, pct)
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
		fmt.Println() // salto línea final
		return nil
	}))

	// 6) printTable(arrayOfArrays)
	//    Imprime una tabla. Ejemplo de formateo.
	//    asume un array con filas, cada fila => array de celdas.
	env.Set("printTable", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("printTable necesita (arrayOfArrays)")
		}
		rows, ok := args[0].([]interface{})
		if !ok {
			panic("printTable: primer arg debe ser array (cada elem => array celdas)")
		}
		// Calcular anchos máximos por columna
		// 1) averiguar cuántas columnas máximo
		maxCols := 0
		tableData := make([][]string, len(rows))
		for i, row := range rows {
			rowArr, ok2 := row.([]interface{})
			if !ok2 {
				panic("printTable: cada fila debe ser un array")
			}
			tableData[i] = make([]string, len(rowArr))
			if len(rowArr) > maxCols {
				maxCols = len(rowArr)
			}
			// convertir a string
			for j, cell := range rowArr {
				tableData[i][j] = fmt.Sprint(cell)
			}
		}
		// 2) ancho máximo por columna
		colWidths := make([]int, maxCols)
		for i := 0; i < len(tableData); i++ {
			for j := 0; j < len(tableData[i]); j++ {
				cellLen := len(tableData[i][j])
				if cellLen > colWidths[j] {
					colWidths[j] = cellLen
				}
			}
		}
		// 3) imprimir
		for i := 0; i < len(tableData); i++ {
			row := tableData[i]
			for j := 0; j < len(row); j++ {
				fmt.Printf("%-*s ", colWidths[j], row[j])
			}
			fmt.Println()
		}
		return nil
	}))

	// 7) printAlign(str, align, width)
	//    align => "left", "right", "center"
	env.Set("printAlign", BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("printAlign(str, align, width)")
		}
		s, ok1 := args[0].(string)
		alignOpt, ok2 := args[1].(string)
		width := int(toFloat(args[2]))
		if !(ok1 && ok2) {
			panic("printAlign: (string, align, width)")
		}
		if width < len(s) {
			width = len(s)
		}
		space := width - len(s)
		switch strings.ToLower(alignOpt) {
		case "left":
			// s + spaces
			fmt.Println(s + strings.Repeat(" ", space))
		case "right":
			// spaces + s
			fmt.Println(strings.Repeat(" ", space) + s)
		case "center":
			leftPad := space / 2
			rightPad := space - leftPad
			fmt.Println(strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad))
		default:
			panic("printAlign: align debe ser 'left','right' o 'center'")
		}
		return nil
	}))
}
