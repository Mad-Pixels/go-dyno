package templates

var MainTemplate = `package main

import (
	"bytes"
	"fmt"
	"text/template"
	
	"github.com/your-org/templates"
)

// Функция для генерации всего кода
func GenerateCode(data interface{}) (string, error) {
	// Создаем буфер для результирующего кода
	var result bytes.Buffer
	
	// Создаем список всех шаблонов
	templateParts := []string{
		templates.HeaderTemplate,
		templates.ConstantsTemplate,
		templates.SchemaStructsTemplate,
		templates.QueryBuilderBaseTemplate,
		templates.QueryBuilderIndexMethodsTemplate,
		templates.QueryBuilderBuildMethodsTemplate,
		templates.QueryBuilderExecuteMethodsTemplate,
		templates.QueryBuilderFilterMethodsTemplate,
		templates.ItemOperationsTemplate,
		templates.StreamOperationsTemplate,
		templates.KeyUtilsTemplate,
		templates.TriggerHandlersTemplate,
		templates.AttributeUtilsTemplate,
	}
	
	// Объединяем все шаблоны
	combinedTemplate := strings.Join(templateParts, "\n")
	
	// Создаем и компилируем шаблон
	tmpl, err := template.New("dynamodb-code").Funcs(template.FuncMap{
		"ToCamelCase":      ToCamelCase,
		"ToLowerCamelCase": ToLowerCamelCase,
		"SafeName":         SafeName,
		"TypeGo":           TypeGo,
		"TypeGoAttr":       TypeGoAttr,
	}).Parse(combinedTemplate)
	
	if err != nil {
		return "", fmt.Errorf("ошибка компиляции шаблона: %v", err)
	}
	
	// Выполняем шаблон
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения шаблона: %v", err)
	}
	
	return result.String(), nil
}
`

// Функции-помощники для шаблонов
var HelperFunctionsTemplate = `package helpers

import (
	"regexp"
	"strings"
)

// ToCamelCase преобразует строку в CamelCase
func ToCamelCase(s string) string {
	// Разделяем строку на части по не-буквенным и не-цифровым символам
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	parts := re.Split(s, -1)
	
	var result string
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Преобразуем первую букву в верхний регистр
		result += strings.ToUpper(part[:1]) + part[1:]
	}
	
	return result
}

// ToLowerCamelCase преобразует строку в lowerCamelCase
func ToLowerCamelCase(s string) string {
	camelCase := ToCamelCase(s)
	if len(camelCase) == 0 {
		return ""
	}
	return strings.ToLower(camelCase[:1]) + camelCase[1:]
}

// SafeName делает имя безопасным для использования в качестве идентификатора
func SafeName(s string) string {
	// Заменяем специальные символы на подчеркивание
	re := regexp.MustCompile("[^a-zA-Z0-9]+")
	return re.ReplaceAllString(s, "_")
}

// TypeGo возвращает Go-тип для DynamoDB типа
func TypeGo(dbType string) string {
	switch dbType {
	case "S":
		return "string"
	case "N":
		return "int"
	case "B":
		return "bool"
	default:
		return "interface{}"
	}
}

// TypeGoAttr возвращает Go-тип для атрибута по его имени
func TypeGoAttr(attrName string, attributes []Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return TypeGo(attr.Type)
		}
	}
	return "interface{}"
}
`
