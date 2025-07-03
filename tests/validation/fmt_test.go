package validation

// TestGeneratedCodeFormatting validates that complete DynamoDB code generation produces properly formatted Go code.
//
// Test process:
//  1. Reads JSON schema files from .tmpl/ directory (skipping files with 'invalid-' prefix)
//  2. Generates complete Go code using DynamoDB templates
//  3. Runs formatting validation (go fmt, goimports, gofumpt)
//
// This ensures generated code is production-ready and passes all standard Go formatting tools.
// func TestGeneratedCodeFormatting(t *testing.T) {
// 	schemaFiles, err := filepath.Glob(filepath.Join(EXAMPLES, "*.json"))
// 	require.NoError(t, err, "Failed to read template files")
// 	require.NotEmpty(t, schemaFiles, "No JSON files found in %s", EXAMPLES)

// 	for _, schemaFile := range schemaFiles {
// 		schemaFile := schemaFile
// 		schemaName := strings.TrimSuffix(filepath.Base(schemaFile), ".json")

// 		if strings.HasPrefix(schemaName, "invalid-") {
// 			t.Logf("Skipping invalid schema for compilation test: %s", schemaName)
// 			continue
// 		}

// 		t.Run(schemaName, func(t *testing.T) {
// 			t.Parallel()

// 			dynamoSchema, err := generator.Load(schemaFile)
// 			require.NoError(t, err, "Failed to load schema: %s", schemaFile)

// 			templateMap := v2.TemplateMap{
// 				PackageName:      dynamoSchema.PackageName(),
// 				TableName:        dynamoSchema.TableName(),
// 				HashKey:          dynamoSchema.HashKey(),
// 				RangeKey:         dynamoSchema.RangeKey(),
// 				Attributes:       dynamoSchema.Attributes(),
// 				CommonAttributes: dynamoSchema.CommonAttributes(),
// 				AllAttributes:    dynamoSchema.AllAttributes(),
// 				SecondaryIndexes: dynamoSchema.SecondaryIndexes(),
// 			}

// 			generatedCode := tmpl.MustParseTemplateFormattedToString(v2.CodeTemplate, templateMap)
// 			require.NotEmpty(t, generatedCode, "Generated code is empty")
// 			AllFormattersUnchanged(t, generatedCode)
// 		})
// 	}
// }
