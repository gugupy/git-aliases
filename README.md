# git aliases

```text
         _  _              _   _
  __ _  (_)| |_      __ _ | | (_)  __ _  ___   ___  ___
 / _` | | || __|    / _` || | | | / _` |/ __| / _ \/ __|
| (_| | | || |_    | (_| || | | || (_| |\__ \|  __/\__ \
 \__, | |_| \__|    \__,_||_| |_| \__,_||___/ \___||___/
 |___/
Manage your global Git aliases
```

An interactive Go CLI for managing useful aliases in Git's global
configuration. Once installed, run it as `git aliases`. It builds to a single ~4 MB binary for macOS, Linux, and Windows.

![git aliases installing built-in aliases, adding a custom alias, and listing global aliases](docs/demo.gif)

## Requirements

- Git installed and available on `PATH`
- Go 1.26 or later, only to install from source or for development (Go 1.21+
  downloads the required toolchain automatically). A built binary needs only
  Git.

## Install

### Download a binary

Get it from the [download page](https://gugupy.github.io/git-aliases/), or
download the binary for your platform from the
[Releases page](https://github.com/gugupy/git-aliases/releases), rename it to
`git-aliases` (`git-aliases.exe` on Windows), and put it in a folder on your
`PATH`. For example, on an Apple Silicon Mac:

```bash
curl -L -o git-aliases https://github.com/gugupy/git-aliases/releases/latest/download/git-aliases-darwin-arm64
chmod +x git-aliases
mv git-aliases /usr/local/bin/
git aliases
```

Available builds: `darwin-arm64`, `darwin-amd64`, `linux-amd64`,
`linux-arm64`, and `windows-amd64.exe`. Each release also includes
`checksums.txt` with SHA-256 sums. On macOS, if Gatekeeper blocks the
downloaded binary, run `xattr -d com.apple.quarantine git-aliases`.

### Install from source

Install the `git-aliases` command into Go's bin folder
(`$(go env GOPATH)/bin`, usually `~/go/bin`; make sure it is on your `PATH`):

```bash
go install .
git aliases
```

Or build a standalone binary and copy it to any folder on your `PATH`:

```bash
make build
cp dist/git-aliases /usr/local/bin/
git aliases
```

Git runs any program named `git-<name>` on your `PATH` as `git <name>`, so
`git aliases` and `git-aliases` both work.

## Usage

```bash
git aliases              # start the interactive CLI
git aliases -h           # usage (git aliases --help is taken by Git's man pages)
git aliases --version
```

The prompts let you select built-in aliases, add or update custom
aliases, filter global aliases, and remove selected aliases. In multi-select
lists, press Space to toggle an item, Ctrl+A to select all or none, and
Enter to confirm. Each prompt shows its keys at the bottom of its box.

When adding a custom alias, the alias name, a list of common Git commands
(press `/` to search), and the command itself are on one screen. Moving
through the list fills in the command, which you can then edit before saving.
Enter or Tab moves to the next field, Shift+Tab to the previous one, and Esc
returns to the main menu. Ctrl+C quits. The CLI previews
custom commands and asks before replacing an alias with a different command.
Run it in an interactive terminal.

Aliases are written to your global Git configuration and apply to all
repositories for your user account. The built-in aliases are defined in
`internal/gitconfig/gitconfig.go`.

| Alias | Command |
| --- | --- |
| `git st` | `git status --short --branch` |
| `git co` | `git checkout` |
| `git br` | `git branch` |
| `git cp` | `git cherry-pick` |
| `git ci` | `git commit` |
| `git last` | `git log -1 HEAD` |
| `git lg` | `git log --oneline --graph --decorate --all` |
| `git sp` | `git stash push` |
| `git spp` | `git stash pop` |
| `git sl` | `git stash list` |
| `git theirs` | `git checkout --theirs -- .` |
| `git unstage` | `git reset HEAD --` |

During a merge conflict, `git theirs` chooses the incoming side for all
conflicted paths in the current repository. It does not stage the resolutions;
review them and run `git add` when ready.

## Development

Run straight from the source, without building a binary:

```bash
go run .              # start the interactive CLI
go run . --version    # print the version ("dev" when run this way)
go run . --help
```

`go run` compiles to a temporary location and starts the CLI right away; the
first run downloads the dependencies listed in `go.mod`. Nothing is written
to the project folder, so just run it again after each change.

To try changes without touching your real `~/.gitconfig`, point Git at a
throwaway file:

```bash
GIT_CONFIG_GLOBAL=/tmp/test-gitconfig go run .
```

Other tasks:

```bash
make test                      # run the tests
make build                     # build dist/git-aliases for this machine
make build VERSION=1.2.0       # set the version shown in the banner and --version
make release                   # cross-compile for macOS, Linux, and Windows into dist/
make demo                      # re-record docs/demo.gif (needs VHS: brew install vhs)
make clean                     # remove dist/
```

Without `make`, use `go test ./...` and `go build -o dist/git-aliases .`.

### Releasing

Push a version tag to publish a release:

```bash
git tag v1.2.0
git push origin v1.2.0
```

The [Release workflow](.github/workflows/release.yml) runs the tests,
cross-compiles every target with `make release` (the version comes from the
tag, without the leading `v`), and attaches the binaries and `checksums.txt`
to a GitHub release for that tag, with notes generated from the commits.

The download page at `docs/index.html` is served by GitHub Pages and always
links to the latest release, so it needs no changes when you publish one.

Project layout:

- `main.go`: the interactive prompts and menu flow
- `banner.go`: the startup banner and version
- `internal/gitconfig/`: reading and writing global Git aliases, the built-in
  aliases, and the list of common commands, with tests
