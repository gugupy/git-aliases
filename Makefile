BIN := git-aliases
VERSION ?= 1.0.0
LDFLAGS := -s -w -X main.version=$(VERSION)
TARGETS := darwin/arm64 darwin/amd64 linux/amd64 linux/arm64 windows/amd64

.PHONY: build release test clean demo debug

# Build for the current machine.
build:
	go build -trimpath -ldflags="$(LDFLAGS)" -o dist/$(BIN) .

# Cross-compile for every target in TARGETS.
release:
	@for t in $(TARGETS); do \
		os=$${t%/*}; arch=$${t#*/}; ext=; [ $$os = windows ] && ext=.exe; \
		echo "building $$os/$$arch"; \
		CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -trimpath -ldflags="$(LDFLAGS)" \
			-o dist/$(BIN)-$$os-$$arch$$ext . || exit 1; \
	done

test:
	go test ./...

# Record docs/demo.gif with VHS (brew install vhs).
demo: build
	vhs docs/demo.tape

# Run the TUI under a headless Delve server; attach from VS Code
# with "Attach to headless dlv (:2345)".
debug:
	dlv debug --headless --listen=127.0.0.1:2345 --api-version=2 --accept-multiclient .

clean:
	rm -rf dist
