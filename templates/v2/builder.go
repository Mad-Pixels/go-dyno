package templates

var QueryBuilderBaseTemplate = `
type QueryBuilder struct {
    IndexName           string
    KeyConditions       map[string]expression.KeyConditionBuilder
    FilterConditions    []expression.ConditionBuilder
    UsedKeys            map[string]bool
    Attributes          map[string]interface{}
    SortDescending      bool
    LimitValue          *int
    ExclusiveStartKey   map[string]types.AttributeValue
    PreferredSortKey    string
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        KeyConditions:   make(map[string]expression.KeyConditionBuilder),
        UsedKeys:        make(map[string]bool),
        Attributes:      make(map[string]interface{}),
    }
}

{{range .AllAttributes}}
func (qb *QueryBuilder) With{{SafeName .Name | ToCamelCase}}({{SafeName .Name | ToLowerCamelCase}} {{TypeGo .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.Attributes[attrName] = {{SafeName .Name | ToLowerCamelCase}}
    qb.UsedKeys[attrName] = true
    return qb
}
{{end}}
`
