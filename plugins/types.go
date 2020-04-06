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

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Params []KVPair

func (p Params) Get(key string) (value string, ok bool) {
	for _, param := range p {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}

type Config struct {
	DebugDbQueries     bool          `json:"debug_db_queries"`
	DbConnectionString string        `json:"db_connection_string"`
	SuperuserPskToken  string        `json:"superuser_psk_token"`
	JwtProviders       []JWTProvider `json:"jwt_providers"`
	JwkProviders       []JWKProvider `json:"jwk_providers"`
	AppConfig          Params        `json:"app_config"`
}

type PluginMakeHandler func(config Config, findActor perms.GRPCActorFinder) (*grpc.Server, http.HandlerFunc, error)
