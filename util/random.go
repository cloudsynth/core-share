package util

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/oklog/ulid/v2"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func GenerateRandomString(randBytes int) (string, error) {
	b, err := GenerateRandomBytes(randBytes)
	return base64.URLEncoding.EncodeToString(b), err
}

// TODO: WARNING => Javascript regex ids rely on it/this....
func RandomULID() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

func IsULID(data string) bool {
	if data == "" {
		return false
	}
	_, err := ulid.Parse(data)
	return err == nil
}
