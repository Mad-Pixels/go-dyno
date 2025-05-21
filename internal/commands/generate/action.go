package generate

import (
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils"

	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) (err error) {
	var (
		cfgFl  = getFlagCfgValue(ctx)
		destFl = getFlagDestValue(ctx)
	)
	if err = utils.IsFileOrError(cfgFl); err != nil {
		return err
	}
	if err = utils.IsDirOrCreate(destFl); err != nil {
		return err
	}

	return nil
}

// type DynamoSchema struct {
// 	TableName        string           `json:"table_name"`
// 	HashKey          string           `json:"hash_key"`
// 	RangeKey         string           `json:"range_key"`
// 	Attributes       []Attribute      `json:"attributes"`
// 	CommonAttributes []Attribute      `json:"common_attributes"`
// 	SecondaryIndexes []SecondaryIndex `json:"secondary_indexes"`
// }

// type Attribute struct {
// 	Name string `json:"name"`
// 	Type string `json:"type"`
// }

// type CompositeKeyPart struct {
// 	IsConstant bool
// 	Value      string
// }

// type SecondaryIndex struct {
// 	Name             string `json:"name"`
// 	HashKey          string `json:"hash_key"`
// 	HashKeyParts     []CompositeKeyPart
// 	RangeKey         string `json:"range_key"`
// 	RangeKeyParts    []CompositeKeyPart
// 	ProjectionType   string   `json:"projection_type"`
// 	NonKeyAttributes []string `json:"non_key_attributes"`
// }

// const (
// 	jsonTag = "`json:\"S\"`"
// )

// func action(ctx *cli.Context) (err error) {
// 	var (
// 		cfgFl  = getFlagCfgValue(ctx)
// 		destFl = getFlagDestValue(ctx)
// 	)
//
// 	if err = utils.IsFileOrError(cfgFl); err != nil {
// 		return err
// 	}
// 	if err = utils.IsDirOrCreate(destFl); err != nil {
// 		return err
// 	}
//
// 	// processSchemaFile(cfgFl, destFl)
// 	return nil
// }

