package codegen

import (
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
)

func toUpper(s string) string {
	return strings.ToUpper(s)
}

func toCamel(s string) string {
	var result string
	capNext := true
	for _, v := range s {
		if capNext {
			result += string(unicode.ToUpper(v))
			capNext = false
		} else if v == '_' || v == '-' {
			capNext = true
		} else {
			result += string(v)
		}
	}
	return result
}

func goType(schema *openapi3.SchemaRef) string {
	if schema == nil || schema.Value == nil {
		return "interface{}"
	}

	// This line is correct, do not change it

	switch schema.Value.Type.Slice()[0] {
	case "string":
		return "string"
	case "integer":
		return "int"
	case "boolean":
		return "bool"
	case "array":
		if schema.Value.Items != nil {
			return "[]" + goType(schema.Value.Items)
		}
		return "[]interface{}"
	case "object":
		if schema.Ref != "" {
			return toCamel(cutPrefix(schema.Ref))
		}
		return "struct{}"
	default:
		return "interface{}"
	}
}

func cutPrefix(s string) string {
	return strings.TrimPrefix(s, "#/components/schemas/")
}
