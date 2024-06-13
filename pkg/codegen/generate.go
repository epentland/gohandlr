package codegen

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"text/template"

	views "github.com/epentland/gohandlr/pkg/codegen/templates"
)

func generateStructs(openAPIStructs OpenAPIStructs, packageName string) {
	templatesFs := views.FS
	tmpl, err := template.New("struct").Funcs(funcMap).ParseFS(templatesFs, "*.tmpl")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	// Generate struct file
	structFile, err := os.Create("./" + packageName + "/structs.go")
	if err != nil {
		log.Fatalf("Error creating struct file: %v", err)
	}
	defer structFile.Close()

	var structBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&structBuf, "structs", openAPIStructs)
	if err != nil {
		log.Fatalf("Error executing struct template: %v", err)
	}
	formattedStruct, err := format.Source(structBuf.Bytes())
	if err != nil {
		log.Fatalf("Error formatting struct file: %v", err)
	}
	_, err = structFile.Write(formattedStruct)
	if err != nil {
		log.Fatalf("Error writing formatted struct file: %v", err)
	}

	// Generate handler file
	handlerFile, err := os.Create("./" + packageName + "/handlers.go")
	if err != nil {
		log.Fatalf("Error creating handler file: %v", err)
	}
	defer handlerFile.Close()

	var handlerBuf bytes.Buffer
	err = tmpl.ExecuteTemplate(&handlerBuf, "handlers", openAPIStructs)
	if err != nil {
		log.Fatalf("Error executing handler template: %v", err)
	}
	formattedHandler, err := format.Source(handlerBuf.Bytes())
	if err != nil {
		log.Fatalf("Error formatting handler file: %v", err)
	}
	_, err = handlerFile.Write(formattedHandler)
	if err != nil {
		log.Fatalf("Error writing formatted handler file: %v", err)
	}
}
