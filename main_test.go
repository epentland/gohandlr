package main

import (
	"context"
	"testing"

	"github.com/epentland/twirp/handle"
	"github.com/go-chi/chi/v5"
)

func TestFunc(t *testing.T) {
	r := chi.NewRouter()
	hand, err := handle.NewHandleStruct[handle.Nil, handle.Nil, handle.Nil]()
	if err != nil {
		panic(err)
	}

	hand.Handle(r.Post, "/user/{id}", func(ctx context.Context, e1, e2 handle.Nil) (handle.Nil, error) {
		return handle.Nil{}, nil
	})
}
