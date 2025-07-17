package r2libs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func RegisterJSON(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"parse": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.parse needs (jsonString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.parse: first argument must be string")
			}

			var result interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.parse: error parsing JSON: %v", err))
			}

			return convertJSONToR2(result)
		}),

		"stringify": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.stringify needs (value)")
			}

			value := convertR2ToJSON(args[0])

			indent := ""
			if len(args) > 1 {
				if replacer, ok := args[1].([]interface{}); ok {
					// Handle replacer array (simplified implementation)
					_ = replacer
				}
			}

			if len(args) > 2 {
				if space, ok := args[2].(float64); ok {
					indent = strings.Repeat(" ", int(space))
				} else if space, ok := args[2].(string); ok {
					indent = space
				}
			}

			var result []byte
			var err error

			if indent != "" {
				result, err = json.MarshalIndent(value, "", indent)
			} else {
				result, err = json.Marshal(value)
			}

			if err != nil {
				panic(fmt.Sprintf("JSON.stringify: error stringifying: %v", err))
			}

			return string(result)
		}),

		"parseArray": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.parseArray needs (jsonArrayString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.parseArray: first argument must be string")
			}

			var result []interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.parseArray: error parsing JSON array: %v", err))
			}

			return convertJSONArrayToR2(result)
		}),

		"parseObject": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.parseObject needs (jsonObjectString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.parseObject: first argument must be string")
			}

			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.parseObject: error parsing JSON object: %v", err))
			}

			return convertJSONObjectToR2(result)
		}),

		"validate": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.validate needs (jsonString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.validate: first argument must be string")
			}

			var result interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			return err == nil
		}),

		"getKeys": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.getKeys needs (jsonObjectString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.getKeys: first argument must be string")
			}

			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.getKeys: error parsing JSON object: %v", err))
			}

			var keys []interface{}
			for key := range result {
				keys = append(keys, key)
			}

			return keys
		}),

		"getValue": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JSON.getValue needs (jsonObjectString, key)")
			}
			jsonStr, ok1 := args[0].(string)
			key, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JSON.getValue: arguments must be strings")
			}

			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.getValue: error parsing JSON object: %v", err))
			}

			if value, exists := result[key]; exists {
				return convertJSONToR2(value)
			}

			return nil
		}),

		"setValue": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("JSON.setValue needs (jsonObjectString, key, value)")
			}
			jsonStr, ok1 := args[0].(string)
			key, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JSON.setValue: first two arguments must be strings")
			}

			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.setValue: error parsing JSON object: %v", err))
			}

			result[key] = convertR2ToJSON(args[2])

			updated, err := json.Marshal(result)
			if err != nil {
				panic(fmt.Sprintf("JSON.setValue: error marshaling result: %v", err))
			}

			return string(updated)
		}),

		"deleteKey": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JSON.deleteKey needs (jsonObjectString, key)")
			}
			jsonStr, ok1 := args[0].(string)
			key, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JSON.deleteKey: arguments must be strings")
			}

			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.deleteKey: error parsing JSON object: %v", err))
			}

			delete(result, key)

			updated, err := json.Marshal(result)
			if err != nil {
				panic(fmt.Sprintf("JSON.deleteKey: error marshaling result: %v", err))
			}

			return string(updated)
		}),

		"hasKey": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JSON.hasKey needs (jsonObjectString, key)")
			}
			jsonStr, ok1 := args[0].(string)
			key, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JSON.hasKey: arguments must be strings")
			}

			var result map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				panic(fmt.Sprintf("JSON.hasKey: error parsing JSON object: %v", err))
			}

			_, exists := result[key]
			return exists
		}),

		"merge": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JSON.merge needs at least 2 JSON objects")
			}

			var merged map[string]interface{}

			for i, arg := range args {
				jsonStr, ok := arg.(string)
				if !ok {
					panic(fmt.Sprintf("JSON.merge: argument %d must be string", i))
				}

				var obj map[string]interface{}
				err := json.Unmarshal([]byte(jsonStr), &obj)
				if err != nil {
					panic(fmt.Sprintf("JSON.merge: error parsing JSON object %d: %v", i, err))
				}

				if i == 0 {
					merged = obj
				} else {
					for key, value := range obj {
						merged[key] = value
					}
				}
			}

			result, err := json.Marshal(merged)
			if err != nil {
				panic(fmt.Sprintf("JSON.merge: error marshaling result: %v", err))
			}

			return string(result)
		}),

		"deepMerge": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JSON.deepMerge needs at least 2 JSON objects")
			}

			var merged map[string]interface{}

			for i, arg := range args {
				jsonStr, ok := arg.(string)
				if !ok {
					panic(fmt.Sprintf("JSON.deepMerge: argument %d must be string", i))
				}

				var obj map[string]interface{}
				err := json.Unmarshal([]byte(jsonStr), &obj)
				if err != nil {
					panic(fmt.Sprintf("JSON.deepMerge: error parsing JSON object %d: %v", i, err))
				}

				if i == 0 {
					merged = obj
				} else {
					merged = deepMergeObjects(merged, obj)
				}
			}

			result, err := json.Marshal(merged)
			if err != nil {
				panic(fmt.Sprintf("JSON.deepMerge: error marshaling result: %v", err))
			}

			return string(result)
		}),

		"flatten": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.flatten needs (jsonObjectString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.flatten: first argument must be string")
			}

			separator := "."
			if len(args) > 1 {
				if sep, ok := args[1].(string); ok {
					separator = sep
				}
			}

			var obj map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.flatten: error parsing JSON object: %v", err))
			}

			flattened := flattenObject(obj, "", separator)

			result, err := json.Marshal(flattened)
			if err != nil {
				panic(fmt.Sprintf("JSON.flatten: error marshaling result: %v", err))
			}

			return string(result)
		}),

		"unflatten": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.unflatten needs (jsonObjectString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.unflatten: first argument must be string")
			}

			separator := "."
			if len(args) > 1 {
				if sep, ok := args[1].(string); ok {
					separator = sep
				}
			}

			var flattened map[string]interface{}
			err := json.Unmarshal([]byte(jsonStr), &flattened)
			if err != nil {
				panic(fmt.Sprintf("JSON.unflatten: error parsing JSON object: %v", err))
			}

			unflattened := unflattenObject(flattened, separator)

			result, err := json.Marshal(unflattened)
			if err != nil {
				panic(fmt.Sprintf("JSON.unflatten: error marshaling result: %v", err))
			}

			return string(result)
		}),

		"query": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("JSON.query needs (jsonString, path)")
			}
			jsonStr, ok1 := args[0].(string)
			path, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("JSON.query: arguments must be strings")
			}

			var obj interface{}
			err := json.Unmarshal([]byte(jsonStr), &obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.query: error parsing JSON: %v", err))
			}

			result := queryJSONPath(obj, path)
			return convertJSONToR2(result)
		}),

		"size": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.size needs (jsonString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.size: first argument must be string")
			}

			var obj interface{}
			err := json.Unmarshal([]byte(jsonStr), &obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.size: error parsing JSON: %v", err))
			}

			return float64(calculateJSONSize(obj))
		}),

		"type": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.type needs (jsonString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.type: first argument must be string")
			}

			var obj interface{}
			err := json.Unmarshal([]byte(jsonStr), &obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.type: error parsing JSON: %v", err))
			}

			return getJSONType(obj)
		}),

		"pretty": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.pretty needs (jsonString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.pretty: first argument must be string")
			}

			indent := "  "
			if len(args) > 1 {
				if ind, ok := args[1].(string); ok {
					indent = ind
				}
			}

			var obj interface{}
			err := json.Unmarshal([]byte(jsonStr), &obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.pretty: error parsing JSON: %v", err))
			}

			result, err := json.MarshalIndent(obj, "", indent)
			if err != nil {
				panic(fmt.Sprintf("JSON.pretty: error formatting JSON: %v", err))
			}

			return string(result)
		}),

		"minify": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("JSON.minify needs (jsonString)")
			}
			jsonStr, ok := args[0].(string)
			if !ok {
				panic("JSON.minify: first argument must be string")
			}

			var obj interface{}
			err := json.Unmarshal([]byte(jsonStr), &obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.minify: error parsing JSON: %v", err))
			}

			result, err := json.Marshal(obj)
			if err != nil {
				panic(fmt.Sprintf("JSON.minify: error minifying JSON: %v", err))
			}

			return string(result)
		}),
	}

	RegisterModule(env, "json", functions)
}

