package n3_test

import (
	"embed"
	_ "embed"
	"fmt"
	"github.com/0x51-dev/rdf/n3"
	"testing"
)

var (
	//go:embed testdata/example1.n3
	example1 string

	//go:embed testdata/*.n3
	examples embed.FS
)

func Example_example1() {
	doc, _ := n3.ParseDocument(example1)
	fmt.Println(doc)
	// Output:
	// <http://example.org/spiderman> <http://example.org/enemyOf> <http://example.org/green-goblin> .
	// <http://example.org/spiderman> <http://xmlns.com/foaf/0.1/name> "Spiderman" .
	// <http://example.org/green-goblin> <http://xmlns.com/foaf/0.1/name> "Green Goblin" .
}

func TestExamples(t *testing.T) {
	// Amount of triples in each example (manually counted).
	triples := []int{
		3, 1, 1, 4, 2, 2, 6, 3, 3, 3,
		3, 4, 4, 4, 6, 2, 6, 2, 2, 2,
		4, 2, 3, 4, 7, 7, 4, 5, 7, 6,
		5, 3, 3, 2, 5, 2, 4, 2,
	}
	entries, err := examples.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	em := make(map[string]string)
	for _, f := range entries {
		b, err := examples.ReadFile("testdata/" + f.Name())
		if err != nil {
			t.Fatal(err)
		}
		em[f.Name()] = string(b)
	}
	for i, n := range triples {
		t.Run(fmt.Sprintf("example%d", i+1), func(t *testing.T) {
			raw := em[fmt.Sprintf("example%d.n3", i+1)]
			doc, err := n3.ParseDocument(raw)
			if err != nil {
				t.Fatal(err)
			}
			if len(doc) != n {
				t.Fatal(n, len(doc))
			}

			{ // fmt.Stringer
				doc2, err := n3.ParseDocument(doc.String())
				if err != nil {
					t.Fatal(doc.String())
				}
				if len(doc2) != n {
					t.Fatal(n, len(doc2))
				}
			}
		})
	}
}
