package r2lang

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/arturoeanton/go-r2lang/pkg/r2libs"
)

func RunCode(filename string) {

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading the file %s: %v\n", filename, err)
		os.Exit(1)
	}
	code := string(data)

	env := r2core.NewEnvironment()
	env.Set("true", true)
	env.Set("false", false)
	env.Set("nil", nil)
	env.Set("null", nil)
	env.Dir = filepath.Dir(filename)

	// Registrar otras librerías si las tienes:
	r2libs.RegisterLib(env)
	r2libs.RegisterStd(env)
	r2libs.RegisterIO(env)
	r2libs.RegisterHTTPClient(env)
	r2libs.RegisterRequests(env)
	r2libs.RegisterString(env)
	r2libs.RegisterMath(env)
	r2libs.RegisterRand(env)
	r2libs.RegisterTest(env)
	r2libs.RegisterHTTP(env)
	r2libs.RegisterPrint(env)
	r2libs.RegisterOS(env)
	r2libs.RegisterHack(env)
	r2libs.RegisterConcurrency(env)
	r2libs.RegisterCollections(env)
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
	parser := r2core.NewParser(code)
	env.Run(parser)
}
