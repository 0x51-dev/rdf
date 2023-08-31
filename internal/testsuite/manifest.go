package testsuite

import (
	"fmt"
	ttl "github.com/0x51-dev/rdf/turtle"
	"github.com/0x51-dev/rdf/turtle/grammar/ir"
)

type ApprovalType string

const (
	Approved ApprovalType = "rdft:Approved"
	Proposed ApprovalType = "rdft:Proposed"
	Rejected ApprovalType = "rdft:Rejected"
)

type Manifest struct {
	Keys    []string
	Entries map[string]*Test
}

func LoadManifest(raw string) (*Manifest, error) {
	doc, err := ttl.ParseDocument(raw)
	if err != nil {
		return nil, err
	}
	sm, err := doc.SubjectMap()
	if err != nil {
		return nil, err
	}
	manifest, ok := sm[""]
	if !ok {
		return nil, fmt.Errorf("manifest: no root")
	}
	pom, err := manifest.PredicateObjectMap()
	if err != nil {
		return nil, err
	}
	mfEntries, ok := pom["mf:entries"]
	if !ok {
		return nil, fmt.Errorf("manifest: no entries")
	}
	var keys []string
	entries := make(map[string]*Test)
	switch t := mfEntries[0].(type) {
	case ir.Collection:
		for _, entry := range t {
			var name string
			switch t := entry.(type) {
			case *ir.IRI:
				name = t.Value
			default:
				return nil, fmt.Errorf("manifest: entry name not an IRI")
			}
			keys = append(keys, name)
			entries[name], err = NewTest(sm[name])
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("manifest: entries not a collection")
	}
	return &Manifest{Keys: keys, Entries: entries}, nil
}

type Test struct {
	Type     string
	Name     string
	Comment  string
	Approval ApprovalType
	Action   string
}

func NewTest(triple ir.Triple) (*Test, error) {
	pom, err := triple.PredicateObjectMap()
	if err != nil {
		return nil, err
	}
	typ, ok := pom["rdf:type"]
	if !ok {
		// Alternate syntax used in NQuads tests.
		a, ok := pom["a"]
		if !ok {
			return nil, fmt.Errorf("test: no type")
		}
		typ = a
	}
	name, ok := pom["mf:name"]
	if !ok {
		return nil, fmt.Errorf("test: no name")
	}
	comment := pom["rdfs:comment"]
	approval := pom["rdft:approval"]
	action, ok := pom["mf:action"]
	if !ok {
		return nil, fmt.Errorf("test: no action")
	}
	return &Test{
		Type:     typ.String(),
		Name:     (name[0].(*ir.StringLiteral)).Value,
		Comment:  comment.String(),
		Approval: ApprovalType(approval.String()),
		Action:   (action[0].(*ir.IRI)).Value,
	}, nil
}
