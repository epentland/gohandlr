package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewMux()
	r.Use(middleware.Logger)

	err := http.ListenAndServe(":8083", r)
	if err != nil {
		fmt.Println(err)
	}
}
