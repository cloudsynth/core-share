package util

import (
	"github.com/kofalt/go-memoize"
	"google.golang.org/grpc"
	"github.com/jinzhu/gorm"
	"time"
	"strings"
	"net"
)

var rpcConnMemoizer = memoize.NewMemoizer(time.Hour, time.Hour)
func CachedGrpcConn(endpoint string) (*grpc.ClientConn, error) {
	connI, err, _ := rpcConnMemoizer.Memoize(endpoint, func() (interface{}, error) {
		dialOpts := []grpc.DialOption{grpc.WithInsecure()}
		if strings.HasPrefix(endpoint, "/"){
				dialOpts = append(dialOpts, grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
					return net.DialTimeout("unix", addr, timeout)
				}))
		}
		return grpc.Dial(endpoint, dialOpts...)
	})
	if err != nil {
		return nil, err
	}
	return connI.(*grpc.ClientConn), err
}


var connMemoizer = memoize.NewMemoizer(time.Hour, time.Hour)
func CachedDbConn(dialect, connString string) (*gorm.DB, error) {
	connI, err, _ := connMemoizer.Memoize(dialect+connString, func() (interface{}, error) {
		dbConn, err := gorm.Open(dialect, connString)
		if err != nil{
			return nil, err
		}
		dbConn.DB().SetMaxOpenConns(8)
		dbConn.DB().SetMaxIdleConns(8)
		// Ensure bulk removes don't happen
		// https://github.com/go-gorm/gorm/issues/1152
		dbConn = dbConn.BlockGlobalUpdate(true)
		return dbConn, nil
	})
	if err != nil {
		return nil, err
	}
	return connI.(*gorm.DB), err
}


