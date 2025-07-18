package r2libs

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterCSV(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"parse": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("CSV.parse needs (csvString)")
			}
			csvStr, ok := args[0].(string)
			if !ok {
				panic("CSV.parse: first argument must be string")
			}

			delimiter := ","
			if len(args) > 1 {
				if d, ok := args[1].(string); ok {
					delimiter = d
				}
			}

			hasHeader := true
			if len(args) > 2 {
				if h, ok := args[2].(bool); ok {
					hasHeader = h
				}
			}

			reader := csv.NewReader(strings.NewReader(csvStr))
			reader.Comma = rune(delimiter[0])

			records, err := reader.ReadAll()
			if err != nil {
				panic(fmt.Sprintf("CSV.parse: error parsing CSV: %v", err))
			}

			if len(records) == 0 {
				return []interface{}{}
			}

			var result []interface{}

			if hasHeader && len(records) > 1 {
				headers := records[0]
				for i := 1; i < len(records); i++ {
					row := make(map[string]interface{})
					for j, value := range records[i] {
						if j < len(headers) {
							row[headers[j]] = convertCSVValue(value)
						}
					}
					result = append(result, row)
				}
			} else {
				startIndex := 0
				if hasHeader {
					startIndex = 1
				}

				for i := startIndex; i < len(records); i++ {
					var row []interface{}
					for _, value := range records[i] {
						row = append(row, convertCSVValue(value))
					}
					result = append(result, row)
				}
			}

			return result
		}),

		"stringify": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("CSV.stringify needs (array)")
			}

			data, ok := args[0].([]interface{})
			if !ok {
				panic("CSV.stringify: first argument must be array")
			}

			delimiter := ","
			if len(args) > 1 {
				if d, ok := args[1].(string); ok {
					delimiter = d
				}
			}

			includeHeaders := true
			if len(args) > 2 {
				if h, ok := args[2].(bool); ok {
					includeHeaders = h
				}
			}

			if len(data) == 0 {
				return ""
			}

			var result strings.Builder
			writer := csv.NewWriter(&result)
			writer.Comma = rune(delimiter[0])

			// Check if data is array of objects or array of arrays
			if firstRow, ok := data[0].(map[string]interface{}); ok {
				// Array of objects
				var headers []string
				for key := range firstRow {
					headers = append(headers, key)
				}

				if includeHeaders {
					writer.Write(headers)
				}

				for _, item := range data {
					if rowMap, ok := item.(map[string]interface{}); ok {
						var row []string
						for _, header := range headers {
							if value, exists := rowMap[header]; exists {
								row = append(row, fmt.Sprintf("%v", value))
							} else {
								row = append(row, "")
							}
						}
						writer.Write(row)
					}
				}
			} else {
				// Array of arrays
				for _, item := range data {
					if rowArray, ok := item.([]interface{}); ok {
						var row []string
						for _, value := range rowArray {
							row = append(row, fmt.Sprintf("%v", value))
						}
						writer.Write(row)
					}
				}
			}

			writer.Flush()
			if err := writer.Error(); err != nil {
				panic(fmt.Sprintf("CSV.stringify: error writing CSV: %v", err))
			}

			return result.String()
		}),

		"readFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("CSV.readFile needs (filePath)")
			}
			filePath, ok := args[0].(string)
			if !ok {
				panic("CSV.readFile: filePath must be string")
			}

			delimiter := ","
			if len(args) > 1 {
				if d, ok := args[1].(string); ok {
					delimiter = d
				}
			}

			hasHeader := true
			if len(args) > 2 {
				if h, ok := args[2].(bool); ok {
					hasHeader = h
				}
			}

			file, err := os.Open(filePath)
			if err != nil {
				panic(fmt.Sprintf("CSV.readFile: error opening file '%s': %v", filePath, err))
			}
			defer file.Close()

			reader := csv.NewReader(file)
			reader.Comma = rune(delimiter[0])

			records, err := reader.ReadAll()
			if err != nil {
				panic(fmt.Sprintf("CSV.readFile: error reading CSV: %v", err))
			}

			if len(records) == 0 {
				return []interface{}{}
			}

			var result []interface{}

			if hasHeader && len(records) > 1 {
				headers := records[0]
				for i := 1; i < len(records); i++ {
					row := make(map[string]interface{})
					for j, value := range records[i] {
						if j < len(headers) {
							row[headers[j]] = convertCSVValue(value)
						}
					}
					result = append(result, row)
				}
			} else {
				startIndex := 0
				if hasHeader {
					startIndex = 1
				}

				for i := startIndex; i < len(records); i++ {
					var row []interface{}
					for _, value := range records[i] {
						row = append(row, convertCSVValue(value))
					}
					result = append(result, row)
				}
			}

			return result
		}),

		"writeFile": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("CSV.writeFile needs (filePath, data)")
			}
			filePath, ok1 := args[0].(string)
			data, ok2 := args[1].([]interface{})
			if !ok1 || !ok2 {
				panic("CSV.writeFile: arguments must be (string, array)")
			}

			delimiter := ","
			if len(args) > 2 {
				if d, ok := args[2].(string); ok {
					delimiter = d
				}
			}

			includeHeaders := true
			if len(args) > 3 {
				if h, ok := args[3].(bool); ok {
					includeHeaders = h
				}
			}

			file, err := os.Create(filePath)
			if err != nil {
				panic(fmt.Sprintf("CSV.writeFile: error creating file '%s': %v", filePath, err))
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			writer.Comma = rune(delimiter[0])
			defer writer.Flush()

			if len(data) == 0 {
				return nil
			}

			// Check if data is array of objects or array of arrays
			if firstRow, ok := data[0].(map[string]interface{}); ok {
				// Array of objects
				var headers []string
				for key := range firstRow {
					headers = append(headers, key)
				}

				if includeHeaders {
					writer.Write(headers)
				}

				for _, item := range data {
					if rowMap, ok := item.(map[string]interface{}); ok {
						var row []string
						for _, header := range headers {
							if value, exists := rowMap[header]; exists {
								row = append(row, fmt.Sprintf("%v", value))
							} else {
								row = append(row, "")
							}
						}
						writer.Write(row)
					}
				}
			} else {
				// Array of arrays
				for _, item := range data {
					if rowArray, ok := item.([]interface{}); ok {
						var row []string
						for _, value := range rowArray {
							row = append(row, fmt.Sprintf("%v", value))
						}
						writer.Write(row)
					}
				}
			}

			if err := writer.Error(); err != nil {
				panic(fmt.Sprintf("CSV.writeFile: error writing CSV: %v", err))
			}

			return nil
		}),

		"getHeaders": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("CSV.getHeaders needs (csvData)")
			}
			data, ok := args[0].([]interface{})
			if !ok {
				panic("CSV.getHeaders: first argument must be array")
			}

			if len(data) == 0 {
				return []interface{}{}
			}

			if firstRow, ok := data[0].(map[string]interface{}); ok {
				var headers []interface{}
				for key := range firstRow {
					headers = append(headers, key)
				}
				return headers
			}

			return []interface{}{}
		}),

		"getColumn": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("CSV.getColumn needs (csvData, columnName)")
			}
			data, ok1 := args[0].([]interface{})
			columnName, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("CSV.getColumn: arguments must be (array, string)")
			}

			var result []interface{}
			for _, row := range data {
				if rowMap, ok := row.(map[string]interface{}); ok {
					if value, exists := rowMap[columnName]; exists {
						result = append(result, value)
					} else {
						result = append(result, nil)
					}
				}
			}

			return result
		}),

		"filter": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("CSV.filter needs (csvData, filterFunction)")
			}
			data, ok := args[0].([]interface{})
			fn, okFn := args[1].(*r2core.UserFunction)
			if !ok || !okFn {
				panic("CSV.filter: arguments must be (array, function)")
			}

			var result []interface{}
			for _, row := range data {
				if toBool(fn.Call(row)) {
					result = append(result, row)
				}
			}

			return result
		}),

		"map": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("CSV.map needs (csvData, mapFunction)")
			}
			data, ok := args[0].([]interface{})
			fn, okFn := args[1].(*r2core.UserFunction)
			if !ok || !okFn {
				panic("CSV.map: arguments must be (array, function)")
			}

			var result []interface{}
			for _, row := range data {
				result = append(result, fn.Call(row))
			}

			return result
		}),

		"sort": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("CSV.sort needs (csvData, columnName)")
			}
			data, ok1 := args[0].([]interface{})
			columnName, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("CSV.sort: arguments must be (array, string)")
			}

			ascending := true
			if len(args) > 2 {
				if asc, ok := args[2].(bool); ok {
					ascending = asc
				}
			}

			// Create a copy to avoid modifying original
			result := make([]interface{}, len(data))
			copy(result, data)

			// Simple bubble sort for CSV data
			for i := 0; i < len(result)-1; i++ {
				for j := 0; j < len(result)-i-1; j++ {
					row1, ok1 := result[j].(map[string]interface{})
					row2, ok2 := result[j+1].(map[string]interface{})

					if !ok1 || !ok2 {
						continue
					}

					val1, exists1 := row1[columnName]
					val2, exists2 := row2[columnName]

					if !exists1 || !exists2 {
						continue
					}

					shouldSwap := false

					// Try to compare as numbers first
					if num1, ok1 := val1.(float64); ok1 {
						if num2, ok2 := val2.(float64); ok2 {
							if ascending {
								shouldSwap = num1 > num2
							} else {
								shouldSwap = num1 < num2
							}
						}
					} else {
						// Compare as strings
						str1 := fmt.Sprintf("%v", val1)
						str2 := fmt.Sprintf("%v", val2)
						if ascending {
							shouldSwap = str1 > str2
						} else {
							shouldSwap = str1 < str2
						}
					}

					if shouldSwap {
						result[j], result[j+1] = result[j+1], result[j]
					}
				}
			}

			return result
		}),

		"groupBy": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("CSV.groupBy needs (csvData, columnName)")
			}
			data, ok1 := args[0].([]interface{})
			columnName, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("CSV.groupBy: arguments must be (array, string)")
			}

			groups := make(map[string][]interface{})

			for _, row := range data {
				if rowMap, ok := row.(map[string]interface{}); ok {
					if value, exists := rowMap[columnName]; exists {
						key := fmt.Sprintf("%v", value)
						groups[key] = append(groups[key], row)
					}
				}
			}

			result := make(map[string]interface{})
			for key, group := range groups {
				result[key] = group
			}

			return result
		}),

		"aggregate": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("CSV.aggregate needs (csvData, columnName, operation)")
			}
			data, ok1 := args[0].([]interface{})
			columnName, ok2 := args[1].(string)
			operation, ok3 := args[2].(string)
			if !ok1 || !ok2 || !ok3 {
				panic("CSV.aggregate: arguments must be (array, string, string)")
			}

			var values []float64
			for _, row := range data {
				if rowMap, ok := row.(map[string]interface{}); ok {
					if value, exists := rowMap[columnName]; exists {
						if num, ok := value.(float64); ok {
							values = append(values, num)
						} else if str, ok := value.(string); ok {
							if parsed, err := strconv.ParseFloat(str, 64); err == nil {
								values = append(values, parsed)
							}
						}
					}
				}
			}

			if len(values) == 0 {
				return nil
			}

			switch operation {
			case "sum":
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				return sum
			case "avg", "average":
				sum := 0.0
				for _, v := range values {
					sum += v
				}
				return sum / float64(len(values))
			case "min":
				min := values[0]
				for _, v := range values {
					if v < min {
						min = v
					}
				}
				return min
			case "max":
				max := values[0]
				for _, v := range values {
					if v > max {
						max = v
					}
				}
				return max
			case "count":
				return float64(len(values))
			default:
				panic(fmt.Sprintf("CSV.aggregate: unsupported operation '%s'", operation))
			}
		}),

		"validate": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("CSV.validate needs (csvString)")
			}
			csvStr, ok := args[0].(string)
			if !ok {
				panic("CSV.validate: first argument must be string")
			}

			delimiter := ","
			if len(args) > 1 {
				if d, ok := args[1].(string); ok {
					delimiter = d
				}
			}

			reader := csv.NewReader(strings.NewReader(csvStr))
			reader.Comma = rune(delimiter[0])

			var errors []interface{}
			lineNum := 0
			var expectedColumns int

			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				lineNum++

				if err != nil {
					errors = append(errors, map[string]interface{}{
						"line":  float64(lineNum),
						"error": err.Error(),
					})
					continue
				}

				if lineNum == 1 {
					expectedColumns = len(record)
				} else if len(record) != expectedColumns {
					errors = append(errors, map[string]interface{}{
						"line":     float64(lineNum),
						"error":    "Column count mismatch",
						"expected": float64(expectedColumns),
						"actual":   float64(len(record)),
					})
				}
			}

			return map[string]interface{}{
				"valid":  len(errors) == 0,
				"errors": errors,
			}
		}),
	}

	RegisterModule(env, "csv", functions)
}

// Helper function to convert CSV values to appropriate types
func convertCSVValue(value string) interface{} {
	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	// Try to parse as number
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num
	}

	// Try to parse as boolean
	if value == "true" || value == "TRUE" {
		return true
	}
	if value == "false" || value == "FALSE" {
		return false
	}

	// Return as string
	return value
}
