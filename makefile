SHELL := /bin/bash
RAGEL := ragel -I common

export GO_TEST=env GOTRACEBACK=all go test $(GO_ARGS)

.PHONY: build
build: parser/machine.go
	@gofmt -w -s ./parser

.PHONY: clean
clean:
	@$(RM) parser/machine.go
	@$(RM) -R parser/docs

.PHONY: dots
dots:
	@mkdir -p parser/docs

.PHONY: docs
docs: dots parser/docs/minimal_types.png parser/docs/falco_types.png parser/docs/conventional_types.png parser/docs/free_form_types.png parser/docs/body.png parser/docs/trailer_beg.png parser/docs/trailer_end.png

.PHONY: snake2camel
snake2camel:
	@awk -i inplace '{ \
	while ( match($$0, /(.*)([a-z]+[0-9]*)_([a-zA-Z0-9])(.*)/, cap) ) \
	$$0 = cap[1] cap[2] toupper(cap[3]) cap[4]; \
	print \
	}' $(file)

parser/docs/minimal_types.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M main $< -o $@

parser/docs/minimal_types.png: parser/docs/minimal_types.dot
	dot $< -Tpng -o $@

parser/docs/falco_types.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M falco_types_main $< -o $@

parser/docs/falco_types.png: parser/docs/falco_types.dot
	dot $< -Tpng -o $@

parser/docs/conventional_types.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M conventional_types_main $< -o $@

parser/docs/conventional_types.png: parser/docs/conventional_types.dot
	dot $< -Tpng -o $@

parser/docs/free_form_types.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M free_form_types_main $< -o $@

parser/docs/free_form_types.png: parser/docs/free_form_types.dot
	dot $< -Tpng -o $@

parser/docs/body.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M body $< -o $@

parser/docs/body.png: parser/docs/body.dot
	dot $< -Tpng -o $@

parser/docs/trailer_beg.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M trailer_beg $< -o $@

parser/docs/trailer_beg.png: parser/docs/trailer_beg.dot
	dot $< -Tpng -o $@

parser/docs/trailer_end.dot: parser/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp -M trailer_end $< -o $@

parser/docs/trailer_end.png: parser/docs/trailer_end.dot
	dot $< -Tpng -o $@

parser/machine.go: parser/machine.go.rl common/common.rl

parser/machine.go:
	$(RAGEL) -Z -G2 -e -o $@ $<
	@sed -i '' '/^\/\/line/d' $@
	$(MAKE) file=$@ snake2camel

.PHONY: tests
tests:
	$(GO_TEST) ./...

.PHONY: bench
bench: parser/machine.go parser/perf_test.go
	go test -bench=. -run=Bench -benchmem -benchtime=5s ./parser
