package templates

var QueryBuilderFilterMethodsTemplate = `
{{range .AllAttributes}}
{{if eq (TypeGo .Type) "int"}}
func (qb *QueryBuilder) With{{SafeName .Name | ToCamelCase}}Between(start, end {{TypeGo .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.KeyConditions[attrName] = expression.Key(attrName).Between(expression.Value(start), expression.Value(end))
    qb.UsedKeys[attrName] = true
    return qb
}

func (qb *QueryBuilder) With{{SafeName .Name | ToCamelCase}}GreaterThan(value {{TypeGo .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.KeyConditions[attrName] = expression.Key(attrName).GreaterThan(expression.Value(value))
    qb.UsedKeys[attrName] = true
    return qb
}

func (qb *QueryBuilder) With{{SafeName .Name | ToCamelCase}}LessThan(value {{TypeGo .Type}}) *QueryBuilder {
    attrName := "{{.Name}}"
    qb.KeyConditions[attrName] = expression.Key(attrName).LessThan(expression.Value(value))
    qb.UsedKeys[attrName] = true
    return qb
}
{{end}}
{{end}}

func (qb *QueryBuilder) OrderByDesc() *QueryBuilder {
    qb.SortDescending = true
    return qb
}

func (qb *QueryBuilder) OrderByAsc() *QueryBuilder {
    qb.SortDescending = false
    return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.LimitValue = &limit
    return qb
}

func (qb *QueryBuilder) StartFrom(lastEvaluatedKey map[string]types.AttributeValue) *QueryBuilder {
    qb.ExclusiveStartKey = lastEvaluatedKey
    return qb
}
`
