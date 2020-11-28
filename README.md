# go-conventionalcommits

[![License](https://img.shields.io/github/license/leodido/go-conventionalcommits?style=for-the-badge)](LICENSE) [![Go Report](https://goreportcard.com/badge/github.com/leodido/go-conventionalcommits?style=for-the-badge)](https://goreportcard.com/report/github.com/leodido/go-conventionalcommits)

**A parser for [Conventional Commits v1.0](https://www.conventionalcommits.org/en/v1.0.0/#specification) commit messages**.

> Fu powers to parse your commits!

This repository provides libraries to parse your commit messages according to the Conventional Commits v1.0 specification.

Wanna parse only the first line of your commits?

Use the [leodido/go-conventionalcommits/slim](slim/) package.

Wanna parse the full commit message corpus according to the specification?

Use the [leodido/go-conventionalcommits/full](full/) package (**WIP**).

## Installation

```console
go get github.com/leodido/go-conventionalcommits
```

## Docs

TBD.

## Usage

TBD.

### Options

TBD.

## Performances

To run the benchmark suite execute the following command.

```console
make bench
```

All the parsers have the best effort mode on.

On my machine<sup>[1](#mymachine)</sup>, these are the results for the `slim` parser with the default - ie., `minimal`, commit message types.

```
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

```
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

---

* <a name="mymachine">[1]</a>: Intel Core i7-8850H CPU @ 2.60GHz

---

[![Analytics](https://ga-beacon.appspot.com/UA-49657176-1/go-conventionalcommits?flat)](https://github.com/igrigorik/ga-beacon)