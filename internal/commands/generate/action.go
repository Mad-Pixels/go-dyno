package generate

import (
	"path/filepath"

	"github.com/Mad-Pixels/go-dyno/internal/schema"
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
	dynamoSchema, err := schema.LoadSchema(cfgFl)
	if err != nil {
		return err
	}
	if err = utils.IsDirOrCreate(
		filepath.Join(destFl, dynamoSchema.Directory()),
	); err != nil {
		return err
	}

	return nil
}

// func processSchemaFile(jsonPath, rootDir string) {
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
