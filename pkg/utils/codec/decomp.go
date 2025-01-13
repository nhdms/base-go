package codec

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Compress -
func Compress(input []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(input); err != nil {
		return input
	}
	if err := gz.Flush(); err != nil {
		return input
	}
	if err := gz.Close(); err != nil {
		return input
	}
	return b.Bytes()
}

// Decompress -
func Decompress(input []byte) []byte {
	if input == nil {
		return []byte{}
	}
	br := bytes.NewReader(input)
	gz, err := gzip.NewReader(br)
	if err != nil {
		return []byte{}
	}
	out, err := ioutil.ReadAll(gz)
	if err != nil {
		return []byte{}
	}
	return out
}
