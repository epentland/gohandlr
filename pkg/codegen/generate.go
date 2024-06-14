package codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	views "github.com/epentland/gohandlr/pkg/codegen/templates"
)

// findAndAddComment finds a function with a specific name and adds a comment above it, overwriting any existing comments.
func findAndAddComment(f *ast.File, funcName, comment string) bool {
	found := false
	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == funcName {
			found = true
			// Create a new comment group
			cg := &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Slash: fn.Pos() - 1, // Place the comment just before the function
						Text:  comment,
					},
				},
			}

			// Loop through the existing comments to see if there's already a comment group
			if f.Comments != nil {
				for i, c := range f.Comments {
					if c.End() == fn.Pos()-1 {
						// Overwrite the existing comment group
						f.Comments[i] = cg
						return false // stop the traversal
					}
				}
			} else {
				// Add the new comment group
				f.Comments = append(f.Comments, cg)
			}
			return false // stop the traversal
		}
		return true
	})
	return found
}

func addHandlerFunction(content string, endpoint Endpoint, tmpl *template.Template) string {
	functionName := "process" + endpoint.OperationID
	fmt.Println(functionName)
	if !strings.Contains(content, "func "+functionName) {
		buf := &bytes.Buffer{}
		fmt.Println("adding function:", functionName)
		err := tmpl.ExecuteTemplate(buf, "ProcessEndpoint", endpoint)
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}

		content += "\n" + buf.String()
		fmt.Println("adding new content", buf.String())
	}
	return content
}

// addHandlerToRegisterHandlers reads the RegisterHandlers function as a string and appends the handler call if it doesn't exist
func addHandlerToRegisterHandlers(content, handlerName string) string {
	registerFuncStart := "func RegisterHandlers(r *chi.Mux) {"
	handlerCall := fmt.Sprintf("	r.MethodFunc(%s())\n", handlerName)

	// Check if the handlerName exists in the content without considering the parameters
	if !strings.Contains(content, handlerName) {
		content = strings.Replace(content, registerFuncStart, registerFuncStart+"\n"+handlerCall, 1)
	}

	return content
}

func generateStructs(openAPIStructs OpenAPIStructs, packageName string) {
	tmpl := parseTemplates()

	generateFile(tmpl, "structs", "./"+packageName+"/structs.go", openAPIStructs)
	generateFile(tmpl, "handlers", "./"+packageName+"/handlers.go", openAPIStructs)

	// If process.go does not exist
	if _, err := os.Stat("./" + packageName + "/process.go"); os.IsNotExist(err) {
		generateFile(tmpl, "process", "./"+packageName+"/process.go", openAPIStructs)
		return
	}

	// Read the existing file content
	existingContent, err := os.ReadFile("./" + packageName + "/process.go")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	contentStr := string(existingContent)

	// Buffer to hold the new content
	var newContent bytes.Buffer
	newContent.WriteString(contentStr)

	for _, endpoints := range openAPIStructs.Endpoints {
		for _, endpoint := range endpoints {

			contentStr = addHandlerFunction(contentStr, endpoint, tmpl)

			// Add the handler to RegisterHandlers if it doesn't exist
			handlerName := "Handle" + endpoint.OperationID
			contentStr = addHandlerToRegisterHandlers(contentStr, handlerName)
		}
	}

	// Format the new content
	formattedContent, err := format.Source([]byte(contentStr))
	if err != nil {
		fmt.Println("Error formatting content:", err)
		return
	}

	// Write the formatted content to the file
	err = os.WriteFile("./"+packageName+"/process.go", formattedContent, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func parseTemplates() *template.Template {
	templatesFs := views.FS
	tmpl, err := template.New("struct").Funcs(funcMap).ParseFS(templatesFs, "*.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	return tmpl
}

func generateFile(tmpl *template.Template, templateName, fileName string, data interface{}) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating file %s: %v", fileName, err)
	}
	defer file.Close()

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		log.Fatalf("Error executing template %s: %v", templateName, err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Error formatting file %s: %v", fileName, err)
	}

	_, err = file.Write(formatted)
	if err != nil {
		log.Fatalf("Error writing formatted file %s: %v", fileName, err)
	}
}
