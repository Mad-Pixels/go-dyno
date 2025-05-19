package templates

var QueryBuilderBuildMethodsTemplate = `
func (qb *QueryBuilder) Build() (string, expression.KeyConditionBuilder, *expression.ConditionBuilder, map[string]types.AttributeValue, error) {
    var filterCond *expression.ConditionBuilder

    // Определяем индексы по приоритету: сначала проверяем индексы с большим количеством атрибутов ключа
    sortedIndexes := make([]SecondaryIndex, len(TableSchema.SecondaryIndexes))
    copy(sortedIndexes, TableSchema.SecondaryIndexes)
    
    // Сортируем индексы по убыванию количества частей в ключе (индексы с большим числом ключевых атрибутов - приоритетнее)
    // Сортируем индексы, отдавая предпочтение тем, у которых RangeKey совпадает с предпочтительным ключом сортировки
    sort.Slice(sortedIndexes, func(i, j int) bool {
        // Если задан предпочтительный ключ сортировки
        if qb.PreferredSortKey != "" {
            // Проверяем совпадение RangeKey с предпочтительным ключом
            iMatches := sortedIndexes[i].RangeKey == qb.PreferredSortKey
            jMatches := sortedIndexes[j].RangeKey == qb.PreferredSortKey
            
            // Если только один из индексов совпадает, предпочитаем его
            if iMatches && !jMatches {
                return true
            }
            if !iMatches && jMatches {
                return false
            }
        }
        
        // В остальных случаях сортируем по количеству атрибутов в ключе
        iParts := 0
        if sortedIndexes[i].HashKeyParts != nil {
            iParts += len(sortedIndexes[i].HashKeyParts)
        }
        if sortedIndexes[i].RangeKeyParts != nil {
            iParts += len(sortedIndexes[i].RangeKeyParts)
        }
        
        jParts := 0
        if sortedIndexes[j].HashKeyParts != nil {
            jParts += len(sortedIndexes[j].HashKeyParts)
        }
        if sortedIndexes[j].RangeKeyParts != nil {
            jParts += len(sortedIndexes[j].RangeKeyParts)
        }
        
        return iParts > jParts
    })

    // Пытаемся найти индекс, подходящий для запроса
    for _, idx := range sortedIndexes {
        var hashKeyCondition, rangeKeyCondition *expression.KeyConditionBuilder
        var hashKeyMatch, rangeKeyMatch bool

        // Проверка HashKey
        if idx.HashKeyParts != nil {
            if qb.hasAllKeys(idx.HashKeyParts) {
                cond := qb.buildCompositeKeyCondition(idx.HashKeyParts)
                hashKeyCondition = &cond
                hashKeyMatch = true
            }
        } else if idx.HashKey != "" && qb.UsedKeys[idx.HashKey] {
            cond := expression.Key(idx.HashKey).Equal(expression.Value(qb.Attributes[idx.HashKey]))
            hashKeyCondition = &cond
            hashKeyMatch = true
        }

        if !hashKeyMatch {
            continue // Этот индекс не подходит
        }

        // Проверка RangeKey с особым приоритетом для сортировки
        // Если у нас есть несколько подходящих индексов с одинаковым ключом хеша,
        // мы выберем тот, который лучше соответствует условиям сортировки
        if idx.RangeKeyParts != nil {
            if qb.hasAllKeys(idx.RangeKeyParts) {
                cond := qb.buildCompositeKeyCondition(idx.RangeKeyParts)
                rangeKeyCondition = &cond
                rangeKeyMatch = true
            }
        } else if idx.RangeKey != "" {
            // Проверяем, есть ли конкретное условие по ключу сортировки
            if qb.UsedKeys[idx.RangeKey] {
                if cond, exists := qb.KeyConditions[idx.RangeKey]; exists {
                    rangeKeyCondition = &cond
                    rangeKeyMatch = true
                } else {
                    // Индекс все равно считается подходящим, даже если нет условия по ключу сортировки
                    rangeKeyMatch = true
                }
            } else {
                // Индекс все равно считается подходящим, даже если нет условия по ключу сортировки
                rangeKeyMatch = true
            }
        } else {
            rangeKeyMatch = true
        }

        if !rangeKeyMatch {
            continue
        }

        // Нашли подходящий индекс!
        keyCondition := *hashKeyCondition
        if rangeKeyCondition != nil {
            keyCondition = keyCondition.And(*rangeKeyCondition)
        }

        // Собираем фильтры для всех атрибутов, которые не являются частью ключа
        for attrName, value := range qb.Attributes {
            // Проверяем, не является ли атрибут частью ключа индекса
            isPartOfHashKey := false
            isPartOfRangeKey := false
            
            if idx.HashKeyParts != nil {
                for _, part := range idx.HashKeyParts {
                    if !part.IsConstant && part.Value == attrName {
                        isPartOfHashKey = true
                        break
                    }
                }
            } else if attrName == idx.HashKey {
                isPartOfHashKey = true
            }
            
            if idx.RangeKeyParts != nil {
                for _, part := range idx.RangeKeyParts {
                    if !part.IsConstant && part.Value == attrName {
                        isPartOfRangeKey = true
                        break
                    }
                }
            } else if attrName == idx.RangeKey {
                isPartOfRangeKey = true
            }
            
            // Если атрибут не является частью ключа, добавляем его в фильтр
            if !isPartOfHashKey && !isPartOfRangeKey {
                cond := expression.Name(attrName).Equal(expression.Value(value))
                qb.FilterConditions = append(qb.FilterConditions, cond)
            }
        }

        // Объединяем все условия фильтрации
        if len(qb.FilterConditions) > 0 {
            combinedFilter := qb.FilterConditions[0]
            for _, cond := range qb.FilterConditions[1:] {
                combinedFilter = combinedFilter.And(cond)
            }
            filterCond = &combinedFilter
        }

        return idx.Name, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }

    // Если ни один вторичный индекс не подходит, пробуем использовать основной ключ таблицы
    if qb.UsedKeys[TableSchema.HashKey] {
        indexName := ""
        keyCondition := expression.Key(TableSchema.HashKey).Equal(expression.Value(qb.Attributes[TableSchema.HashKey]))

        // Добавляем условие по ключу диапазона, если он есть
        if TableSchema.RangeKey != "" && qb.UsedKeys[TableSchema.RangeKey] {
            if cond, exists := qb.KeyConditions[TableSchema.RangeKey]; exists {
                keyCondition = keyCondition.And(cond)
            } else {
                keyCondition = keyCondition.And(expression.Key(TableSchema.RangeKey).Equal(expression.Value(qb.Attributes[TableSchema.RangeKey])))
            }
        }

        // Добавляем все остальные атрибуты как фильтры
        for attrName, value := range qb.Attributes {
            if attrName != TableSchema.HashKey && attrName != TableSchema.RangeKey {
                cond := expression.Name(attrName).Equal(expression.Value(value))
                qb.FilterConditions = append(qb.FilterConditions, cond)
            }
        }

        // Объединяем все условия фильтрации
        if len(qb.FilterConditions) > 0 {
            combinedFilter := qb.FilterConditions[0]
            for _, cond := range qb.FilterConditions[1:] {
                combinedFilter = combinedFilter.And(cond)
            }
            filterCond = &combinedFilter
        }

        return indexName, keyCondition, filterCond, qb.ExclusiveStartKey, nil
    }

    return "", expression.KeyConditionBuilder{}, nil, nil, fmt.Errorf("no suitable index found for the provided keys")
}

func (qb *QueryBuilder) hasAllKeys(parts []CompositeKeyPart) bool {
    for _, part := range parts {
        if !part.IsConstant && !qb.UsedKeys[part.Value] {
            return false
        }
    }
    return true
}

func (qb *QueryBuilder) buildCompositeKeyCondition(parts []CompositeKeyPart) expression.KeyConditionBuilder {
    var compositeKeyValue string
    for i, part := range parts {
        var valueStr string
        if part.IsConstant {
            valueStr = part.Value
        } else {
            value := qb.Attributes[part.Value]
            valueStr = fmt.Sprintf("%v", value)
        }
        if i > 0 {
            compositeKeyValue += "#"
        }
        compositeKeyValue += valueStr
    }
    compositeKeyName := qb.getCompositeKeyName(parts)
    return expression.Key(compositeKeyName).Equal(expression.Value(compositeKeyValue))
}

func (qb *QueryBuilder) getCompositeKeyName(parts []CompositeKeyPart) string {
    var names []string
    for _, part := range parts {
        names = append(names, part.Value)
    }
    return strings.Join(names, "#")
}
`
