package net

import (
	"context"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func GrpcLoggingInterceptors() (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor) {
	unary := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		v, err := handler(ctx, req)
		end := time.Now()
		status := 200
		logPrefix := log.Info()
		if err != nil{
			status = 500
			logPrefix = log.Error().Err(err)
		}
		logPrefix.
			Str("path", info.FullMethod).
			Int("code", status).
			Dur("timeMs", end.Sub(start)).
			Str("method", "GRPC").Send()
		return v, err
	}

	stream := func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, stream)
		end := time.Now()
		status := 200
		logPrefix := log.Info()
		if err != nil{
			status = 500
			logPrefix = log.Error().Err(err)
		}
		logPrefix.
			Str("path", info.FullMethod).
			Int("code", status).
			Dur("timeMs", end.Sub(start)).
			Str("method", "GRPC").Send()
		return err
	}
	return unary, stream
}

func HttpLoggingHandler(f http.HandlerFunc) http.Handler {
	c := alice.New()
	c = c.Append(hlog.NewHandler(log.Logger))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		if r.Method == "PRI" { // Http2 passthrough headers
			return
		}
		hlog.FromRequest(r).Info().
			Str("path", r.URL.Path).
			Int("code", status).
			Dur("timeMs", duration).
			Str("method", r.Method).Msg("")
	}))
	//c = c.Append(hlog.RemoteAddrHandler("ip"))
	//c = c.Append(hlog.UserAgentHandler("user_agent"))
	//c = c.Append(hlog.RefererHandler("referer"))

	// Here is your final handler
	h := c.Then(f)
	return h
}
