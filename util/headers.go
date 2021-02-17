package util

import (
	"context"
	"github.com/cloudsynth/core-share/coretypes"
	"google.golang.org/grpc/metadata"
	"strings"
	"net/http"
)

func GRPCMDToHTTPHeaders(md metadata.MD) http.Header {
	httpHeaders := http.Header{} // Normalizes headers for us.
	for key, values := range md {
		if len(values) == 0 {
			continue
		}
		httpHeaders.Set(key, values[0])
	}
	return httpHeaders
}

func IPV4FromHeaders(headers http.Header) string {
	// source: https://github.com/un33k/django-ipware
	var headerIpSources = []string{
		"HTTP-X-FORWARDED-FOR",
		"X-FORWARDED-FOR",
		"HTTP-CLIENT-IP",
		"HTTP-X-REAL-IP",
		"HTTP-X-FORWARDED",
		"HTTP-X-CLUSTER-CLIENT-IP",
		"HTTP-FORWARDED-FOR",
		"HTTP-FORWARDED",
		"HTTP-VIA",
		"REMOTE-ADDR",
	}

	for _, headerSource := range headerIpSources {
		result := headers.Get(headerSource)
		if result != "" {
			resultParsed := strings.Split(result, ",")
			if len(resultParsed) != 0 {
				return strings.TrimSpace(resultParsed[0])
			}
		}
	}
	return ""
}

func HttpHeadersToRequestMetaHeaders(h http.Header) coretypes.Params {
	finalHeaders := coretypes.Params{}
	for key, values := range h {
		if len(values) > 0 {
			finalHeaders	= append(finalHeaders, coretypes.KVPair{Key: key, Value: values[0]})
		}
	}
	return finalHeaders
}


func IncomingToOutgoingContext(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(context.Background(), md)
}
