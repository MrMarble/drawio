package drawio

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"net/url"
	"strings"
)

func decode(data string) string {
	// decode base64 string
	uEnc, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}

	inflated := inflate(uEnc)

	// url decode
	decoded, err := url.QueryUnescape(string(inflated))
	if err != nil {
		panic(err)
	}

	return decoded
}

func encode(data string) string {
	// url encode
	encoded := url.QueryEscape(data)
	encoded = strings.Replace(encoded, "+", "%20", -1)

	deflated := deflate(encoded)

	// Base 64 encode
	return base64.StdEncoding.EncodeToString(deflated)
}

// Inflate utility that decompresses a string using the flate algo
func inflate(deflated []byte) []byte {
	var b bytes.Buffer
	r := flate.NewReader(bytes.NewReader(deflated))
	b.ReadFrom(r)
	r.Close()
	return b.Bytes()
}

// Deflate utility that compresses a string using the flate algo
func deflate(inflated string) []byte {
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, -1)
	w.Write([]byte(inflated))
	w.Close()
	return b.Bytes()
}
