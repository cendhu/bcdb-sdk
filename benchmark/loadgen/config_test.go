package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig("./config.yml")
	require.NoError(t, err)
	spew.Dump(c)
}
