package codegen

import (
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/jinzhu/inflection"
)

type Parameter struct {
	Name string
	Type string
	Tag  string
}

type RequestBody struct {
	Name   string
	Fields map[string]string
}

type Endpoint struct {
	Path        string
	Method      string
	OperationID string
	Params      []Parameter
	Body        *RequestBody
	State       int
	Response    *RequestBody // Add this field to handle response
}

type Component struct {
	Name   string
	Fields map[string]string
}

type OpenAPIStructs struct {
	Endpoints  map[string][]Endpoint
	Components []Component
}

var funcMap = template.FuncMap{
	"ToCamel": toCamel,
	"ToUpper": toUpper,
}

func GenerateCode(openapiPath string) {
	packagePath := "./handlr"
	packageName := "handlr"
	// Check if directory exists
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		// Create directory
		errDir := os.MkdirAll(packagePath, 0755)
		if errDir != nil {
			log.Fatalf("Error creating directory: %v", errDir)
		}
	}
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(openapiPath)
	if err != nil {
		log.Fatalf("Error loading OpenAPI document: %v", err)
	}

	err = doc.Validate(loader.Context)
	if err != nil {
		log.Fatalf("Error validating OpenAPI document: %v", err)
	}

	openAPIStructs := extractEndpointsAndComponents(doc)
	generateStructs(openAPIStructs, packageName)
}

func getTag(tags []string) string {
	if len(tags) > 0 {
		return tags[0]
	} else {
		return "default"
	}
}

func extractEndpointsAndComponents(doc *openapi3.T) OpenAPIStructs {
	endpoints := make(map[string][]Endpoint, 0)
	var components []Component
	var hasRequest bool

	for path, pathItem := range doc.Paths.Map() {
		for method, operation := range pathItem.Operations() {
			var params []Parameter
			for _, param := range operation.Parameters {
				params = append(params, Parameter{
					Name: param.Value.Name,
					Type: goType(param.Value.Schema),
					Tag:  param.Value.In,
				})
				hasRequest = true
			}

			var requestBody *RequestBody
			if operation.RequestBody != nil {
				for contentType, content := range operation.RequestBody.Value.Content {
					if contentType == "application/json" {
						schemaRef := content.Schema
						fields := make(map[string]string)
						for fieldName, fieldSchema := range schemaRef.Value.Properties {
							fields[fieldName] = goType(fieldSchema)
						}

						requestBody = &RequestBody{
							Name:   cutPrefix(schemaRef.Ref),
							Fields: fields,
						}
						hasRequest = true
					}
				}
			}

			var responseBody *RequestBody
			if len(operation.Responses.Map()) > 0 {
				for _, response := range operation.Responses.Map() {
					for contentType, content := range response.Value.Content {
						if contentType == "application/json" {
							schemaRef := content.Schema
							fields := make(map[string]string)
							for fieldName, fieldSchema := range schemaRef.Value.Properties {
								fields[fieldName] = goType(fieldSchema)
							}
							name := cutPrefix(schemaRef.Ref)
							if len(schemaRef.Value.Type.Slice()) > 0 && schemaRef.Value.Type.Slice()[0] == "array" {
								name = "[]" + name
							}
							responseBody = &RequestBody{
								Name:   name,
								Fields: fields,
							}
						}
					}
				}
			}

			// Split the path into segments
			segments := strings.Split(path, "/")

			// Iterate over the segments and capitalize the first letter of each one
			for i, segment := range segments {
				if len(segment) > 0 {
					segments[i] = strings.Title(segment)
				}
			}

			// Join the segments back together
			paths := strings.Join(segments, "")

			operationId := strings.Title(strings.ToLower(method)) + paths

			re := regexp.MustCompile(`\{(.*?)\}`)
			operationId = re.ReplaceAllStringFunc(operationId, func(s string) string {
				return strings.Trim(s, "{}")
			})
			// no request no response
			t := 0
			if hasRequest {
				t = 1
			}
			if responseBody != nil {
				t = 2
			}
			if hasRequest && responseBody != nil {
				t = 3
			}
			tag := getTag(operation.Tags)
			endpoints[tag] = append(endpoints[tag], Endpoint{
				Path:        path,
				Method:      method,
				OperationID: operationId,
				Params:      params,
				Body:        requestBody,
				Response:    responseBody,
				State:       t,
			})
		}
	}

	componentMap := processComponents(doc)

	for _, v := range componentMap {
		components = append(components, v)
	}

	return OpenAPIStructs{
		Endpoints:  endpoints,
		Components: components,
	}
}

func processComponents(doc *openapi3.T) map[string]Component {
	componentMap := make(map[string]Component)
	for componentName, componentSchema := range doc.Components.Schemas {
		fields := make(map[string]string)
		singularName := toCamel(componentName)

		pluralName := inflection.Plural(singularName)

		t := componentSchema.Value.Type.Slice()

		if len(t) > 0 && t[0] == "array" {
			// Handle array schema
			itemsSchema := componentSchema.Value.Items.Value
			for fieldName, fieldSchema := range itemsSchema.Properties {
				fields[fieldName] = goType(fieldSchema)
			}
			componentMap[singularName] = Component{
				Name:   singularName,
				Fields: fields,
			}

			componentMap[pluralName] = Component{
				Name: pluralName,
				Fields: map[string]string{
					singularName: "[]" + singularName,
				},
			}
		} else {
			// Handle object schema
			for fieldName, fieldSchema := range componentSchema.Value.Properties {
				fields[fieldName] = goType(fieldSchema)
			}
			componentMap[singularName] = Component{
				Name:   singularName,
				Fields: fields,
			}
		}
	}
	return componentMap
}
