package util

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
	"testing"
)

func TestRandomIDBrowserSafe(t *testing.T) {
	seenUlrs := map[string]bool{}
	for i := 0; i < 10000; i++ {
		result := RandomULID()
		_, ok := seenUlrs[result]
		require.False(t, ok)
		require.True(t, len(result) >= 20)
		// Make sure raw encoded and doesn't contain padding
		require.False(t, strings.Contains(result, "="))
		// For objtree paths we want to make sure that they don't contain a dash
		require.False(t, strings.Contains(result, "-"))
		seenUlrs[result] = true
	}
}

func TestIsRandomID(t *testing.T) {
	failExamples := []string{
		"",
		"foo",
		"=",
		"bazzz",
		// Mast be RAWUrlEncoded i.e no padding
		"T2-dMAOST8mC6Wr9dU3KcQ==",
	}
	for i := 0; i < 1000; i++ {
		randString, err := GenerateRandomString(20 + rand.Intn(5))
		require.Nil(t, err)
		failExamples = append(failExamples, randString)
	}

	for _, ex := range failExamples {
		result := IsULID(ex)
		require.False(t, result)

	}

	for i := 0; i < 1000; i++ {
		result := IsULID(RandomULID())
		require.True(t, result)
	}

	validExamples := []string{
		// Padded and unpadded variants
		"01D78XYFJ1PRM1WPBCBT3VHMNV",
		"01DCHHEYY4VREC6EJ7KSYDKCNV",
	}

	for _, valid := range validExamples {
		result := IsULID(valid)
		require.True(t, result, "%s must is valid uuid", valid)

	}
}
