package server

import (
	"net/http"
	"net/textproto"
	"strings"

	"google.golang.org/grpc/metadata"
)

// httpHeadersToGRPCMetadata converts HTTP headers to gRPC metadata.
func httpHeadersToGRPCMetadata(headers http.Header) metadata.MD {
	grpcMetadata := metadata.MD{}
	for key, values := range headers {
		grpcKey, ok := DefaultHeaderMatcher(key)
		if ok {
			for _, value := range values {
				grpcMetadata.Append(grpcKey, value)
			}
		}
	}
	return grpcMetadata
}

const (
	MetadataHeaderPrefix = "Grpc-Metadata-"
	MetadataPrefix       = "grpcgateway-"
)

func DefaultHeaderMatcher(key string) (string, bool) {
	switch key = textproto.CanonicalMIMEHeaderKey(key); {
	case isPermanentHTTPHeader(key):
		return MetadataPrefix + key, true
	case strings.HasPrefix(key, MetadataHeaderPrefix):
		return key[len(MetadataHeaderPrefix):], true
	}
	return "", false
}

// isPermanentHTTPHeader checks whether hdr belongs to the list of
// permanent request headers maintained by IANA.
// http://www.iana.org/assignments/message-headers/message-headers.xml
func isPermanentHTTPHeader(hdr string) bool {
	switch hdr {
	case
		"Accept",
		"Accept-Charset",
		"Accept-Language",
		"Accept-Ranges",
		"Authorization",
		"Cache-Control",
		"Content-Type",
		"Cookie",
		"Date",
		"Expect",
		"From",
		"Host",
		"If-Match",
		"If-Modified-Since",
		"If-None-Match",
		"If-Schedule-Tag-Match",
		"If-Unmodified-Since",
		"Max-Forwards",
		"Origin",
		"Pragma",
		"Referer",
		"User-Agent",
		"Via",
		"Warning":
		return true
	}
	return false
}
