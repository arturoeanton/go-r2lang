package r2libs

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// XMLNode represents an XML node
type XMLNode struct {
	Name       string            `xml:"-"`
	Content    string            `xml:",chardata"`
	Attributes map[string]string `xml:"-"`
	Children   []*XMLNode        `xml:"-"`
}

func RegisterXML(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"parse": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.parse needs (xmlString)")
			}
			xmlStr, ok := args[0].(string)
			if !ok {
				panic("XML.parse: first argument must be string")
			}

			decoder := xml.NewDecoder(strings.NewReader(xmlStr))
			var root *XMLNode
			var stack []*XMLNode

			for {
				token, err := decoder.Token()
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(fmt.Sprintf("XML.parse: error parsing XML: %v", err))
				}

				switch element := token.(type) {
				case xml.StartElement:
					node := &XMLNode{
						Name:       element.Name.Local,
						Attributes: make(map[string]string),
						Children:   []*XMLNode{},
					}

					for _, attr := range element.Attr {
						// xmlns / xmlns:prefix declarations are namespace
						// plumbing, not user-facing content attributes; keeping
						// them out of Attributes avoids leaking them into
						// scripts that iterate/read attributes.
						if attr.Name.Space == "xmlns" || (attr.Name.Space == "" && attr.Name.Local == "xmlns") {
							continue
						}
						node.Attributes[attr.Name.Local] = attr.Value
					}

					if len(stack) == 0 {
						if root != nil {
							panic("XML.parse: multiple root elements found; XML documents must have exactly one root element")
						}
						root = node
					} else {
						parent := stack[len(stack)-1]
						parent.Children = append(parent.Children, node)
					}
					stack = append(stack, node)

				case xml.EndElement:
					if len(stack) > 0 {
						stack = stack[:len(stack)-1]
					}

				case xml.CharData:
					if len(stack) > 0 {
						content := strings.TrimSpace(string(element))
						if content != "" {
							stack[len(stack)-1].Content = content
						}
					}
				}
			}

			return convertXMLNodeToR2(root)
		}),

		"stringify": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.stringify needs (xmlObject)")
			}

			obj, ok := args[0].(map[string]interface{})
			if !ok {
				panic("XML.stringify: first argument must be object")
			}

			pretty := false
			if len(args) > 1 {
				if p, ok := args[1].(bool); ok {
					pretty = p
				}
			}

			xmlStr := convertR2ToXMLString(obj, "", pretty)
			return xmlStr
		}),

		"validate": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.validate needs (xmlString)")
			}
			xmlStr, ok := args[0].(string)
			if !ok {
				panic("XML.validate: first argument must be string")
			}

			decoder := xml.NewDecoder(strings.NewReader(xmlStr))
			for {
				_, err := decoder.Token()
				if err == io.EOF {
					return true
				}
				if err != nil {
					return false
				}
			}
		}),

		"getAttribute": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.getAttribute needs (xmlObject, attributeName)")
			}
			obj, ok1 := args[0].(map[string]interface{})
			attrName, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("XML.getAttribute: arguments must be (object, string)")
			}

			if attrs, exists := obj["attributes"]; exists {
				if attrMap, ok := attrs.(map[string]interface{}); ok {
					if value, found := attrMap[attrName]; found {
						return value
					}
				}
			}
			return nil
		}),

		"setAttribute": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 3 {
				panic("XML.setAttribute needs (xmlObject, attributeName, value)")
			}
			obj, ok1 := args[0].(map[string]interface{})
			attrName, ok2 := args[1].(string)
			value, ok3 := args[2].(string)
			if !ok1 || !ok2 || !ok3 {
				panic("XML.setAttribute: arguments must be (object, string, string)")
			}

			if obj["attributes"] == nil {
				obj["attributes"] = make(map[string]interface{})
			}

			if attrMap, ok := obj["attributes"].(map[string]interface{}); ok {
				attrMap[attrName] = value
			}

			return obj
		}),

		"getChildren": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.getChildren needs (xmlObject)")
			}
			obj, ok := args[0].(map[string]interface{})
			if !ok {
				panic("XML.getChildren: first argument must be object")
			}

			if children, exists := obj["children"]; exists {
				return children
			}
			return []interface{}{}
		}),

		"getChildByName": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.getChildByName needs (xmlObject, tagName)")
			}
			obj, ok1 := args[0].(map[string]interface{})
			tagName, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("XML.getChildByName: arguments must be (object, string)")
			}

			if children, exists := obj["children"]; exists {
				if childArray, ok := children.([]interface{}); ok {
					for _, child := range childArray {
						if childObj, ok := child.(map[string]interface{}); ok {
							if name, nameExists := childObj["name"]; nameExists {
								if name == tagName {
									return childObj
								}
							}
						}
					}
				}
			}
			return nil
		}),

		"getChildrenByName": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.getChildrenByName needs (xmlObject, tagName)")
			}
			obj, ok1 := args[0].(map[string]interface{})
			tagName, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("XML.getChildrenByName: arguments must be (object, string)")
			}

			var result []interface{}
			if children, exists := obj["children"]; exists {
				if childArray, ok := children.([]interface{}); ok {
					for _, child := range childArray {
						if childObj, ok := child.(map[string]interface{}); ok {
							if name, nameExists := childObj["name"]; nameExists {
								if name == tagName {
									result = append(result, childObj)
								}
							}
						}
					}
				}
			}
			return result
		}),

		"addChild": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.addChild needs (xmlObject, childObject)")
			}
			parent, ok1 := args[0].(map[string]interface{})
			child, ok2 := args[1].(map[string]interface{})
			if !ok1 || !ok2 {
				panic("XML.addChild: arguments must be objects")
			}

			if parent["children"] == nil {
				parent["children"] = []interface{}{}
			}

			if childArray, ok := parent["children"].([]interface{}); ok {
				parent["children"] = append(childArray, child)
			}

			return parent
		}),

		"removeChild": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.removeChild needs (xmlObject, childIndex)")
			}
			parent, ok1 := args[0].(map[string]interface{})
			index := int(toFloat(args[1]))
			if !ok1 {
				panic("XML.removeChild: first argument must be object")
			}

			if children, exists := parent["children"]; exists {
				if childArray, ok := children.([]interface{}); ok {
					if index >= 0 && index < len(childArray) {
						newChildren := append(childArray[:index], childArray[index+1:]...)
						parent["children"] = newChildren
					}
				}
			}

			return parent
		}),

		"createNode": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.createNode needs (tagName)")
			}
			tagName, ok := args[0].(string)
			if !ok {
				panic("XML.createNode: tagName must be string")
			}

			content := ""
			if len(args) > 1 {
				if c, ok := args[1].(string); ok {
					content = c
				}
			}

			node := map[string]interface{}{
				"name":       tagName,
				"content":    content,
				"attributes": make(map[string]interface{}),
				"children":   []interface{}{},
			}

			return node
		}),

		"findByPath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.findByPath needs (xmlObject, path)")
			}
			root, ok1 := args[0].(map[string]interface{})
			path, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("XML.findByPath: arguments must be (object, string)")
			}

			return findXMLByPath(root, path)
		}),

		"xpath": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("XML.xpath needs (xmlObject, xpathExpression)")
			}
			root, ok1 := args[0].(map[string]interface{})
			xpath, ok2 := args[1].(string)
			if !ok1 || !ok2 {
				panic("XML.xpath: arguments must be (object, string)")
			}

			// Simplified XPath implementation
			return simpleXPath(root, xpath)
		}),

		"toJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.toJSON needs (xmlObject)")
			}
			obj, ok := args[0].(map[string]interface{})
			if !ok {
				panic("XML.toJSON: first argument must be object")
			}

			return convertXMLToJSON(obj)
		}),

		"fromJSON": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.fromJSON needs (jsonObject)")
			}
			obj, ok := args[0].(map[string]interface{})
			if !ok {
				panic("XML.fromJSON: first argument must be object")
			}

			return convertJSONToXML(obj)
		}),

		"minify": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.minify needs (xmlString)")
			}
			xmlStr, ok := args[0].(string)
			if !ok {
				panic("XML.minify: first argument must be string")
			}

			// Remove unnecessary whitespace
			lines := strings.Split(xmlStr, "\n")
			var result []string
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					result = append(result, trimmed)
				}
			}
			return strings.Join(result, "")
		}),

		"pretty": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 1 {
				panic("XML.pretty needs (xmlString)")
			}
			xmlStr, ok := args[0].(string)
			if !ok {
				panic("XML.pretty: first argument must be string")
			}

			indent := "  "
			if len(args) > 1 {
				if ind, ok := args[1].(string); ok {
					indent = ind
				}
			}

			return formatXMLPretty(xmlStr, indent)
		}),
	}

	RegisterModule(env, "xml", functions)
}

