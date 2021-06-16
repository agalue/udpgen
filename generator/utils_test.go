package generator

import (
	"log"
	"regexp"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseSimple(t *testing.T) {
	assert.Equal(t, "ABC 123", parseText("ABC 123"))
}

func TestParseRandomNumber(t *testing.T) {
	out := parseText("ABC [rand_n:100]")
	matches, err := regexp.Match("ABC [0-9]+", []byte(out))
	assert.NilError(t, err)
	assert.Assert(t, matches)
}

func TestParseRandomIP(t *testing.T) {
	out := parseText("ABC [rand_ipv4]")
	matches, err := regexp.Match("ABC [0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+", []byte(out))
	assert.NilError(t, err)
	assert.Assert(t, matches)
}

func TestParseComplexText(t *testing.T) {
	text := "%%SEC-6-IPACCESSLOGP: list in[rand_n:150] denied tcp [rand_ipv4]([rand_n:65536]) -> [rand_ipv4]([rand_n:65536]), 1 packet"
	out := parseText(text)
	log.Println(out)
	assert.Assert(t, strings.Contains(out, "rand_") == false)
}