func convertJSONToR2(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = convertJSONToR2(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = convertJSONToR2(val)
		}
		return result
	case float64:
		return v
	case string:
		return v
	case bool:
		return v
	case nil:
		return nil
	default:
		return v
	}
}

func convertJSONArrayToR2(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, val := range arr {
		result[i] = convertJSONToR2(val)
	}
	return result
}

func convertJSONObjectToR2(obj map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, val := range obj {
		result[key] = convertJSONToR2(val)
	}
	return result
}

func convertR2ToJSON(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = convertR2ToJSON(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = convertR2ToJSON(val)
		}
		return result
	case float64:
		return v
	case string:
		return v
	case bool:
		return v
	case nil:
		return nil
	case *r2core.DateValue:
		return v.Time.Format("2006-01-02T15:04:05Z07:00")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func deepMergeObjects(obj1, obj2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, val := range obj1 {
		result[key] = val
	}

	for key, val := range obj2 {
		if existing, exists := result[key]; exists {
			if existingMap, ok := existing.(map[string]interface{}); ok {
				if valMap, ok := val.(map[string]interface{}); ok {
					result[key] = deepMergeObjects(existingMap, valMap)
					continue
				}
			}
		}
		result[key] = val
	}

	return result
}

func flattenObject(obj map[string]interface{}, prefix, separator string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, val := range obj {
		newKey := key
		if prefix != "" {
			newKey = prefix + separator + key
		}

		switch v := val.(type) {
		case map[string]interface{}:
			nested := flattenObject(v, newKey, separator)
			for nestedKey, nestedVal := range nested {
				result[nestedKey] = nestedVal
			}
		case []interface{}:
			for i, item := range v {
				itemKey := newKey + separator + strconv.Itoa(i)
				if itemMap, ok := item.(map[string]interface{}); ok {
					nested := flattenObject(itemMap, itemKey, separator)
					for nestedKey, nestedVal := range nested {
						result[nestedKey] = nestedVal
					}
				} else {
					result[itemKey] = item
				}
			}
		default:
			result[newKey] = val
		}
	}

	return result
}

func unflattenObject(flattened map[string]interface{}, separator string) map[string]interface{} {
	result := make(map[string]interface{})

	for key, val := range flattened {
		parts := strings.Split(key, separator)
		current := result

		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = val
			} else {
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				if nested, ok := current[part].(map[string]interface{}); ok {
					current = nested
				}
			}
		}
	}

	return result
}