// Helper functions
func convertXMLNodeToR2(node *XMLNode) map[string]interface{} {
	if node == nil {
		return nil
	}

	result := map[string]interface{}{
		"name":       node.Name,
		"content":    node.Content,
		"attributes": make(map[string]interface{}),
		"children":   []interface{}{},
	}

	// Convert attributes
	for key, value := range node.Attributes {
		result["attributes"].(map[string]interface{})[key] = value
	}

	// Convert children
	children := make([]interface{}, len(node.Children))
	for i, child := range node.Children {
		children[i] = convertXMLNodeToR2(child)
	}
	result["children"] = children

	return result
}

func escapeXMLText(s string) string {
	var buf strings.Builder
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		panic(fmt.Sprintf("XML.stringify: failed to escape text: %v", err))
	}
	return buf.String()
}

func escapeXMLAttr(s string) string {
	return escapeXMLText(s)
}

func convertR2ToXMLString(obj map[string]interface{}, indent string, pretty bool) string {
	return convertR2ToXMLStringSeen(obj, indent, pretty, make(map[uintptr]bool))
}

// seen tracks the active ancestor chain (by underlying map pointer) to reject
// self-referential nodes (e.g. via XML.addChild(node, node)), which would
// otherwise recurse forever and crash the process with a stack overflow.
func convertR2ToXMLStringSeen(obj map[string]interface{}, indent string, pretty bool, seen map[uintptr]bool) string {
	name, hasName := obj["name"].(string)
	if !hasName {
		return ""
	}

	ptr := reflect.ValueOf(obj).Pointer()
	if seen[ptr] {
		panic("XML.stringify: circular reference detected")
	}
	seen[ptr] = true
	defer delete(seen, ptr)

	var result strings.Builder

	if pretty {
		result.WriteString(indent)
	}

	result.WriteString("<")
	result.WriteString(name)

	// Add attributes
	if attrs, hasAttrs := obj["attributes"].(map[string]interface{}); hasAttrs {
		for key, value := range attrs {
			result.WriteString(fmt.Sprintf(" %s=\"%s\"", key, escapeXMLAttr(fmt.Sprintf("%v", value))))
		}
	}

	content, hasContent := obj["content"].(string)
	children, hasChildren := obj["children"].([]interface{})

	if !hasContent && (!hasChildren || len(children) == 0) {
		result.WriteString("/>")
		if pretty {
			result.WriteString("\n")
		}
		return result.String()
	}

	result.WriteString(">")

	if hasContent && content != "" {
		result.WriteString(escapeXMLText(content))
	}

	if hasChildren && len(children) > 0 {
		if pretty {
			result.WriteString("\n")
		}

		nextIndent := indent
		if pretty {
			nextIndent = indent + "  "
		}

		for _, child := range children {
			if childObj, ok := child.(map[string]interface{}); ok {
				result.WriteString(convertR2ToXMLStringSeen(childObj, nextIndent, pretty, seen))
			}
		}

		if pretty {
			result.WriteString(indent)
		}
	}

	result.WriteString("</")
	result.WriteString(name)
	result.WriteString(">")
	if pretty {
		result.WriteString("\n")
	}

	return result.String()
}

