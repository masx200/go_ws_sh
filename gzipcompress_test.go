package main

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/klauspost/pgzip"
	"github.com/zeebo/assert"
)

func TestGzipCompress(t *testing.T) {
	compressed := &bytes.Buffer{}
	var input = []byte("hello world")
	log.Println("input:",input)
	// log.Println(input)
	w := pgzip.NewWriter(compressed)
	defer w.Close()
	var _, err = io.Copy(w, bytes.NewReader(input))
	if err != nil {
		t.Fatal(err)
		return
	}

	w.Flush()
	w.Close()
	x := compressed.Bytes()
	log.Println(
		"compressed:",
		x)

	reader, err := pgzip.NewReader(bytes.NewReader(x))
	if err != nil {
		t.Fatal(err)
		return
	}
	defer reader.Close()
	var buf *bytes.Buffer = &bytes.Buffer{}
	_, err = io.Copy(buf, reader)
	reader.Close()
	if err != nil {
		t.Fatal(err)
		return
	}
	log.Println(
		"decompressed:",
		buf.Bytes())
	assert.Equal(t, input, buf.Bytes())
}
