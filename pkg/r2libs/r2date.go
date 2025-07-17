package r2libs

import (
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterDate(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"Date": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) == 0 {
				return createDateObject()
			}
			if len(args) == 1 {
				if str, ok := args[0].(string); ok {
					t, err := time.Parse(time.RFC3339, str)
					if err != nil {
						t, err = time.Parse("2006-01-02", str)
						if err != nil {
							return nil
						}
					}
					return &r2core.DateValue{Time: t}
				}
				if ts, ok := args[0].(float64); ok {
					return &r2core.DateValue{Time: time.Unix(int64(ts/1000), int64(ts)%1000*1000000)}
				}
			}
			return createDateFromArgs(args...)
		}),
		"format": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				return nil
			}
			date, ok := args[0].(*r2core.DateValue)
			if !ok {
				return nil
			}
			format, ok := args[1].(string)
			if !ok {
				return nil
			}
			goFormat := r2core.ConvertToGoFormat(format)
			return date.Time.Format(goFormat)
		}),
	}

	RegisterModule(env, "date", functions)
}

func createDateFromArgs(args ...interface{}) *r2core.DateValue {
	if len(args) < 1 {
		return &r2core.DateValue{Time: time.Now()}
	}

	year := int(toFloat(args[0]))
	month := time.January
	day := 1
	hour, minute, second := 0, 0, 0

	if len(args) > 1 {
		month = time.Month(int(toFloat(args[1])) + 1)
	}
	if len(args) > 2 {
		day = int(toFloat(args[2]))
	}
	if len(args) > 3 {
		hour = int(toFloat(args[3]))
	}
	if len(args) > 4 {
		minute = int(toFloat(args[4]))
	}
	if len(args) > 5 {
		second = int(toFloat(args[5]))
	}

	return &r2core.DateValue{Time: time.Date(year, month, day, hour, minute, second, 0, time.Local)}
}

