package httpserver

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

var errUnsupportedContentEncoding = errors.New("unsupported Content-Encoding")

func readRequestBody(r *http.Request) ([]byte, error) {
	var body io.Reader = r.Body
	switch encoding := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Encoding"))); encoding {
	case "", "identity":
	case "gzip":
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, fmt.Errorf("decode gzip request: %w", err)
		}
		defer reader.Close()
		body = reader
	default:
		return nil, fmt.Errorf("%w: %q", errUnsupportedContentEncoding, encoding)
	}

	return io.ReadAll(body)
}

func writeResponse(w http.ResponseWriter, r *http.Request, code int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Vary", "Accept-Encoding")

	var body io.Writer = w
	if acceptsGzip(r.Header) {
		w.Header().Set("Content-Encoding", "gzip")
		writer := gzip.NewWriter(w)
		defer writer.Close()
		body = writer
	}

	w.WriteHeader(code)
	json.NewEncoder(body).Encode(response)
}

func acceptsGzip(header http.Header) bool {
	wildcard := false
	for value := range strings.SplitSeq(strings.Join(header.Values("Accept-Encoding"), ","), ",") {
		encoding, params, err := mime.ParseMediaType(strings.TrimSpace(value))
		if err != nil {
			continue
		}

		quality := 1.0
		if value, ok := params["q"]; ok {
			quality, err = strconv.ParseFloat(value, 64)
			if err != nil || !(quality >= 0 && quality <= 1) {
				continue
			}
		}

		switch encoding {
		case "gzip":
			return quality > 0
		case "*":
			wildcard = quality > 0
		}
	}

	return wildcard
}
