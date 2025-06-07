package v2

// UpdateItemInputTemplate generates utility functions for creating DynamoDB UpdateItemInput.
// This template creates helper functions for updating items using complete objects,
// raw values, conditional updates, and expression builders.
const UpdateItemInputTemplate = `
// UpdateItemInput creates UpdateItemInput from a SchemaItem (updates all non-key attributes).
//
// Parameters:
//   - item: SchemaItem with updated values
//
// Returns:
//   - *dynamodb.UpdateItemInput: Ready for UpdateItem operation
//   - error: If key extraction or marshaling fails
//
// Example:
//   item := SchemaItem{Id: "user123", Created: 1640995200, Name: "Updated Name", Age: 30}
//   updateInput, err := UpdateItemInput(item)
//   _, err = client.UpdateItem(ctx, updateInput)
func UpdateItemInput(item SchemaItem) (*dynamodb.UpdateItemInput, error) {
   key, err := KeyInput(item)
   if err != nil {
       return nil, fmt.Errorf("failed to create key from item for update: %v", err)
   }
   
   // Marshal the entire item to get all attributes
   allAttributes, err := attributevalue.MarshalMap(item)
   if err != nil {
       return nil, fmt.Errorf("failed to marshal item for update: %v", err)
   }
   
   // Remove key attributes from update
   updates := make(map[string]types.AttributeValue)
   for attrName, attrValue := range allAttributes {
       if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
           updates[attrName] = attrValue
       }
   }
   
   if len(updates) == 0 {
       return nil, fmt.Errorf("no non-key attributes to update")
   }
   
   var updateExpressionParts []string
   expressionAttributeNames := make(map[string]string)
   expressionAttributeValues := make(map[string]types.AttributeValue)
   
   i := 0
   for attrName, attrValue := range updates {
       nameKey := fmt.Sprintf("#attr%d", i)
       valueKey := fmt.Sprintf(":val%d", i)
       
       updateExpressionParts = append(updateExpressionParts, fmt.Sprintf("%s = %s", nameKey, valueKey))
       expressionAttributeNames[nameKey] = attrName
       expressionAttributeValues[valueKey] = attrValue
       i++
   }
   
   updateExpression := "SET " + strings.Join(updateExpressionParts, ", ")
   
   return &dynamodb.UpdateItemInput{
       TableName:                 aws.String(TableSchema.TableName),
       Key:                       key,
       UpdateExpression:          aws.String(updateExpression),
       ExpressionAttributeNames:  expressionAttributeNames,
       ExpressionAttributeValues: expressionAttributeValues,
   }, nil
}

// UpdateItemInputFromRaw creates UpdateItemInput from raw key values and update map.
//
// Parameters:
//   - hashKeyValue: Hash key value (any type)
//   - rangeKeyValue: Range key value (any type, nil if no range key)
//   - updates: Map of attribute names to new values
//
// Returns:
//   - *dynamodb.UpdateItemInput: Ready for UpdateItem operation
//   - error: If key creation or marshaling fails
//
// Example:
//   updates := map[string]interface{}{"name": "Updated Name", "age": 30}
//   updateInput, err := UpdateItemInputFromRaw("user123", 1640995200, updates)
//   _, err = client.UpdateItem(ctx, updateInput)
func UpdateItemInputFromRaw(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}) (*dynamodb.UpdateItemInput, error) {
   key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
   if err != nil {
       return nil, fmt.Errorf("failed to create key for update: %v", err)
   }
   
   if len(updates) == 0 {
       return nil, fmt.Errorf("no updates provided")
   }
   
   var updateExpressionParts []string
   expressionAttributeNames := make(map[string]string)
   expressionAttributeValues := make(map[string]types.AttributeValue)
   
   i := 0
   for attrName, value := range updates {
       nameKey := fmt.Sprintf("#attr%d", i)
       valueKey := fmt.Sprintf(":val%d", i)
       
       updateExpressionParts = append(updateExpressionParts, fmt.Sprintf("%s = %s", nameKey, valueKey))
       expressionAttributeNames[nameKey] = attrName
       
       av, err := attributevalue.Marshal(value)
       if err != nil {
           return nil, fmt.Errorf("failed to marshal update value for %s: %v", attrName, err)
       }
       expressionAttributeValues[valueKey] = av
       i++
   }
   
   updateExpression := "SET " + strings.Join(updateExpressionParts, ", ")
   
   return &dynamodb.UpdateItemInput{
       TableName:                 aws.String(TableSchema.TableName),
       Key:                       key,
       UpdateExpression:          aws.String(updateExpression),
       ExpressionAttributeNames:  expressionAttributeNames,
       ExpressionAttributeValues: expressionAttributeValues,
   }, nil
}

// UpdateItemInputWithCondition creates UpdateItemInput with conditional expression (optimistic locking).
//
// Parameters:
//   - hashKeyValue: Hash key value
//   - rangeKeyValue: Range key value (nil if no range key)
//   - updates: Map of attributes to update
//   - conditionExpression: DynamoDB condition (e.g., "#version = :v")
//   - conditionAttributeNames: Attribute name mappings for condition
//   - conditionAttributeValues: Attribute value mappings for condition
//
// Returns:
//   - *dynamodb.UpdateItemInput: Ready for conditional UpdateItem
//   - error: If input preparation fails
//
// Example:
//   updates := map[string]interface{}{"name": "New Name", "version": 2}
//   updateInput, err := UpdateItemInputWithCondition(
//       "user123", 1640995200, updates,
//       "#version = :currentVersion",
//       map[string]string{"#version": "version"},
//       map[string]types.AttributeValue{":currentVersion": &types.AttributeValueMemberN{Value: "1"}},
//   )
//   _, err = client.UpdateItem(ctx, updateInput)
func UpdateItemInputWithCondition(hashKeyValue interface{}, rangeKeyValue interface{}, updates map[string]interface{}, conditionExpression string, conditionAttributeNames map[string]string, conditionAttributeValues map[string]types.AttributeValue) (*dynamodb.UpdateItemInput, error) {
   updateInput, err := UpdateItemInputFromRaw(hashKeyValue, rangeKeyValue, updates)
   if err != nil {
       return nil, err
   }
   
   updateInput.ConditionExpression = aws.String(conditionExpression)
   
   // Merge condition attribute names with update attribute names
   if conditionAttributeNames != nil {
       for key, value := range conditionAttributeNames {
           updateInput.ExpressionAttributeNames[key] = value
       }
   }
   
   // Merge condition attribute values with update attribute values
   if conditionAttributeValues != nil {
       for key, value := range conditionAttributeValues {
           updateInput.ExpressionAttributeValues[key] = value
       }
   }
   
   return updateInput, nil
}

// UpdateItemInputWithExpression creates UpdateItemInput using DynamoDB expression builder (advanced).
//
// Parameters:
//   - hashKeyValue: Hash key value
//   - rangeKeyValue: Range key value (nil if no range key)
//   - updateBuilder: DynamoDB expression update builder
//   - conditionBuilder: Optional condition builder (nil if none)
//
// Returns:
//   - *dynamodb.UpdateItemInput: Ready for expression-based UpdateItem
//   - error: If expression building fails
//
// Example:
//   update := expression.Set(expression.Name("name"), expression.Value("New Name")).
//             Add(expression.Name("views"), expression.Value(1))
//   condition := expression.Equal(expression.Name("version"), expression.Value(1))
//   updateInput, err := UpdateItemInputWithExpression("user123", 1640995200, update, &condition)
//   _, err = client.UpdateItem(ctx, updateInput)
func UpdateItemInputWithExpression(hashKeyValue interface{}, rangeKeyValue interface{}, updateBuilder expression.UpdateBuilder, conditionBuilder *expression.ConditionBuilder) (*dynamodb.UpdateItemInput, error) {
   key, err := KeyInputFromRaw(hashKeyValue, rangeKeyValue)
   if err != nil {
       return nil, fmt.Errorf("failed to create key for expression update: %v", err)
   }
   
   var expr expression.Expression
   if conditionBuilder != nil {
       expr, err = expression.NewBuilder().
           WithUpdate(updateBuilder).
           WithCondition(*conditionBuilder).
           Build()
   } else {
       expr, err = expression.NewBuilder().
           WithUpdate(updateBuilder).
           Build()
   }
   
   if err != nil {
       return nil, fmt.Errorf("failed to build update expression: %v", err)
   }
   
   input := &dynamodb.UpdateItemInput{
       TableName:                 aws.String(TableSchema.TableName),
       Key:                       key,
       UpdateExpression:          expr.Update(),
       ExpressionAttributeNames:  expr.Names(),
       ExpressionAttributeValues: expr.Values(),
   }
   
   if conditionBuilder != nil {
       input.ConditionExpression = expr.Condition()
   }
   
   return input, nil
}
`
