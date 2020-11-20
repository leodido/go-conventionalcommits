SHELL := /bin/bash
RAGEL := ragel -I common

.PHONY: build
build: slim/machine.go
	@gofmt -w -s ./slim

.PHONY: clean
clean:
	@$(RM) slim/machine.go
	@$(RM) -R slim/docs

.PHONY: dots
dots:
	@mkdir -p slim/docs
	$(MAKE) -s slim/docs/main.dot

.PHONY: docs
docs: dots slim/docs/main.png

.PHONY: snake2camel
snake2camel:
	@awk -i inplace '{ \
	while ( match($$0, /(.*)([a-z]+[0-9]*)_([a-zA-Z0-9])(.*)/, cap) ) \
	$$0 = cap[1] cap[2] toupper(cap[3]) cap[4]; \
	print \
	}' $(file)

slim/docs/main.dot: slim/machine.go.rl common/common.rl
	$(RAGEL) -Z -Vp $< -o $@

slim/docs/main.png: slim/docs/main.dot
	dot $< -Tpng -o $@

slim/machine.go: slim/machine.go.rl common/common.rl

slim/machine.go:
	$(RAGEL) -Z -G2 -e -o $@ $<
	@sed -i '/^\/\/line/d' $@
	$(MAKE) file=$@ snake2camel