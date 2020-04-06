package plugins

import (
	"github.com/cloudsynth/core-share/perms"
	"google.golang.org/grpc"
	"net/http"
)

type JWTProvider struct {
	Key            string `json:"key"`
	ExpectedIssuer string `json:"expected_issuer"`
}

type JWKProvider struct {
	JWKURI           string   `json:"jwk_uri"`
	ExpectedIssuer   string   `json:"expected_issuer"`
	ExpectedAudience []string `json:"expcted_audience"`
	ExpectedAlgo     string   `json:"expected_algo"`
}

type Config struct {
	DebugDbQueries    bool
	DbDialect         string
	DbDialectArgs     []string
	SuperUserPSKToken string
	JWTProviders      []JWTProvider
	JWKProviders      []JWKProvider
	AppConfig		  map[string]string
}

type PluginMakeHandler func(config Config, findActor  perms.GRPCActorFinder) (*grpc.Server, http.HandlerFunc, error)
