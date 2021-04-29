package parser

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/leodido/go-conventionalcommits"
)

func output(out interface{}) {
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerAddresses = true
	spew.Dump(out)
}

func Example_minimal_withoutbody() {
	i := []byte(`fix!: something`)
	m, _ := NewMachine().Parse(i)
	output(m)
	fmt.Println("there are breaking changes?", m.IsBreakingChange())
	// Output:
	// (*conventionalcommits.ConventionalCommit)({
	//  Type: (string) (len=3) "fix",
	//  Description: (string) (len=9) "something",
	//  Scope: (*string)(<nil>),
	//  Exclamation: (bool) true,
	//  Body: (*string)(<nil>),
	//  Footers: (map[string][]string) <nil>
	// })
	// there are breaking changes? true
}

func Example_multiline_body() {
	i := []byte(`fix: x

see the issue for details

but first a newline
and then two blank lines:

typos fixed.`)
	m, _ := NewMachine().Parse(i)
	output(m)
	// Output:
	// (*conventionalcommits.ConventionalCommit)({
	//  Type: (string) (len=3) "fix",
	//  Description: (string) (len=1) "x",
	//  Scope: (*string)(<nil>),
	//  Exclamation: (bool) false,
	//  Body: (*string)((len=86) "see the issue for details\n\nbut first a newline\nand then two blank lines:\n\ntypos fixed."),
	//  Footers: (map[string][]string) <nil>
	// })
}

func Example_full_conventional() {
	i := []byte(`fix: correct minor typos in code

see the issue for details
on typos fixed.

Reviewed-by: Z
Refs #133`)
	opts := []conventionalcommits.MachineOption{
		WithTypes(conventionalcommits.TypesConventional),
	}
	m, _ := NewMachine(opts...).Parse(i)
	output(m)
	// Output:
	// (*conventionalcommits.ConventionalCommit)({
	//  Type: (string) (len=3) "fix",
	//  Description: (string) (len=27) "correct minor typos in code",
	//  Scope: (*string)(<nil>),
	//  Exclamation: (bool) false,
	//  Body: (*string)((len=41) "see the issue for details\non typos fixed."),
	//  Footers: (map[string][]string) (len=2) {
	//   (string) (len=11) "reviewed-by": ([]string) (len=1) {
	//    (string) (len=1) "Z"
	//   },
	//   (string) (len=4) "refs": ([]string) (len=1) {
	//    (string) (len=3) "133"
	//   }
	//  }
	// })
}
