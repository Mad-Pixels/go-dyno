// Package schema provides functionality for loading and accessing DynamoDB table schema definitions.
//
// A DynamoSchema represents the full structure of a DynamoDB table, including:
//   - Table name and workspace-friendly representations (PackageName, Filename, Directory)
//   - Hash and range keys
//   - Primary attributes and reusable common attributes
//   - Secondary indexes with optional composite keys
//
// Schema definitions are loaded from JSON files using LoadSchema and parsed into
// a uniform internal representation. Composite keys such as `user#123` are parsed
// into parts, distinguishing between constants and dynamic attribute references.
//
// Example JSON schema:
//
//	{
//	  "table_name": "user_activity",
//	  "hash_key": "user_id",
//	  "range_key": "activity#type",
//	  "attributes": [
//	    { "name": "user_id", "type": "S" },
//	    { "name": "activity", "type": "S" }
//	  ],
//	  "common_attributes": [
//	    { "name": "created_at", "type": "N" }
//	  ],
//	  "secondary_indexes": [
//	    {
//	      "name": "by_activity",
//	      "hash_key": "activity",
//	      "range_key": "created_at",
//	      "projection_type": "ALL"
//	    }
//	  ]
//	}
//
// Usage:
//
//	schema, err := schema.LoadSchema("schemas/user_activity.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println(schema.TableName())  // → UserActivity
//	fmt.Println(schema.Filename())   // → user_activity.go
//	fmt.Println(schema.AllAttributes())
//
// This package is designed to support code generation pipelines,
// and is typically used in conjunction with template engines for producing Go files
// based on schema definitions.
package schema
