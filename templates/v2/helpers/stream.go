package helpers

// StreamHelpersTemplate ...
const StreamHelpersTemplate = `
// ExtractFromDynamoDBStreamEvent extracts SchemaItem from DynamoDB stream event
// Converts stream AttributeValues to DynamoDB types and uses AWS SDK unmarshaler
func ExtractFromDynamoDBStreamEvent(dbEvent events.DynamoDBEventRecord) (*SchemaItem, error) {
    if dbEvent.Change.NewImage == nil {
        return nil, fmt.Errorf("new image is nil in the event")
    }
  
    // Convert stream AttributeValues to DynamoDB types
    dynamoAttrs := toDynamoMap(dbEvent.Change.NewImage)
    
    var item SchemaItem
    if err := attributevalue.UnmarshalMap(dynamoAttrs, &item); err != nil {
        return nil, fmt.Errorf("failed to unmarshal DynamoDB stream event: %v", err)
    }
  
    return &item, nil
}

// ExtractOldFromDynamoDBStreamEvent extracts old SchemaItem from DynamoDB stream event
// Converts stream AttributeValues to DynamoDB types and uses AWS SDK unmarshaler
func ExtractOldFromDynamoDBStreamEvent(dbEvent events.DynamoDBEventRecord) (*SchemaItem, error) {
    if dbEvent.Change.OldImage == nil {
        return nil, fmt.Errorf("old image is nil in the event")
    }
  
    // Convert stream AttributeValues to DynamoDB types
    dynamoAttrs := toDynamoMap(dbEvent.Change.OldImage)
    
    var item SchemaItem
    if err := attributevalue.UnmarshalMap(dynamoAttrs, &item); err != nil {
        return nil, fmt.Errorf("failed to unmarshal old DynamoDB stream event: %v", err)
    }
  
    return &item, nil
}

// toDynamoMap converts events.DynamoDBAttributeValue to types.AttributeValue
func toDynamoMap(streamAttrs map[string]events.DynamoDBAttributeValue) map[string]types.AttributeValue {
    dynamoAttrs := make(map[string]types.AttributeValue, len(streamAttrs))
    
    for key, streamAttr := range streamAttrs {
        dynamoAttrs[key] = toDynamoAttr(streamAttr)
    }
    
    return dynamoAttrs
}

// toDynamoAttr converts single events.DynamoDBAttributeValue to types.AttributeValue
func toDynamoAttr(streamAttr events.DynamoDBAttributeValue) types.AttributeValue {
    // Use DataType to properly identify the attribute type
    switch streamAttr.DataType() {
    case events.DataTypeString:
        return &types.AttributeValueMemberS{Value: streamAttr.String()}
    case events.DataTypeNumber:
        return &types.AttributeValueMemberN{Value: streamAttr.Number()}
    case events.DataTypeBoolean:
        return &types.AttributeValueMemberBOOL{Value: streamAttr.Boolean()}
    case events.DataTypeStringSet:
        return &types.AttributeValueMemberSS{Value: streamAttr.StringSet()}
    case events.DataTypeNumberSet:
        return &types.AttributeValueMemberNS{Value: streamAttr.NumberSet()}
    case events.DataTypeBinarySet:
        return &types.AttributeValueMemberBS{Value: streamAttr.BinarySet()}
    case events.DataTypeBinary:
        return &types.AttributeValueMemberB{Value: streamAttr.Binary()}
    case events.DataTypeList:
        list := make([]types.AttributeValue, len(streamAttr.List()))
        for i, item := range streamAttr.List() {
            list[i] = toDynamoAttr(item)
        }
        return &types.AttributeValueMemberL{Value: list}
    case events.DataTypeMap:
        m := make(map[string]types.AttributeValue, len(streamAttr.Map()))
        for k, v := range streamAttr.Map() {
            m[k] = toDynamoAttr(v)
        }
        return &types.AttributeValueMemberM{Value: m}
    case events.DataTypeNull:
        return &types.AttributeValueMemberNULL{Value: true}
    default:
        // Fallback for unknown types
        return &types.AttributeValueMemberNULL{Value: true}
    }
}

// IsFieldModified checks if a specific field was modified in a MODIFY event
func IsFieldModified(dbEvent events.DynamoDBEventRecord, fieldName string) bool {
    if dbEvent.EventName != "MODIFY" {
        return false
    }
    
    if dbEvent.Change.OldImage == nil || dbEvent.Change.NewImage == nil {
        return false
    }
    
    oldVal, oldExists := dbEvent.Change.OldImage[fieldName]
    newVal, newExists := dbEvent.Change.NewImage[fieldName]
    
    // Field was added
    if !oldExists && newExists {
        return true
    }
    
    // Field was removed
    if oldExists && !newExists {
        return true
    }
    
    // Field exists in both - check if values differ
    if oldExists && newExists {
        return !streamAttributeValuesEqual(oldVal, newVal)
    }
    
    return false
}

// streamAttributeValuesEqual compares two stream AttributeValues for equality
func streamAttributeValuesEqual(a, b events.DynamoDBAttributeValue) bool {
    // First check if data types are the same
    if a.DataType() != b.DataType() {
        return false
    }
    
    // Compare based on data type
    switch a.DataType() {
    case events.DataTypeString:
        return a.String() == b.String()
    case events.DataTypeNumber:
        return a.Number() == b.Number()
    case events.DataTypeBoolean:
        return a.Boolean() == b.Boolean()
    case events.DataTypeStringSet:
        aSet, bSet := a.StringSet(), b.StringSet()
        if len(aSet) != len(bSet) {
            return false
        }
        setMap := make(map[string]bool, len(aSet))
        for _, item := range aSet {
            setMap[item] = true
        }
        for _, item := range bSet {
            if !setMap[item] {
                return false
            }
        }
        return true
    case events.DataTypeNumberSet:
        aSet, bSet := a.NumberSet(), b.NumberSet()
        if len(aSet) != len(bSet) {
            return false
        }
        setMap := make(map[string]bool, len(aSet))
        for _, item := range aSet {
            setMap[item] = true
        }
        for _, item := range bSet {
            if !setMap[item] {
                return false
            }
        }
        return true
    case events.DataTypeNull:
        return true // Both are null
    default:
        // For complex types (List, Map, Binary), fall back to simple comparison
        // In real scenarios, you might want more sophisticated comparison
        return false
    }
}

// GetBoolFieldChanged checks if a boolean field changed from false to true
func GetBoolFieldChanged(dbEvent events.DynamoDBEventRecord, fieldName string) bool {
    if dbEvent.EventName != "MODIFY" {
        return false
    }
    
    if dbEvent.Change.OldImage == nil || dbEvent.Change.NewImage == nil {
        return false
    }
    
    oldValue := false
    if oldVal, ok := dbEvent.Change.OldImage[fieldName]; ok {
        oldValue = oldVal.Boolean()
    }
    
    newValue := false
    if newVal, ok := dbEvent.Change.NewImage[fieldName]; ok {
        newValue = newVal.Boolean()
    }
    
    return !oldValue && newValue
}

// ExtractBothFromDynamoDBStreamEvent extracts both old and new items from stream event
func ExtractBothFromDynamoDBStreamEvent(dbEvent events.DynamoDBEventRecord) (*SchemaItem, *SchemaItem, error) {
    var oldItem, newItem *SchemaItem
    var err error
    
    if dbEvent.Change.OldImage != nil {
        oldItem, err = ExtractOldFromDynamoDBStreamEvent(dbEvent)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to extract old item: %v", err)
        }
    }
    
    if dbEvent.Change.NewImage != nil {
        newItem, err = ExtractFromDynamoDBStreamEvent(dbEvent)
        if err != nil {
            return nil, nil, fmt.Errorf("failed to extract new item: %v", err)
        }
    }
    
    return oldItem, newItem, nil
}

// CreateTriggerHandler creates a handler function for DynamoDB stream events
func CreateTriggerHandler(
    onInsert func(context.Context, *SchemaItem) error,
    onModify func(context.Context, *SchemaItem, *SchemaItem) error,
    onDelete func(context.Context, map[string]events.DynamoDBAttributeValue) error,
) func(ctx context.Context, event events.DynamoDBEvent) error {
    return func(ctx context.Context, event events.DynamoDBEvent) error {
        for _, record := range event.Records {
            switch record.EventName {
            case "INSERT":
                if onInsert != nil {
                    item, err := ExtractFromDynamoDBStreamEvent(record)
                    if err != nil {
                        return err
                    }
                    if err := onInsert(ctx, item); err != nil {
                        return err
                    }
                }
                
            case "MODIFY":
                if onModify != nil {
                    oldItem, newItem, err := ExtractBothFromDynamoDBStreamEvent(record)
                    if err != nil {
                        return err
                    }
                    
                    if err := onModify(ctx, oldItem, newItem); err != nil {
                        return err
                    }
                }
                
            case "REMOVE":
                if onDelete != nil {
                    if err := onDelete(ctx, record.Change.OldImage); err != nil {
                        return err
                    }
                }
            }
        }
        return nil
    }
}
`
