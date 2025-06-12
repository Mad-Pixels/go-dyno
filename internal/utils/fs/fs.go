// Package fs provides file system utilities for the go-dyno project.
//
// It includes:
//   - File operations: reading, writing, validation
//   - Directory operations: creation, removal, validation
//   - Path manipulation: extension handling, validation
//   - JSON operations: parsing files directly into structs
//
// All functions provide structured logging with context information
// and consistent error handling patterns.
//
// Example usage:
//
//	// Read and parse JSON
//	var config Config
//	err := fs.ReadAndParseJSON("config.json", &config)
//
//	// Create directories and files
//	err = fs.IsDirOrCreate("output/generated")
//	err = fs.WriteToFile("output/generated/model.go", codeBytes)
//
//	// Path manipulation
//	goFile := fs.AddFileExt("model", ".go") // "model.go"
package fs