// func processSchemaFile(jsonPath, rootDir string) {
// 	jsonFile, err := os.ReadFile(jsonPath)
// 	if err != nil {
// 		fmt.Printf("Failed to read json %s: %v\n", jsonPath, err)
// 		return
// 	}
//
// 	var schema DynamoSchema
// 	err = json.Unmarshal(jsonFile, &schema)
// 	if err != nil {
// 		fmt.Printf("Failed to parse json %s: %v\n", jsonPath, err)
// 		return
// 	}
//
// 	packageName := strings.ReplaceAll(schema.TableName, "-", "")
// 	packageDir := filepath.Join(rootDir, "gen", packageName)
//
// 	if err = os.MkdirAll(packageDir, os.ModePerm); err != nil {
// 		fmt.Printf("Failed to create directory %s: %v\n", packageDir, err)
// 		return
// 	}
// 	outputPath := filepath.Join(packageDir, packageName+".go")
//
// 	funcMap := template.FuncMap{
// 		"ToCamelCase":      toCamelCase,
// 		"ToLowerCamelCase": toLowerCamelCase,
// 		"SafeName":         safeName,
// 		"TypeGo":           typeGo,
// 		"TypeZero":         typeZero,
// 		"TypeGoAttr":       typeGoAttr,
// 	}
// 	allAttributes := append(schema.Attributes, schema.CommonAttributes...)
//
// 	for i, idx := range schema.SecondaryIndexes {
// 		schema.SecondaryIndexes[i].HashKeyParts = parseCompositeKey(idx.HashKey, allAttributes)
// 		schema.SecondaryIndexes[i].RangeKeyParts = parseCompositeKey(idx.RangeKey, allAttributes)
//
// 		fmt.Printf("Index: %s, HashKeyParts: %+v, RangeKeyParts: %+v\n", idx.Name, schema.SecondaryIndexes[i].HashKeyParts, schema.SecondaryIndexes[i].RangeKeyParts)
// 	}
//
// 	schemaMap := map[string]interface{}{
// 		"PackageName":      packageName,
// 		"TableName":        schema.TableName,
// 		"HashKey":          schema.HashKey,
// 		"RangeKey":         schema.RangeKey,
// 		"Attributes":       schema.Attributes,
// 		"CommonAttributes": schema.CommonAttributes,
// 		"AllAttributes":    allAttributes,
// 		"SecondaryIndexes": schema.SecondaryIndexes,
// 	}
//
// 	tmpl, err := template.New("schema").Funcs(funcMap).Parse(tmpl.CodeTemplate)
// 	if err != nil {
// 		fmt.Printf("Failed to parse template: %v\n", err)
// 		return
// 	}
//
// 	outputFile, err := os.Create(outputPath)
// 	if err != nil {
// 		fmt.Printf("Failed to create output file %s: %v\n", outputPath, err)
// 		return
// 	}
// 	defer outputFile.Close()
//
// 	err = tmpl.Execute(outputFile, schemaMap)
// 	if err != nil {
// 		fmt.Printf("Failed to execute template for %s: %v\n", outputPath, err)
// 		return
// 	}
// 	fmt.Printf("Successfully generated %s!\n", schema.TableName)
// }
//
// func parseCompositeKey(key string, allAttributes []Attribute) []CompositeKeyPart {
// 	if key == "" {
// 		return nil
// 	}
// 	parts := strings.Split(key, "#")
// 	var result []CompositeKeyPart
// 	for _, part := range parts {
// 		if isAttribute(part, allAttributes) {
// 			result = append(result, CompositeKeyPart{IsConstant: false, Value: part})
// 		} else {
// 			result = append(result, CompositeKeyPart{IsConstant: true, Value: part})
// 		}
// 	}
// 	return result
// }
//
// func isAttribute(name string, attributes []Attribute) bool {
// 	for _, attr := range attributes {
// 		if attr.Name == name {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// func toCamelCase(s string) string {
// 	var result string
// 	capitalizeNext := true
// 	for _, r := range s {
// 		if r == '_' || r == '-' {
// 			capitalizeNext = true
// 		} else if capitalizeNext {
// 			result += string(unicode.ToUpper(r))
// 			capitalizeNext = false
// 		} else {
// 			result += string(r)
// 		}
// 	}
// 	return result
// }
//
// func toLowerCamelCase(s string) string {
// 	if s == "" {
// 		return ""
// 	}
// 	s = toCamelCase(s)
// 	return strings.ToLower(s[:1]) + s[1:]
// }
//
// var reservedWords = map[string]bool{
// 	// List of Go reserved words
// 	"break":       true,
// 	"default":     true,
// 	"func":        true,
// 	"interface":   true,
// 	"select":      true,
// 	"case":        true,
// 	"defer":       true,
// 	"go":          true,
// 	"map":         true,
// 	"struct":      true,
// 	"chan":        true,
// 	"else":        true,
// 	"goto":        true,
// 	"package":     true,
// 	"switch":      true,
// 	"const":       true,
// 	"fallthrough": true,
// 	"if":          true,
// 	"range":       true,
// 	"type":        true,
// 	"continue":    true,
// 	"for":         true,
// 	"import":      true,
// 	"return":      true,
// 	"var":         true,
// }
//
// func safeName(s string) string {
// 	// Сначала заменяем # на _
// 	s = strings.ReplaceAll(s, "#", "_")
//
// 	if reservedWords[s] {
// 		return s + "_"
// 	}
// 	return s
// }
//
// func typeGo(dynamoType string) string {
// 	switch dynamoType {
// 	case "S":
// 		return "string"
// 	case "N":
// 		return "int"
// 	case "B":
// 		return "bool"
// 	default:
// 		return "interface{}"
// 	}
// }
//
// func typeZero(dynamoType string) string {
// 	switch dynamoType {
// 	case "S":
// 		return `""`
// 	case "N":
// 		return "0"
// 	case "B":
// 		return "false"
// 	default:
// 		return "nil"
// 	}
// }
//
// func typeGoAttr(attrName string, attributes []Attribute) string {
// 	for _, attr := range attributes {
// 		if attr.Name == attrName {
// 			return typeGo(attr.Type)
// 		}
// 	}
// 	return "interface{}"
// }
