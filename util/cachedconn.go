package util

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/kofalt/go-memoize"
	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strings"
)


var connMemoizer = memoize.NewMemoizer(cache.NoExpiration, -1)

func CachedGrpcConn(endpoint string) (*grpc.ClientConn, error) {
	connI, err, _ := connMemoizer.Memoize("grpc" + endpoint, func() (interface{}, error) {
		if strings.HasPrefix(endpoint, "/"){
		}
		return grpc.Dial(endpoint, grpc.WithInsecure())
	})
	if err != nil {
		return nil, err
	}
	return connI.(*grpc.ClientConn), err
}


func CachedDbConn(dialect, connString string) (*gorm.DB, error) {
	connI, err, _ := connMemoizer.Memoize("gorm" + dialect+connString, func() (interface{}, error) {
		if dialect != "postgres"{
			return nil, errors.New("unsupported dialect")
		}
		dbConn, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
		if err != nil{
			return nil, err
		}
		db, err := dbConn.DB()
		if err != nil{
			errors.Wrap(err, "unable to grab underlying connection")
		}
		db.SetMaxOpenConns(8)
		db.SetMaxIdleConns(8)
		// Ensure bulk removes don't happen
		// https://github.com/go-gorm/gorm/issues/1152
		//dbConn = dbConn.BlockGlobalUpdate(true) // Not needed in new gorm
		return dbConn, nil
	})
	if err != nil {
		return nil, err
	}
	return connI.(*gorm.DB), err
}

func CachedTransport(endpoint string) (*http.Transport, error) {
	connI, err, _ := connMemoizer.Memoize("http" + endpoint, func() (interface{}, error) {
		return &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				network := "tcp"
				if strings.HasPrefix(endpoint, "/"){
					network = "unix"
				}
				return net.Dial(network, endpoint)
			},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return connI.(*http.Transport), err
}

