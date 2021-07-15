package perms

import (
	"context"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/metadata"
)

const GRPCTokenKey = "token"

func logAuthError(err error) {
	//logger.Debugf("Ran into auth error: %s", err)
}

func GetTokenFromGrpcContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	tokens, ok := md[GRPCTokenKey]
	if !ok || len(tokens) != 1 || tokens[0] == "" {
		return ""
	}
	return tokens[0]
}

func IsAtLeastLevel(readers, writers, owners []string, actor Actor, level PermissionLevel) bool {
	if actor.IsSuperUser() {
		return true
	}
	actorIdentities := append([]string{actor.IdentitySubject()}, actor.MemberSubjects()...)
	actorIdentities = funk.FilterString(actorIdentities, func(s string) bool { return s != "" })
	// case where empty strings are returned by the find actor
	if len(actorIdentities) == 0 {
		return false
	}
	if level == LevelSuperUser {
		return actor.IsSuperUser()
	} else if level == LevelOwner {
		return len(funk.IntersectString(owners, actorIdentities)) > 0 ||
			actor.IsSuperUser()

	} else if level == LevelWriter {
		return len(funk.IntersectString(owners, actorIdentities)) > 0 ||
			len(funk.IntersectString(writers, actorIdentities)) > 0 ||
			actor.IsSuperUser()
	} else if level == LevelReader {
		return len(funk.IntersectString(owners, actorIdentities)) > 0 ||
			len(funk.IntersectString(writers, actorIdentities)) > 0 ||
			len(funk.IntersectString(readers, actorIdentities)) > 0 ||
			actor.IsSuperUser()
	}
	return false
}