func findXMLByPath(root map[string]interface{}, path string) interface{} {
	if path == "" || path == "/" {
		return root
	}

	parts := strings.Split(strings.Trim(path, "/"), "/")
	current := root

	for _, part := range parts {
		if part == "" {
			continue
		}

		found := false
		if children, hasChildren := current["children"].([]interface{}); hasChildren {
			for _, child := range children {
				if childObj, ok := child.(map[string]interface{}); ok {
					if name, hasName := childObj["name"].(string); hasName && name == part {
						current = childObj
						found = true
						break
					}
				}
			}
		}

		if !found {
			return nil
		}
	}

	return current
}

func simpleXPath(root map[string]interface{}, xpath string) []interface{} {
	var results []interface{}

	// Very simplified XPath - just supports //tagname
	if strings.HasPrefix(xpath, "//") {
		tagName := strings.TrimPrefix(xpath, "//")
		results = findAllByTagName(root, tagName)
	} else if strings.HasPrefix(xpath, "/") {
		// Absolute path
		path := strings.TrimPrefix(xpath, "/")
		if result := findXMLByPath(root, path); result != nil {
			results = append(results, result)
		}
	}

	return results
}

func findAllByTagName(node map[string]interface{}, tagName string) []interface{} {
	return findAllByTagNameSeen(node, tagName, make(map[uintptr]bool))
}

