package tmpl

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustParseTemplate_BasicTemplate(t *testing.T) {
	var b bytes.Buffer
	tmpl := "Hello, {{ .Name }}!"
	vars := map[string]string{"Name": "World"}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "Hello, World!", b.String())
}

func TestMustParseTemplate_WithBuiltinFunctions(t *testing.T) {
	var b bytes.Buffer
	tmpl := "{{ ToUpperCamelCase .Field }}: {{ ToSafeName .Key }}"
	vars := map[string]string{"Field": "user_name", "Key": "1invalid"}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "UserName: x1invalid", b.String())
}

func TestMustParseTemplate_JoinFunction(t *testing.T) {
	var b bytes.Buffer
	tmpl := `{{ Join .Items ", " }}`
	vars := map[string][]string{"Items": {"apple", "banana", "cherry"}}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "apple, banana, cherry", b.String())
}

func TestMustParseTemplateFormatted_SimpleStruct(t *testing.T) {
	var b bytes.Buffer
	tmpl := "package main\n\ntype {{ .Name }} struct {\n{{ .Field }} {{ .Type }}\n}"
	vars := map[string]string{"Name": "User", "Field": "Name", "Type": "string"}

	MustParseTemplateFormatted(&b, tmpl, vars)

	result := b.String()
	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "type User struct")
	assert.Contains(t, result, "Name string")
}

func TestMustParseTemplateToString_ReturnString(t *testing.T) {
	tmpl := "Value: {{ .Value }}"
	vars := map[string]int{"Value": 42}

	result := MustParseTemplateToString(tmpl, vars)
	assert.Equal(t, "Value: 42", result)
}

func TestMustParseTemplateFormattedToString_ReturnString(t *testing.T) {
	tmpl := "package main\nfunc test(){\nreturn\n}"

	result := MustParseTemplateFormattedToString(tmpl, nil)

	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "func test() {")
	assert.Contains(t, result, "return")
	assert.Contains(t, result, "}")
}

func TestMustParseTemplate_SliceFunction(t *testing.T) {
	var b bytes.Buffer
	tmpl := `{{ Slice .Text 2 }}`
	vars := map[string]string{"Text": "[]int"}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "int", b.String())
}

func TestMustParseTemplate_ToUpperFunction(t *testing.T) {
	var b bytes.Buffer
	tmpl := `{{ ToUpper .Text }}`
	vars := map[string]string{"Text": "hello"}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "HELLO", b.String())
}

func TestMustParseTemplate_EmptyTemplate(t *testing.T) {
	var b bytes.Buffer
	tmpl := ""
	vars := map[string]string{}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "", b.String())
}

func TestMustParseTemplate_NoVariables(t *testing.T) {
	var b bytes.Buffer
	tmpl := "Static text"

	MustParseTemplate(&b, tmpl, nil)
	assert.Equal(t, "Static text", b.String())
}

func TestMustParseTemplate_MultipleBuiltinFunctions(t *testing.T) {
	var b bytes.Buffer
	tmpl := `{{ ToUpperCamelCase .Field1 }} and {{ ToLowerCamelCase .Field2 }}`
	vars := map[string]string{"Field1": "first_name", "Field2": "last_name"}

	MustParseTemplate(&b, tmpl, vars)
	assert.Equal(t, "FirstName and lastName", b.String())
}

func TestMustParseTemplateFormatted_ImportsFormatting(t *testing.T) {
	var b bytes.Buffer
	tmpl := `package main
import "fmt"
import "os"
func main() {
fmt.Println("test")
os.Exit(0)
}`

	MustParseTemplateFormatted(&b, tmpl, nil)

	result := b.String()
	assert.Contains(t, result, "import (")
	assert.Contains(t, result, `"fmt"`)
	assert.Contains(t, result, `"os"`)
	assert.Contains(t, result, "fmt.Println")
	assert.Contains(t, result, "os.Exit")
}

func TestMustParseTemplateFormatted_StructFormatting(t *testing.T) {
	var b bytes.Buffer
	tmpl := `package main
type User struct{
Name string
Age int
Email string
}`

	MustParseTemplateFormatted(&b, tmpl, nil)

	result := b.String()
	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "type User struct {")
	assert.Contains(t, result, "Name")
	assert.Contains(t, result, "Age")
	assert.Contains(t, result, "Email")
	assert.Contains(t, result, "\t")
}

func TestMustParseTemplateFormatted_FunctionFormatting(t *testing.T) {
	var b bytes.Buffer
	tmpl := `package main
import "strconv"
func getValue(x int)string{
if x>0{
return strconv.Itoa(x)
}
return "zero"
}`

	MustParseTemplateFormatted(&b, tmpl, nil)

	result := b.String()
	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "func getValue(x int) string {")
	assert.Contains(t, result, "if x > 0 {")
	assert.Contains(t, result, "return strconv.Itoa(x)")
	assert.Contains(t, result, `return "zero"`)
	assert.Contains(t, result, "strconv")
}
