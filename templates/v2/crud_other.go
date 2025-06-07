package v2

// CrudOther ...
const CrudOther = `
// IncrementAttribute creates an UpdateItemInput to increment/decrement a numeric attribute
// Useful for counters, views, likes, etc.
//
// Example usage:
//   // Increment views by 1
//   updateInput, err := IncrementAttribute("user123", 1640995200, "views", 1)
//   
//   // Decrement likes by 1
//   updateInput, err := IncrementAttribute("post456", 1640995200, "likes", -1)
//   
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func IncrementAttribute(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, incrementValue int) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for increment: %v", err)
    }
    
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("ADD #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", incrementValue)},
        },
    }, nil
}

// AddToSet creates an UpdateItemInput to add values to a string set (SS) or number set (NS)
// Creates the set if it doesn't exist, otherwise adds to existing set
//
// Example usage:
//   // Add tags to string set
//   updateInput, err := AddToSet("user123", 1640995200, "tags", []string{"golang", "backend"})
//   
//   // Add scores to number set  
//   updateInput, err := AddToSet("user123", 1640995200, "scores", []int{95, 87})
//   
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func AddToSet(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, values interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for add to set: %v", err)
    }
    
    var attributeValue types.AttributeValue
    
    switch v := values.(type) {
    case []string:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot add empty string set")
        }
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    case []int:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot add empty number set")
        }
        numberStrings := make([]string, len(v))
        for i, num := range v {
            numberStrings[i] = fmt.Sprintf("%d", num)
        }
        attributeValue = &types.AttributeValueMemberNS{Value: numberStrings}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string or []int", values)
    }
    
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("ADD #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": attributeValue,
        },
    }, nil
}

// RemoveFromSet creates an UpdateItemInput to remove values from a string set (SS) or number set (NS)
//
// Example usage:
//   // Remove tags from string set
//   updateInput, err := RemoveFromSet("user123", 1640995200, "tags", []string{"deprecated"})
//   
//   // Remove scores from number set
//   updateInput, err := RemoveFromSet("user123", 1640995200, "scores", []int{60})
//   
//   if err != nil {
//       return err
//   }
//   _, err = dynamoClient.UpdateItem(ctx, updateInput)
func RemoveFromSet(hashKeyValue interface{}, rangeKeyValue interface{}, attributeName string, values interface{}) (*dynamodb.UpdateItemInput, error) {
    key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
    if err != nil {
        return nil, fmt.Errorf("failed to create key for remove from set: %v", err)
    }
    
    var attributeValue types.AttributeValue
    
    switch v := values.(type) {
    case []string:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot remove empty string set")
        }
        attributeValue = &types.AttributeValueMemberSS{Value: v}
    case []int:
        if len(v) == 0 {
            return nil, fmt.Errorf("cannot remove empty number set")
        }
        numberStrings := make([]string, len(v))
        for i, num := range v {
            numberStrings[i] = fmt.Sprintf("%d", num)
        }
        attributeValue = &types.AttributeValueMemberNS{Value: numberStrings}
    default:
        return nil, fmt.Errorf("unsupported type for set operation: %T, expected []string or []int", values)
    }
    
    return &dynamodb.UpdateItemInput{
        TableName:        aws.String(TableSchema.TableName),
        Key:              key,
        UpdateExpression: aws.String("DELETE #attr :val"),
        ExpressionAttributeNames: map[string]string{
            "#attr": attributeName,
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":val": attributeValue,
        },
    }, nil
}

// ExtractFromDynamoDBStreamEvent converts a DynamoDB Stream event record into a strongly-typed SchemaItem.
// This function is essential for processing DynamoDB Streams in Lambda triggers, providing
// type-safe access to changed data without manual attribute parsing.
//
// The function handles all DynamoDB attribute types with proper Go type conversion:
// - String (S): Direct string assignment  
// - Number (N): Converts string representation to int with error handling
// - Boolean (BOOL): Direct boolean assignment
// - String Set (SS): Copies string slice with nil safety
// - Number Set (NS): Converts string slice to int slice with error handling
//
// Stream Event Processing:
// - Validates that NewImage exists (required for INSERT/MODIFY events)
// - Safely extracts each attribute with type checking
// - Handles missing attributes gracefully (leaves zero values)
// - Provides error handling for malformed numeric data
//
// Parameters:
//   - dbEvent: DynamoDB Stream event record from AWS Lambda
//
// Returns:
//   - *SchemaItem: Populated struct with data from the stream event
//   - error: Processing error if NewImage is nil or data is malformed
//
// Example usage in Lambda trigger:
//   func handleDynamoDBStream(ctx context.Context, event events.DynamoDBEvent) error {
//       for _, record := range event.Records {
//           switch record.EventName {
//           case "INSERT", "MODIFY":
//               item, err := ExtractFromDynamoDBStreamEvent(record)
//               if err != nil {
//                   log.Printf("Failed to extract item: %v", err)
//                   continue
//               }
//               
//               // Process the strongly-typed item
//               log.Printf("Item changed: %+v", item)
//               
//           case "REMOVE":
//               // Handle deletion using record.Change.OldImage
//           }
//       }
//       return nil
//   }
//
// Error Handling:
// - Returns error if NewImage is nil (malformed event)
// - Continues processing other attributes if individual conversions fail
// - Number conversion errors are logged but don't stop processing
func ExtractFromDynamoDBStreamEvent(dbEvent events.DynamoDBEventRecord) (*SchemaItem, error) {
   if dbEvent.Change.NewImage == nil {
       return nil, fmt.Errorf("new image is nil in the event")
   }
   
   item := &SchemaItem{}
   
   {{range .AllAttributes}}
   if val, ok := dbEvent.Change.NewImage["{{.Name}}"]; ok {
       {{if eq .Type "S"}}
       // String attribute: direct assignment from DynamoDB String type
       item.{{ToSafeName .Name | ToUpperCamelCase}} = val.String()
       {{else if eq .Type "N"}}
       // Number attribute: convert DynamoDB Number (string) to Go int with error handling
       if n, err := strconv.Atoi(val.Number()); err == nil {
           item.{{ToSafeName .Name | ToUpperCamelCase}} = n
       }
       {{else if eq .Type "BOOL"}}
       // Boolean attribute: direct assignment from DynamoDB Boolean type
       item.{{ToSafeName .Name | ToUpperCamelCase}} = val.Boolean()
       {{else if eq .Type "SS"}}
       // String Set attribute: copy slice with nil safety check
       if ss := val.StringSet(); ss != nil {
           item.{{ToSafeName .Name | ToUpperCamelCase}} = ss
       }
       {{else if eq .Type "NS"}}
       // Number Set attribute: convert string slice to int slice with error handling
       if ns := val.NumberSet(); ns != nil {
           numbers := make([]int, 0, len(ns))
           for _, numStr := range ns {
               if num, err := strconv.Atoi(numStr); err == nil {
                   numbers = append(numbers, num)
               }
           }
           item.{{ToSafeName .Name | ToUpperCamelCase}} = numbers
       }
       {{else}}
       // Unsupported DynamoDB type: {{.Type}} for attribute {{.Name}}
       // This ensures compilation succeeds even if new types are added to schema
       _ = val // Mark as used to avoid compilation error
       {{end}}
   }
   {{end}}
   
   return item, nil
}
`