func findAllByTagNameSeen(node map[string]interface{}, tagName string, seen map[uintptr]bool) []interface{} {
	ptr := reflect.ValueOf(node).Pointer()
	if seen[ptr] {
		panic("XML.xpath: circular reference detected")
	}
	seen[ptr] = true
	defer delete(seen, ptr)

	var results []interface{}

	if name, hasName := node["name"].(string); hasName && name == tagName {
		results = append(results, node)
	}

	if children, hasChildren := node["children"].([]interface{}); hasChildren {
		for _, child := range children {
			if childObj, ok := child.(map[string]interface{}); ok {
				results = append(results, findAllByTagNameSeen(childObj, tagName, seen)...)
			}
		}
	}

	return results
}

func convertXMLToJSON(xmlObj map[string]interface{}) map[string]interface{} {
	return convertXMLToJSONSeen(xmlObj, make(map[uintptr]bool))
}

func convertXMLToJSONSeen(xmlObj map[string]interface{}, seen map[uintptr]bool) map[string]interface{} {
	ptr := reflect.ValueOf(xmlObj).Pointer()
	if seen[ptr] {
		panic("XML.toJSON: circular reference detected")
	}
	seen[ptr] = true
	defer delete(seen, ptr)

	result := make(map[string]interface{})

	if name, hasName := xmlObj["name"].(string); hasName {
		nodeData := make(map[string]interface{})

		if content, hasContent := xmlObj["content"].(string); hasContent && content != "" {
			nodeData["_text"] = content
		}

		if attrs, hasAttrs := xmlObj["attributes"].(map[string]interface{}); hasAttrs {
			for key, value := range attrs {
				nodeData["@"+key] = value
			}
		}

		if children, hasChildren := xmlObj["children"].([]interface{}); hasChildren {
			for _, child := range children {
				if childObj, ok := child.(map[string]interface{}); ok {
					childJSON := convertXMLToJSONSeen(childObj, seen)
					for key, value := range childJSON {
						if existing, exists := nodeData[key]; exists {
							if existingList, isList := existing.([]interface{}); isList {
								nodeData[key] = append(existingList, value)
							} else {
								nodeData[key] = []interface{}{existing, value}
							}
							continue
						}
						nodeData[key] = value
					}
				}
			}
		}

		result[name] = nodeData
	}

	return result
}

