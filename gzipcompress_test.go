package main

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/klauspost/pgzip"
)

func TestGzipCompress(t *testing.T) {
	compressed := &bytes.Buffer{}
	var input = []byte("hello world")
	log.Println(input)
	w := pgzip.NewWriter(compressed)
	defer w.Close()
	var _, err = io.Copy(w, bytes.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	w.Flush()
	w.Close()
	x := compressed.Bytes()
	log.Println(x)

	reader, err := pgzip.NewReader(bytes.NewReader(x))
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	var buf *bytes.Buffer = &bytes.Buffer{}
	_, err = io.Copy(buf, reader)
	reader.Close()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(buf.Bytes())
}
