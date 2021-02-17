package coretypes

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"context"
	"encoding/json"
	"strings"
)

type JWTProvider struct {
	Key            string `json:"key,omitempty"`
	ExpectedIssuer string `json:"expected_issuer,omitempty"`
}

type JWKProvider struct {
	JWKURI           string   `json:"jwk_uri,omitempty"`
	ExpectedIssuer   string   `json:"expected_issuer,omitempty"`
	ExpectedAudience []string `json:"expcted_audience,omitempty"`
	ExpectedAlgo     string   `json:"expected_algo,omitempty"`
}

type KVPair struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
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

type AuthConfig struct {
	SuperuserPskToken     string        `json:"superuser_psk_token,omitempty"`
	JwtProviders          []JWTProvider `json:"jwt_providers,omitempty"`
	JwkProviders          []JWKProvider `json:"jwk_providers,omitempty"`
	EnableAnonymousUser   bool          `json:"enable_anonymous_user,omitempty"`
	EnablePublicReadPerms bool          `json:"enable_public_read_perms,omitempty"`
}

type EnabledHook struct {
	Model string `json:"model,omitempty"`
	HookType string `json:"hook_type,omitempty"`
}

type ServerConfig struct {
	DebugDbQueries     bool       `json:"debug_db_queries,omitempty"`
	DbConnectionString string     `json:"db_connection_string,omitempty"`
	DbDialect          string     `json:"db_dialect,omitempty"`
	AuthConfig         AuthConfig `json:"auth_config,omitempty"`
	SelfEndpoint       string      `json:"self_endpoint,omitempty"`
	EnabledHooks       []EnabledHook `json:"enabled_hooks,omitempty"`
	Vars               Params     `json:"vars,omitempty"`
}

func (s *ServerConfig) HookEnabled(model string, hookType string) bool {
	for _, enabledHook := range s.EnabledHooks {
		if strings.ToLower(enabledHook.Model) == strings.ToLower(model) && strings.ToLower(enabledHook.HookType) == strings.ToLower(hookType) {
			return true
		}
	}
	return false
}

type RequestMeta struct {
	ActorId             string   `json:"actor_id,omitempty"`
	ActorMemberSubjects []string `json:"actor_member_subjects,omitempty"`
	ActorIsSuperuser    bool     `json:"actor_is_superuser,omitempty"`
	Token               string   `json:"token,omitempty"`
	IpVfour             string   `json:"ip_vfour,omitempty"`
	Source              string   `json:"source,omitempty"`
	TraceId             string   `json:"trace_id,omitempty"`
	Headers             Params   `json:"headers,omitempty"`
}

type MetaHeader struct {
	ServerConfig ServerConfig `json:"server_config,omitempty"`
	RequestMeta  RequestMeta  `json:"request_meta,omitempty"`
}

const GRPCMetaHeaderKey = "core-meta"

func GetMetaFromIncomingGrpcContext(ctx context.Context) (*MetaHeader, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metaheader: no meta header set.")
	}
	values, ok := md[GRPCMetaHeaderKey]
	if !ok || len(values) != 1 || values[0] == "" {
		return nil, errors.New("metaheader: no meta header value")
	}
	metaJson := values[0]
	result := &MetaHeader{}
	err := json.Unmarshal([]byte(metaJson), result)
	if err != nil {
		return nil, errors.Wrap(err, "metaheader: uable to unmarshal header")
	}
	if result.RequestMeta.ActorId == ""{
		return nil, errors.Wrap(err, "metaheader: could not produce actor")
	}
	return result, nil
}