func convertJSONToXML(jsonObj map[string]interface{}) map[string]interface{} {
	return convertJSONToXMLSeen(jsonObj, make(map[uintptr]bool))
}

func convertJSONToXMLSeen(jsonObj map[string]interface{}, seen map[uintptr]bool) map[string]interface{} {
	// This is a simplified conversion - real implementation would be more complex
	for key, value := range jsonObj {
		if valueMap, ok := value.(map[string]interface{}); ok {
			return buildXMLNodeFromJSONObject(key, valueMap, seen)
		}
	}

	return nil
}

func buildXMLNodeFromJSONObject(name string, valueMap map[string]interface{}, seen map[uintptr]bool) map[string]interface{} {
	ptr := reflect.ValueOf(valueMap).Pointer()
	if seen[ptr] {
		panic("XML.fromJSON: circular reference detected")
	}
	seen[ptr] = true
	defer delete(seen, ptr)

	result := map[string]interface{}{
		"name":       name,
		"content":    "",
		"attributes": make(map[string]interface{}),
		"children":   []interface{}{},
	}

	for subKey, subValue := range valueMap {
		if strings.HasPrefix(subKey, "@") {
			// Attribute
			attrName := strings.TrimPrefix(subKey, "@")
			result["attributes"].(map[string]interface{})[attrName] = subValue
		} else if subKey == "_text" {
			// Text content
			if str, ok := subValue.(string); ok {
				result["content"] = str
			}
		} else {
			// Child element(s): may be a single object, a plain scalar
			// (text-only element), or an array (repeated sibling tags, the
			// shape XML.toJSON itself produces for repeated elements).
			children := result["children"].([]interface{})
			children = append(children, jsonValueToXMLChildren(subKey, subValue, seen)...)
			result["children"] = children
		}
	}

	return result
}

func jsonValueToXMLChildren(name string, value interface{}, seen map[uintptr]bool) []interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return []interface{}{buildXMLNodeFromJSONObject(name, v, seen)}
	case []interface{}:
		children := make([]interface{}, 0, len(v))
		for _, item := range v {
			children = append(children, jsonValueToXMLChildren(name, item, seen)...)
		}
		return children
	case nil:
		return []interface{}{map[string]interface{}{
			"name":       name,
			"content":    "",
			"attributes": make(map[string]interface{}),
			"children":   []interface{}{},
		}}
	default:
		return []interface{}{map[string]interface{}{
			"name":       name,
			"content":    fmt.Sprintf("%v", v),
			"attributes": make(map[string]interface{}),
			"children":   []interface{}{},
		}}
	}
}

func formatXMLPretty(xmlStr string, indent string) string {
	// Simple pretty formatter - would need more sophisticated implementation for production
	var result strings.Builder
	level := 0
	inTag := false

	for i, char := range xmlStr {
		switch char {
		case '<':
			if !inTag {
				if i > 0 && xmlStr[i-1] != '>' {
					result.WriteString("\n")
					result.WriteString(strings.Repeat(indent, level))
				}
				inTag = true
			}
			result.WriteRune(char)
		case '>':
			result.WriteRune(char)
			if inTag {
				inTag = false
				if i+1 < len(xmlStr) && xmlStr[i+1] == '<' {
					if xmlStr[i-1] != '/' {
						level++
					}
				}
			}
		case '/':
			if inTag && i+1 < len(xmlStr) && xmlStr[i+1] == '>' {
				// Self-closing tag
			} else if inTag && i > 0 && xmlStr[i-1] == '<' {
				// Closing tag
				level--
				result.WriteString("\n")
				result.WriteString(strings.Repeat(indent, level))
			}
			result.WriteRune(char)
		default:
			result.WriteRune(char)
		}
	}

	return result.String()
}
