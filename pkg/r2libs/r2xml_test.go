package r2libs

import (
	"testing"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

func xmlTestModule(t *testing.T) map[string]interface{} {
	t.Helper()
	env := r2core.NewEnvironment()
	RegisterXML(env)
	xmlModule, _ := env.Get("xml")
	return xmlModule.(map[string]interface{})
}

func TestXMLParseIgnoresNamespaceDeclarations(t *testing.T) {
	module := xmlTestModule(t)
	parseFunc := module["parse"].(r2core.BuiltinFunction)

	result := parseFunc(`<root xmlns="http://default" xmlns:ns1="http://a" id="5"><foo>x</foo></root>`)
	node, ok := result.(map[string]interface{})
	if !ok {
		t.Fatal("Expected object from XML parse")
	}

	attrs, ok := node["attributes"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected attributes map")
	}

	if _, exists := attrs["xmlns"]; exists {
		t.Error("Expected default xmlns declaration to not leak into attributes")
	}
	if _, exists := attrs["ns1"]; exists {
		t.Error("Expected xmlns:ns1 declaration to not leak into attributes")
	}
	if attrs["id"] != "5" {
		t.Errorf("Expected real attribute 'id' to be preserved, got: %v", attrs["id"])
	}
}

func TestXMLFromJSONRoundTripsRepeatedElements(t *testing.T) {
	module := xmlTestModule(t)
	parseFunc := module["parse"].(r2core.BuiltinFunction)
	toJSONFunc := module["toJSON"].(r2core.BuiltinFunction)
	fromJSONFunc := module["fromJSON"].(r2core.BuiltinFunction)
	stringifyFunc := module["stringify"].(r2core.BuiltinFunction)

	doc := parseFunc(`<root><item>1</item><item>2</item></root>`).(map[string]interface{})
	j := toJSONFunc(doc).(map[string]interface{})
	back := fromJSONFunc(j).(map[string]interface{})

	xmlStr := stringifyFunc(back).(string)
	if xmlStr != "<root><item>1</item><item>2</item></root>" {
		t.Errorf("Expected repeated <item> elements to round-trip through toJSON/fromJSON, got: %s", xmlStr)
	}
}

func TestXMLFromJSONKeepsScalarLeafChildren(t *testing.T) {
	module := xmlTestModule(t)
	fromJSONFunc := module["fromJSON"].(r2core.BuiltinFunction)
	stringifyFunc := module["stringify"].(r2core.BuiltinFunction)

	j := map[string]interface{}{
		"person": map[string]interface{}{
			"name": "John",
			"age":  "30",
		},
	}

	back := fromJSONFunc(j).(map[string]interface{})
	xmlStr := stringifyFunc(back).(string)

	if xmlStr != "<person><name>John</name><age>30</age></person>" &&
		xmlStr != "<person><age>30</age><name>John</name></person>" {
		t.Errorf("Expected plain-string JSON fields to become text-content child elements, got: %s", xmlStr)
	}
}
