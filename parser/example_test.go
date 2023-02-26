// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/leodido/go-conventionalcommits"
)

func output(out interface{}) {
	spew.Config.DisableCapacities = true
	spew.Config.DisablePointerAddresses = true
	spew.Config.SortKeys = true
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

func Example_best_effort() {
	i := []byte(`fix: description
a blank line is mandatory to start the body part of the commit message!`)
	m, e := NewMachine(WithBestEffort()).Parse(i)
	output(m)
	fmt.Println(e)
	// Output:
	// (*conventionalcommits.ConventionalCommit)({
	//  Type: (string) (len=3) "fix",
	//  Description: (string) (len=11) "description",
	//  Scope: (*string)(<nil>),
	//  Exclamation: (bool) false,
	//  Body: (*string)(<nil>),
	//  Footers: (map[string][]string) <nil>
	// })
	// missing a blank line: col=17
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

// fixme > flaky because of the footer map keys.
func Example_full_conventional() {
	i := []byte(`fix: correct minor typos in code

see the issue [0] for details
on typos fixed.

[0]: https://issue

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
	//  Body: (*string)((len=65) "see the issue [0] for details\non typos fixed.\n\n[0]: https://issue"),
	//  Footers: (map[string][]string) (len=2) {
	//   (string) (len=4) "refs": ([]string) (len=1) {
	//    (string) (len=3) "133"
	//   },
	//   (string) (len=11) "reviewed-by": ([]string) (len=1) {
	//    (string) (len=1) "Z"
	//   }
	//  }
	// })
}

func Example_breaking_freeformtype_with_scope() {
	i := []byte(`KVM(nVMX)!: Truncate base/index GPR value on address calc in !64-bit`)
	m, _ := NewMachine(WithTypes(conventionalcommits.TypesFreeForm)).Parse(i)
	output(m)
	// Output:
	// (*conventionalcommits.ConventionalCommit)({
	//  Type: (string) (len=3) "kvm",
	//  Description: (string) (len=56) "Truncate base/index GPR value on address calc in !64-bit",
	//  Scope: (*string)((len=4) "nvmx"),
	//  Exclamation: (bool) true,
	//  Body: (*string)(<nil>),
	//  Footers: (map[string][]string) <nil>
	// })
}
