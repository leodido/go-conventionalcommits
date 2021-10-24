# go-conventionalcommits

[![Build](https://img.shields.io/circleci/build/github/leodido/go-conventionalcommits/develop?style=for-the-badge)](https://app.circleci.com/pipelines/github/leodido/go-conventionalcommits) [![Coverage](https://img.shields.io/codecov/c/github/leodido/go-conventionalcommits/develop?style=for-the-badge)](https://codecov.io/gh/leodido/go-conventionalcommits) [![License](https://img.shields.io/github/license/leodido/go-conventionalcommits?style=for-the-badge)](LICENSE) [![Go Report](https://goreportcard.com/badge/github.com/leodido/go-conventionalcommits?style=for-the-badge)](https://goreportcard.com/report/github.com/leodido/go-conventionalcommits)

**A parser for [Conventional Commits v1.0](https://www.conventionalcommits.org/en/v1.0.0/#specification) commit messages**.

> Fu powers to parse your commits!

This repository provides a library to parse your commit messages according to the Conventional Commits v1.0 specification.

## Installation

```console
go get github.com/leodido/go-conventionalcommits
```

## Docs

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/leodido/go-conventionalcommits)

The [parser/docs](parser/docs/) directory contains `.dot` and `.png` files representing the finite-state machines (FSMs) implementing the parser.

## Usage

### Parse

Your code base uses only single line commit messages like this one?

```console
feat: awesomeness
```

No problem at all since the body and the footer parts are not mandatory:

```go
m, _ := parser.NewMachine().Parse([]byte(`feat: awesomeness`))
```

### Full conventional commit messages

Imagine you have a commit message like this:

```console
docs: correct minor typos

see the issue for details

on docs edits.

Reviewed-by: Z
Refs #133
```

Go with this:

```go
opts := []conventionalcommits.MachineOption{
    WithTypes(conventionalcommits.TypesConventional),
}
res, err := parser.NewMachine(opts...).Parse(i)
```

Or, more simpler:

```go
res, err := parser.NewMachine(WithTypes(conventionalcommits.TypesConventional)).Parse(i)
```

### Types

This library provides support for different types sets:

- **minimal**: fix, feat
- **conventional**: build, ci, chore, docs, feat, fix, perf, refactor, revert, style, test
- **falco**: build, ci, chore, docs, feat, fix, perf, new, revert, update, test, rule

At the moment, those types are at build time. Which means users can't configure them at runtime.

Anyway, there's also a **free-form** types set that accepts any combination of printable characters (before the separator after which the commit description starts) as a valid type.

You can choose the type set passing the `WithTypes(conventionalcommits.TypesConventional)` option as shown above.

### Options

A parser behaviour is configurable by using options.

You can set them calling a function on the parser machine.

```go
p := parser.NewMachine()
p.WithBestEffort()
res, err := p.Parse(i)
```

Or you can provide options to `NewMachine(...)` directly.

```go
p := parser.NewMachine(WithBestEffort())
res, err := p.Parse(i)
```

### Best effort

The best effort mode will make the parser return what it found until the point it errored out,
if it found (at least) a valid type and a valid description.

Let's make an example.

Suppose this input commit message:

```console
fix: description
a blank line is mandatory to start the body part of the commit message!
```

The input does not respect the Conventional Commits v1 specification because it lacks a blank line after the description (before the body).

Anyways, if the parser you're using has the best effort mode enabled, you can still obtain some structured data since at least a valid type and description have been found!

```go
res, err := parser.NewMachine(WithBestEffort()).Parse(i)
```

The result will contain a `ConventionalCommit` struct instance with the `Type` and the `Description` fields populated and ignore the rest after the error column.

The parser will still return the error (with the position information), so that you can eventually use it.

## Performances

To run the benchmark suite execute the following command.

```console
make bench
```

All the parsers have the best effort mode on.

On my machine<sup>[1](#mymachine)</sup>, these are the results for the `slim` parser with the default - ie., `minimal`, commit message types.

```console
[ok]_minimal______________________________________-12          4876018       242 ns/op     147 B/op       5 allocs/op
[ok]_minimal_with_scope___________________________-12          4258562       284 ns/op     163 B/op       6 allocs/op
[ok]_minimal_breaking_with_scope__________________-12          4176747       288 ns/op     163 B/op       6 allocs/op
[ok]_full_with_50_characters_long_description_____-12          1661618       700 ns/op     288 B/op      10 allocs/op
[no]_empty________________________________________-12          4059327       292 ns/op     112 B/op       3 allocs/op
[no]_type_but_missing_colon_______________________-12          2701904       444 ns/op     200 B/op       6 allocs/op
[no]_type_but_missing_description_________________-12          2207985       539 ns/op     288 B/op       8 allocs/op
[no]_type_and_scope_but_missing_description_______-12          1969390       605 ns/op     312 B/op      10 allocs/op
[no]_breaking_with_type_and_scope_but_missing_desc-12          1978302       606 ns/op     312 B/op      10 allocs/op
[~~]_newline_in_description_______________________-12          2115649       563 ns/op     232 B/op      11 allocs/op
[no]_missing_whitespace_in_description____________-12          1997863       595 ns/op     312 B/op      10 allocs/op
```

Using another set of commit message types, for example the `conventional` one, does not have any noticeable impact on performances, as you can see below.

```console
[ok]_minimal______________________________________-12          5297486       228 ns/op     147 B/op       5 allocs/op
[ok]_minimal_with_scope___________________________-12          4498694       267 ns/op     163 B/op       6 allocs/op
[ok]_minimal_breaking_with_scope__________________-12          4431040       273 ns/op     163 B/op       6 allocs/op
[ok]_full_with_50_characters_long_description_____-12          1750111       692 ns/op     288 B/op      10 allocs/op
[no]_empty________________________________________-12          3996532       294 ns/op     112 B/op       3 allocs/op
[no]_type_but_missing_colon_______________________-12          2657913       451 ns/op     200 B/op       6 allocs/op
[no]_type_but_missing_description_________________-12          2172524       553 ns/op     288 B/op       8 allocs/op
[no]_type_and_scope_but_missing_description_______-12          1880526       637 ns/op     312 B/op      10 allocs/op
[no]_breaking_with_type_and_scope_but_missing_desc-12          1879779       635 ns/op     312 B/op      10 allocs/op
[~~]_newline_in_description_______________________-12          2023514       592 ns/op     232 B/op      11 allocs/op
[no]_missing_whitespace_in_description____________-12          1883124       623 ns/op     312 B/op      10 allocs/op
```

As you may notice, this library is very fast at what it does.

Parsing a commit goes from taking about the same amount of time (~299ns) the half-life of polonium-212 takes<sup>[2](#nanosecondwiki)</sup> to less than a microsecond.

---

- <a name="mymachine">[1]</a>: Intel Core i7-8850H CPU @ 2.60GHz
- <a name="nanosecondwiki">[2]</a>: <https://en.wikipedia.org/wiki/nanosecond>

---

[![Analytics](https://ga-beacon.appspot.com/UA-49657176-1/go-conventionalcommits?flat)](https://github.com/igrigorik/ga-beacon)
