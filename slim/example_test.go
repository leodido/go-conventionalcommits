package slim

import (
	"fmt"
	"strconv"
	"strings"

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
	// (*slim.ConventionalCommit)({
	//  Minimal: (conventionalcommits.Minimal) {
	//   Type: (string) (len=3) "fix",
	//   Description: (string) (len=9) "something",
	//   Scope: (*string)(<nil>),
	//   Exclamation: (bool) true
	//  }
	// })
	// there are breaking changes? true
}

func Example_conventional_ignoringbody() {
	i := []byte(`fix: correct minor typos in code

see the issue for details
on typos fixed.

Reviewed-by: Z
Refs #133`)

	opts := []conventionalcommits.MachineOption{
		WithBestEffort(),
		WithTypes(conventionalcommits.TypesConventional),
	}
	m, e := NewMachine(opts...).Parse(i)
	output(m)
	fmt.Println("is result ok?", m.Ok())

	errstr := e.Error()
	fmt.Println(errstr)
	pos := strings.LastIndex(errstr, "=")
	num, _ := strconv.Atoi(errstr[pos+1 : len(errstr)])
	// Not checking pos and num because ain't time for bs
	fmt.Printf("parsing ok until position %d\n", num)
	fmt.Println("ignored body:")
	fmt.Println(string(i[num:len(i)]))

	// Output:
	// (*slim.ConventionalCommit)({
	//  Minimal: (conventionalcommits.Minimal) {
	//   Type: (string) (len=3) "fix",
	//   Description: (string) (len=27) "correct minor typos in code",
	//   Scope: (*string)(<nil>),
	//   Exclamation: (bool) false
	//  }
	// })
	// is result ok? true
	// illegal newline: col=33
	// parsing ok until position 33
	// ignored body:
	//
	// see the issue for details
	// on typos fixed.
	//
	// Reviewed-by: Z
	// Refs #133
}
