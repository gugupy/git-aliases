# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

The product itself is a terminal UI (Go, Bubble Tea / huh / lipgloss). Its one web surface is the static download page at `docs/index.html`, served by GitHub Pages.

## Users

Two audiences, weighted equally:

- **Developers new to aliases.** Comfortable with everyday Git but have never set up aliases or edited `~/.gitconfig`. They want useful defaults installed without reading documentation.
- **Experienced Git users.** Already rely on aliases and want a faster way to add, review, filter, and prune them than hand-editing config.

Copy and UI must work for the first group without slowing down the second: explain plainly, never assume alias knowledge, but keep every flow quick and keyboard-driven.

## Product Purpose

`git aliases` is an interactive CLI that manages aliases in Git's global configuration. Users install a curated set of built-in aliases, add or update custom ones, filter what is configured, and remove what they don't use, all without editing `~/.gitconfig` by hand. Success means a user gets working aliases in one short session and trusts the tool not to clobber their existing config.

## Positioning

Three things together, which a dotfiles alias list or manual config editing does not offer:

1. **A curated built-in set**: a small, opinionated list of useful aliases (`st`, `co`, `lg`, `theirs`, `unstage`, and others) ready to install in one step.
2. **Interactive, no hand-editing**: guided prompts that preview custom commands and confirm before replacing an alias that has a different command.
3. **A native `git aliases` command**: through Git's `git-<name>` PATH convention it runs as a Git subcommand, so it feels like part of Git.

## Operating Context

- Runs in an interactive terminal on macOS, Linux, or Windows. Light and dark terminal themes and narrow windows are all real conditions; the banner already falls back to a one-line title on narrow terminals.
- Writes to the user's global Git config, which affects every repository for that account.
- Distributed as GitHub release binaries (with `checksums.txt`), through the download page, or built from source with Go.
- Developers test against a throwaway config with `GIT_CONFIG_GLOBAL=/tmp/test-gitconfig go run .`.

## Capabilities and Constraints

- Main menu flows: install built-in aliases, add or update a custom alias (name, a searchable list of common Git commands, and an editable command on one screen), filter global aliases, and remove selected aliases.
- Keys: Space toggles, Ctrl+A selects all or none, Enter confirms, Enter/Tab and Shift+Tab move between fields, `/` searches, Esc returns to the main menu, Ctrl+C quits. Each prompt shows its keys at the bottom of its box.
- `git aliases -h` shows usage (`--help` is taken by Git's man pages); `--version` prints the version.
- Single ~4 MB binary per platform (`darwin-arm64`, `darwin-amd64`, `linux-amd64`, `linux-arm64`, `windows-amd64.exe`); a built binary needs only Git.
- Built-in aliases are defined in `internal/gitconfig/gitconfig.go`.
- **The download page stays plain static HTML/CSS**: no framework and no build step. GitHub Pages serves `docs/` as-is, and the page always links to the latest release so it needs no edit per release.

## Brand Commitments

- Name: "git aliases" (lowercase in titles and the banner), run as `git aliases`.
- Tagline in use: "Manage your global Git aliases".
- Existing identity, recorded but not made binding: an ASCII-art banner with a cyan-to-purple gradient (`#22D3EE` to `#C084FC`) and a `#A78BFA` accent, shared by the CLI and the download page.

## Evidence on Hand

- Demo recording: `docs/demo.gif`, generated from `docs/demo.tape` with VHS (`make demo`).
- The built-in alias table and usage docs in `README.md`.
- No testimonials, star counts, user numbers, press, or adoption figures exist. Future work must not invent any.

## Product Principles

1. **Never surprise the user's config.** Preview before writing and confirm before replacing; global config changes affect every repository.
2. **Feel like Git.** Behave as a native subcommand: terse, predictable, keyboard-first.
3. **Curate, don't catalogue.** A small set of aliases worth installing beats an exhaustive list.
4. **Serve both audiences in one flow.** Clear enough for a first-time alias user, fast enough for someone who already has dozens.
5. **Show the real thing.** Credibility comes from the demo, the command table, and the binaries, not from claims.

## Accessibility & Inclusion

- The CLI must stay fully keyboard-operable, with visible key hints on every prompt.
- Output must stay legible in both light and dark terminal themes and degrade gracefully at narrow widths.
- The download page follows the same expectations: keyboard reachable, readable in light and dark color schemes.
