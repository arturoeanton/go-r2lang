package r2libs

import (
	"net"
	"net/url"
	"regexp"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

var validateEmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`)

func RegisterValidate(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"isEmail": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("validate.isEmail: se acepta 1 argumento (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("validate.isEmail: el argumento debe ser un string")
			}
			return validateEmailRegexp.MatchString(s)
		}),

		"isURL": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("validate.isURL: se acepta 1 argumento (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("validate.isURL: el argumento debe ser un string")
			}
			u, err := url.ParseRequestURI(s)
			if err != nil {
				return false
			}
			if u.Scheme == "" || u.Host == "" {
				return false
			}
			return true
		}),

		"isIP": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("validate.isIP: se acepta 1 argumento (str)")
			}
			s, ok := args[0].(string)
			if !ok {
				panic("validate.isIP: el argumento debe ser un string")
			}
			return net.ParseIP(s) != nil
		}),
	}

	RegisterModule(env, "validate", functions)
}
