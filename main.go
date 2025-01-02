package main

import (
	"fmt"
	"github.com/arturoeanton/go-r2lang/r2lang"
	"os"
)

func main() {
	filename := ""
	if len(os.Args) > 1 {
		filename = os.Args[1]
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
