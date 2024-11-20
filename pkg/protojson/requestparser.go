package protojson

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fullstorydev/grpcurl"

	// nolint
	"github.com/golang/protobuf/jsonpb"
)

// NewRequestParser creates a new request parser from the given http.Request and resolver.
func NewRequestParser(r *http.Request, pathName []string, resolver jsonpb.AnyResolver) (grpcurl.RequestParser, error) {
	params := make(map[string]any)

	for _, v := range pathName {
		params[v] = r.PathValue(v)
	}

	body, ok := getBody(r)
	if !ok {
		return buildJsonRequestParser(params, resolver)
	}

	if len(params) == 0 {
		return grpcurl.NewJSONRequestParser(body, resolver), nil
	}

	m := make(map[string]any)
	if err := json.NewDecoder(body).Decode(&m); err != nil && err != io.EOF {
		return nil, err
	}

	for k, v := range params {
		m[k] = v
	}

	return buildJsonRequestParser(m, resolver)
}

func buildJsonRequestParser(m map[string]any, resolver jsonpb.AnyResolver) (
	grpcurl.RequestParser, error,
) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}

	return grpcurl.NewJSONRequestParser(&buf, resolver), nil
}

func getBody(r *http.Request) (io.Reader, bool) {
	if r.Body == nil {
		return nil, false
	}

	if r.ContentLength == 0 {
		return nil, false
	}

	if r.ContentLength > 0 {
		return r.Body, true
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		return nil, false
	}

	if buf.Len() > 0 {
		return &buf, true
	}

	return nil, false
}
