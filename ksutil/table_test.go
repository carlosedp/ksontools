package ksutil

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(&buf)

	table.SetHeader([]string{"name", "version", "Namespace", "SERVER"})
	table.Append([]string{"default", "v1.7.0", "default", "http://default"})
	table.AppendBulk([][]string{
		{"dev", "v1.8.0", "dev", "http://dev"},
		{"east/prod", "v1.8.0", "east/prod", "http://east-prod"},
	})

	table.Render()

	b, err := ioutil.ReadFile("testdata/table/table.txt")
	require.NoError(t, err)

	assert.Equal(t, string(b), buf.String())
}

func TestTable_no_header(t *testing.T) {
	var buf bytes.Buffer
	table := NewTable(&buf)

	table.Append([]string{"default", "v1.7.0", "default", "http://default"})
	table.AppendBulk([][]string{
		{"dev", "v1.8.0", "dev", "http://dev"},
		{"east/prod", "v1.8.0", "east/prod", "http://east-prod"},
	})

	table.Render()

	b, err := ioutil.ReadFile("testdata/table/table_no_header.txt")
	require.NoError(t, err)

	assert.Equal(t, string(b), buf.String())
}
