package r2libs

import (
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterDate(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"Date": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			return createDateObject()
		}),
	}

	RegisterModule(env, "date", functions)
}

func createDateObject() map[string]interface{} {
	obj := make(map[string]interface{})
	obj["now"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return &r2core.DateValue{Time: time.Now()}
	})
	obj["today"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return &r2core.DateValue{Time: time.Now().Truncate(24 * time.Hour)}
	})
	obj["create"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			return nil // Error: requires at least 3 arguments
		}
		year := int(args[0].(float64))
		month := time.Month(int(args[1].(float64)))
		day := int(args[2].(float64))
		hour, minute, second := 0, 0, 0
		if len(args) > 3 {
			hour = int(args[3].(float64))
		}
		if len(args) > 4 {
			minute = int(args[4].(float64))
		}
		if len(args) > 5 {
			second = int(args[5].(float64))
		}
		return &r2core.DateValue{Time: time.Date(year, month, day, hour, minute, second, 0, time.Local)}
	})
	obj["getYear"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil // Error: requires a Date object
		}
		return float64(date.Time.Year())
	})
	obj["getMonth"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil // Error: requires a Date object
		}
		return float64(date.Time.Month())
	})
	obj["getDay"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil // Error: requires a Date object
		}
		return float64(date.Time.Day())
	})
	obj["format"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil // Error: requires a Date object and a format string
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil // Error: requires a Date object
		}
		format, ok := args[1].(string)
		if !ok {
			return nil // Error: requires a string as format
		}
		goFormat := r2core.ConvertToGoFormat(format)
		return date.Time.Format(goFormat)
	})
	obj["timezone"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 4 {
			return nil // Error: requires at least 4 arguments
		}
		timezone, ok := args[0].(string)
		if !ok {
			return nil // Error: requires a string as timezone
		}
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil // Error: invalid timezone
		}
		year := int(args[1].(float64))
		month := time.Month(int(args[2].(float64)))
		day := int(args[3].(float64))
		hour, minute, second := 0, 0, 0
		if len(args) > 4 {
			hour = int(args[4].(float64))
		}
		if len(args) > 5 {
			minute = int(args[5].(float64))
		}
		if len(args) > 6 {
			second = int(args[6].(float64))
		}
		return &r2core.DateValue{Time: time.Date(year, month, day, hour, minute, second, 0, loc)}
	})
	obj["toTimezone"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil // Error: requires a Date object and a timezone string
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil // Error: requires a Date object
		}
		timezone, ok := args[1].(string)
		if !ok {
			return nil // Error: requires a string as timezone
		}
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil // Error: invalid timezone
		}
		return &r2core.DateValue{Time: date.Time.In(loc)}
	})
	return obj
}
