package main

import (
	"context"
	"testing"

	"github.com/epentland/twirp/handle"
	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	ctx := handle.Context[ProcessDataInput, ProcessDataParams]{
		Context: context.Background(),
		Body:    ProcessDataInput{},
		Params:  ProcessDataParams{},
	}

	assert.NotNil(t, ctx)
}
