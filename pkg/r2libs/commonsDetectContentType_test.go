package r2libs

import "testing"

func TestDetectContentType(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  ContentType
	}{
		{"empty", "", Unknown},
		{"whitespace only", "   ", Unknown},
		{"json object", `{"a":1}`, JSONType},
		{"json array", `[1,2,3]`, JSONType},
		{"plain text", "plain text", TextType},
		{"full html document", "<html><body>hi</body></html>", HTMLType},
		{"doctype html", "<!DOCTYPE html><html></html>", HTMLType},
		{"bare body tag", "<body>hi</body>", HTMLType},
		{
			// Regression: a well-formed XML document whose element happens
			// to be named "a" (a perfectly valid XML tag name, and also the
			// HTML anchor tag name) must not be misclassified as HTML just
			// because a loose substring pattern matches "<a". Previously
			// isHTML() ran before isXML() and matched on tag-name substrings
			// alone, so any well-formed XML containing a <div>/<span>/<p>/
			// <a>/<script>/<style> element anywhere was wrongly reported as
			// text/html.
			"well-formed xml with html-like tag name",
			"<order><a>123</a></order>",
			XMLType,
		},
		{"well-formed xml simple root", "<root><value>1</value></root>", XMLType},
		{"malformed unclosed tag falls back to html heuristic", "<a>", HTMLType},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := DetectContentType(c.input)
			if got != c.want {
				t.Errorf("DetectContentType(%q) = %v, want %v", c.input, got, c.want)
			}
		})
	}
}
