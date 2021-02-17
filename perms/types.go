package perms

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

var ErrorTokenNotPassed = fmt.Errorf("api token not passed")
var ErrorTokenInvalid = fmt.Errorf("api token invalid")
var ErrorGeneric = fmt.Errorf("api token invalid")

type PermissionLevel string

const LevelOwner = PermissionLevel("owner")
const LevelWriter = PermissionLevel("writer")
const LevelReader = PermissionLevel("reader")

type Actor interface {
	IdentitySubject() string
	MemberSubjects() []string
	IsSuperUser() bool
}

type QuickActor struct {
	Identity  string
	MemberOf  []string
	SuperUser bool
}

func (q *QuickActor) IdentitySubject() string {
	return q.Identity
}

func (q *QuickActor) MemberSubjects() []string {
	return q.MemberOf
}

func (q *QuickActor) IsSuperUser() bool {
	return q.SuperUser
}

type Policy struct {
	Owners  []string
	Writers []string
	Readers []string
}

const EmailPrefix = "email|"

func EmailActor(email string) Actor {
	return &QuickActor{
		Identity: EmailPrefix + email,
	}
}

func EmailFromIdentity(identity string) (string, error) {
	if strings.HasPrefix(identity, EmailPrefix) {
		return strings.TrimPrefix(identity, EmailPrefix), nil
	}
	return "", errors.New("invalid")
}

func BotSuperUser() Actor {
	return &QuickActor{
		Identity:  "email|bot@cloudsynth.com",
		SuperUser: true,
	}
}

type ActorFinder interface {
	FindActor(token string) (Actor, error)
}