func createDateObject() map[string]interface{} {
	obj := make(map[string]interface{})

	// Constructor function for creating new DateValue instances
	obj["create"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) == 0 {
			return &r2core.DateValue{Time: time.Now()}
		}
		if len(args) == 1 {
			if str, ok := args[0].(string); ok {
				t, err := time.Parse(time.RFC3339, str)
				if err != nil {
					t, err = time.Parse("2006-01-02", str)
					if err != nil {
						return nil
					}
				}
				return &r2core.DateValue{Time: t}
			}
			if ts, ok := args[0].(float64); ok {
				return &r2core.DateValue{Time: time.Unix(int64(ts/1000), int64(ts)%1000*1000000)}
			}
		}
		return createDateFromArgs(args...)
	})

	obj["now"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return float64(time.Now().UnixMilli())
	})

	obj["parse"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			return nil
		}
		str, ok := args[0].(string)
		if !ok {
			return nil
		}
		t, err := time.Parse(time.RFC3339, str)
		if err != nil {
			t, err = time.Parse("2006-01-02", str)
			if err != nil {
				return nil
			}
		}
		return &r2core.DateValue{Time: t}
	})

	obj["UTC"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return createDateFromUTCArgs(args...)
	})

	obj["getTime"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UnixMilli())
	})

	obj["getFullYear"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Year())
	})

	obj["getMonth"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Month() - 1)
	})

	obj["getDate"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Day())
	})

	obj["getDay"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Weekday())
	})

	obj["getHours"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Hour())
	})

	obj["getMinutes"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Minute())
	})

	obj["getSeconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Second())
	})

	obj["getMilliseconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.Nanosecond() / 1000000)
	})

	obj["getTimezoneOffset"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		_, offset := date.Time.Zone()
		return float64(-offset / 60)
	})

	obj["setFullYear"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		year := int(toFloat(args[1]))
		date.Time = time.Date(year, date.Time.Month(), date.Time.Day(),
			date.Time.Hour(), date.Time.Minute(), date.Time.Second(),
			date.Time.Nanosecond(), date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["setMonth"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		month := time.Month(int(toFloat(args[1])) + 1)
		date.Time = time.Date(date.Time.Year(), month, date.Time.Day(),
			date.Time.Hour(), date.Time.Minute(), date.Time.Second(),
			date.Time.Nanosecond(), date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["setDate"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		day := int(toFloat(args[1]))
		date.Time = time.Date(date.Time.Year(), date.Time.Month(), day,
			date.Time.Hour(), date.Time.Minute(), date.Time.Second(),
			date.Time.Nanosecond(), date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["setHours"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		hour := int(toFloat(args[1]))
		date.Time = time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(),
			hour, date.Time.Minute(), date.Time.Second(),
			date.Time.Nanosecond(), date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["setMinutes"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		minute := int(toFloat(args[1]))
		date.Time = time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(),
			date.Time.Hour(), minute, date.Time.Second(),
			date.Time.Nanosecond(), date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["setSeconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		second := int(toFloat(args[1]))
		date.Time = time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(),
			date.Time.Hour(), date.Time.Minute(), second,
			date.Time.Nanosecond(), date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["setMilliseconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		ms := int(toFloat(args[1]))
		date.Time = time.Date(date.Time.Year(), date.Time.Month(), date.Time.Day(),
			date.Time.Hour(), date.Time.Minute(), date.Time.Second(),
			ms*1000000, date.Time.Location())
		return float64(date.Time.UnixMilli())
	})

	obj["toISOString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.UTC().Format(time.RFC3339)
	})

	obj["toDateString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.Format("Mon Jan 02 2006")
	})

	obj["toTimeString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.Format("15:04:05 MST")
	})

	obj["toString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.Format("Mon Jan 02 2006 15:04:05 MST")
	})

	obj["valueOf"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UnixMilli())
	})

	obj["getUTCFullYear"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Year())
	})

	obj["getUTCMonth"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Month() - 1)
	})

	obj["getUTCDate"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Day())
	})

	obj["getUTCDay"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Weekday())
	})

	obj["getUTCHours"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Hour())
	})

	obj["getUTCMinutes"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Minute())
	})

	obj["getUTCSeconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Second())
	})

	obj["getUTCMilliseconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return float64(date.Time.UTC().Nanosecond() / 1000000)
	})

	obj["toLocaleDateString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.Format("1/2/2006")
	})

	obj["toLocaleTimeString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.Format("3:04:05 PM")
	})

	obj["toLocaleString"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		return date.Time.Format("1/2/2006, 3:04:05 PM")
	})

	obj["addYears"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		years := int(toFloat(args[1]))
		newTime := date.Time.AddDate(years, 0, 0)
		return &r2core.DateValue{Time: newTime}
	})

	obj["addMonths"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		months := int(toFloat(args[1]))
		newTime := date.Time.AddDate(0, months, 0)
		return &r2core.DateValue{Time: newTime}
	})

	obj["addDays"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		days := int(toFloat(args[1]))
		newTime := date.Time.AddDate(0, 0, days)
		return &r2core.DateValue{Time: newTime}
	})

	obj["addHours"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		hours := int(toFloat(args[1]))
		newTime := date.Time.Add(time.Duration(hours) * time.Hour)
		return &r2core.DateValue{Time: newTime}
	})

	obj["addMinutes"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		minutes := int(toFloat(args[1]))
		newTime := date.Time.Add(time.Duration(minutes) * time.Minute)
		return &r2core.DateValue{Time: newTime}
	})

	obj["addSeconds"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		seconds := int(toFloat(args[1]))
		newTime := date.Time.Add(time.Duration(seconds) * time.Second)
		return &r2core.DateValue{Time: newTime}
	})

	obj["diff"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date1, ok1 := args[0].(*r2core.DateValue)
		date2, ok2 := args[1].(*r2core.DateValue)
		if !ok1 || !ok2 {
			return nil
		}
		unit := "milliseconds"
		if len(args) > 2 {
			if u, ok := args[2].(string); ok {
				unit = u
			}
		}

		diff := date1.Time.Sub(date2.Time)
		switch unit {
		case "years":
			return float64(diff.Hours() / 24 / 365.25)
		case "months":
			return float64(diff.Hours() / 24 / 30.44)
		case "days":
			return float64(diff.Hours() / 24)
		case "hours":
			return float64(diff.Hours())
		case "minutes":
			return float64(diff.Minutes())
		case "seconds":
			return float64(diff.Seconds())
		default:
			return float64(diff.Nanoseconds() / 1000000)
		}
	})

	obj["format"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		format, ok := args[1].(string)
		if !ok {
			return nil
		}
		goFormat := r2core.ConvertToGoFormat(format)
		return date.Time.Format(goFormat)
	})

	obj["timezone"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 4 {
			return nil
		}
		timezone, ok := args[0].(string)
		if !ok {
			return nil
		}
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil
		}
		year := int(toFloat(args[1]))
		month := time.Month(int(toFloat(args[2])))
		day := int(toFloat(args[3]))
		hour, minute, second := 0, 0, 0
		if len(args) > 4 {
			hour = int(toFloat(args[4]))
		}
		if len(args) > 5 {
			minute = int(toFloat(args[5]))
		}
		if len(args) > 6 {
			second = int(toFloat(args[6]))
		}
		return &r2core.DateValue{Time: time.Date(year, month, day, hour, minute, second, 0, loc)}
	})

	obj["toTimezone"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			return nil
		}
		date, ok := args[0].(*r2core.DateValue)
		if !ok {
			return nil
		}
		timezone, ok := args[1].(string)
		if !ok {
			return nil
		}
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			return nil
		}
		return &r2core.DateValue{Time: date.Time.In(loc)}
	})

	return obj
}

func createDateFromUTCArgs(args ...interface{}) *r2core.DateValue {
	if len(args) < 1 {
		return &r2core.DateValue{Time: time.Now().UTC()}
	}

	year := int(toFloat(args[0]))
	month := time.January
	day := 1
	hour, minute, second := 0, 0, 0

	if len(args) > 1 {
		month = time.Month(int(toFloat(args[1])) + 1)
	}
	if len(args) > 2 {
		day = int(toFloat(args[2]))
	}
	if len(args) > 3 {
		hour = int(toFloat(args[3]))
	}
	if len(args) > 4 {
		minute = int(toFloat(args[4]))
	}
	if len(args) > 5 {
		second = int(toFloat(args[5]))
	}

	return &r2core.DateValue{Time: time.Date(year, month, day, hour, minute, second, 0, time.UTC)}
}
