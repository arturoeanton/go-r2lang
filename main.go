package main

import (
	"fmt"
	"os"

	"github.com/arturoeanton/go-r2lang/pkg/r2lang"
	"github.com/arturoeanton/go-r2lang/pkg/r2repl"
)

func main() {
	filename := ""
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "-repl" {
			outputFlag := true
			if len(os.Args) > 2 {
				if os.Args[2] == "-no-output" {
					outputFlag = false
				}
			}
			r2repl.Repl(outputFlag)
			os.Exit(0)
		}
		filename = cmd

	} else {
		// intentar main.r2
		if _, err := os.Stat("main.r2"); os.IsNotExist(err) {
			fmt.Println("Error: Debes pasar un archivo .r2 o tener main.r2 en el directorio actual.")
			os.Exit(1)
		}

		filename = "main.r2"
	}
	r2lang.RunCode(filename)
}