func queryJSONPath(obj interface{}, path string) interface{} {
	if path == "" || path == "$" {
		return obj
	}

	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	current := obj

	for _, part := range parts {
		if part == "" {
			continue
		}

		if strings.Contains(part, "[") && strings.Contains(part, "]") {
			// Handle array access like "items[0]"
			keyPart := part[:strings.Index(part, "[")]
			indexPart := part[strings.Index(part, "[")+1 : strings.Index(part, "]")]

			if keyPart != "" {
				if objMap, ok := current.(map[string]interface{}); ok {
					current = objMap[keyPart]
				} else {
					return nil
				}
			}

			if arr, ok := current.([]interface{}); ok {
				if index, err := strconv.Atoi(indexPart); err == nil && index >= 0 && index < len(arr) {
					current = arr[index]
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			if objMap, ok := current.(map[string]interface{}); ok {
				current = objMap[part]
			} else {
				return nil
			}
		}
	}

	return current
}

func calculateJSONSize(obj interface{}) int {
	v := reflect.ValueOf(obj)
	switch v.Kind() {
	case reflect.Map:
		return v.Len()
	case reflect.Slice, reflect.Array:
		return v.Len()
	case reflect.String:
		return len(v.String())
	default:
		return 1
	}
}

func getJSONType(obj interface{}) string {
	switch obj.(type) {
	case map[string]interface{}:
		return "object"
	case []interface{}:
		return "array"
	case float64:
		return "number"
	case string:
		return "string"
	case bool:
		return "boolean"
	case nil:
		return "null"
	default:
		return "unknown"
	}
}
