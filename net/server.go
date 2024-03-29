package net

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
	"time"
)

func ServeGrpcAndHttpMultiplexed(addr string, grpcServer *grpc.Server, httpHandler http.Handler, defaultTimeout time.Duration) error {
	log.Info().Msgf("starting listener on %s", addr)
	network := "tcp"
	if strings.HasPrefix(addr, "/"){
		network = "unix"
	}
	mainListener, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	multiplexer := cmux.New(mainListener)
	multiplexer.SetReadTimeout(defaultTimeout)
	multiplexer.HandleError(func(err error) bool {
		if err != nil {
			log.Error().Err(err).Send()
		}
		return true
	})

	grpcListener := multiplexer.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc+proto"),
	)
	httpListener := multiplexer.Match(cmux.Any())

	errChan := make(chan error)
	httpServer := http.Server{
		Handler:           http.TimeoutHandler(httpHandler, defaultTimeout, "timeout"),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	go func() {
		err := grpcServer.Serve(grpcListener)
		errChan <- errors.Wrap(err, "error serving grpc")
	}()
	go func() {
		err := httpServer.Serve(httpListener)
		errChan <- errors.Wrap(err, "error serving http")
	}()
	go func() {
		err := multiplexer.Serve()
		errChan <- errors.Wrap(err, "error serving cmux")
	}()

	closeAll := func() {
		grpcServer.Stop()
		httpServer.Close()
		httpListener.Close()
		grpcListener.Close()
		mainListener.Close()
	}

	// Blocking
	select {
	case err := <-errChan:
		closeAll()
		return err
	}
}
