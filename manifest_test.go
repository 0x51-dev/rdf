package rdf_test

import (
	_ "embed"
	"github.com/0x51-dev/rdf/internal/testsuite"
	"testing"
)

var (
	//go:embed ntriples/testdata/suite/manifest.ttl
	ntriples string
)

func TestLoadManifest(t *testing.T) {
	if _, err := testsuite.LoadManifest(ntriples); err != nil {
		t.Fatal(err)
	}
}
