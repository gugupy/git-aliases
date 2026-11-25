---
name: git aliases
description: Manage your global Git aliases. The web presence is a terminal session of the CLI itself.
colors:
  terminal-ground: "#13111c"
  status-bar: "#1c1928"
  command-well: "#1a1725"
  terminal-ink: "#e9e7f2"
  dim-ink: "#948fab"
  box-rule: "#3a3550"
  prompt-violet: "#a78bfa"
  key-hint-lilac: "#c4b5fd"
  cursor-row: "#262040"
  success-green: "#4ade80"
  banner-cyan: "#22d3ee"
  banner-sky: "#38bdf8"
  banner-blue: "#60a5fa"
  banner-indigo: "#818cf8"
  banner-violet: "#a78bfa"
  banner-purple: "#c084fc"
  light-ground: "#f6f5fa"
  light-status-bar: "#e9e7f2"
  light-command-well: "#edebf4"
  light-ink: "#1e1b2e"
  light-dim-ink: "#5d5873"
  light-box-rule: "#cdc8dd"
  light-prompt-violet: "#6d28d9"
  light-key-hint: "#5b21b6"
  light-cursor-row: "#e6defc"
  light-success-green: "#15803d"
typography:
  banner:
    fontFamily: "Commit Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace"
    fontSize: "min(17px, calc((100vw - 4ch) / 34.5))"
    fontWeight: 700
    lineHeight: 1.2
  lede:
    fontFamily: "Commit Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace"
    fontSize: "1.07rem"
    fontWeight: 400
    lineHeight: 1.6
  title:
    fontFamily: "Commit Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace"
    fontSize: "1rem"
    fontWeight: 700
    lineHeight: 1.65
  body:
    fontFamily: "Commit Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.65
  label:
    fontFamily: "Commit Mono, ui-monospace, SF Mono, Menlo, Consolas, monospace"
    fontSize: "13px"
    fontWeight: 400
    lineHeight: 1.65
rounded:
  row: "4px"
  well: "6px"
  box: "10px"
spacing:
  row: "0.3rem"
  paragraph: "0.9rem"
  box-y: "1.25rem"
  section: "4.5rem"
components:
  prompt-box:
    textColor: "{colors.terminal-ink}"
    rounded: "{rounded.box}"
    padding: "1.25rem 2ch 1rem"
  note-box:
    textColor: "{colors.terminal-ink}"
    rounded: "{rounded.box}"
    padding: "1.25rem 2ch 1rem"
  select-row:
    textColor: "{colors.terminal-ink}"
    rounded: "{rounded.row}"
    padding: "0.3rem 1ch"
  select-row-cursor:
    backgroundColor: "{colors.cursor-row}"
    textColor: "{colors.prompt-violet}"
    rounded: "{rounded.row}"
    padding: "0.3rem 1ch"
  command-block:
    backgroundColor: "{colors.command-well}"
    textColor: "{colors.terminal-ink}"
    rounded: "{rounded.well}"
    padding: "0.85rem 9ch 0.85rem 2ch"
  copy-button:
    backgroundColor: "{colors.command-well}"
    textColor: "{colors.dim-ink}"
    rounded: "{rounded.row}"
    padding: "0.1rem 1ch"
  status-window:
    textColor: "{colors.dim-ink}"
    rounded: "3px"
    padding: "0.15rem 0.5ch"
  status-window-current:
    backgroundColor: "{colors.cursor-row}"
    textColor: "{colors.terminal-ink}"
    rounded: "3px"
    padding: "0.15rem 0.5ch"
---

# Design System: git aliases

## Overview

**Creative North Star: "The Session Is the Page"**

Everything on the web reads as one `git aliases` terminal session. The page opens with the command being typed, the CLI's gradient banner prints, and every section after it is either a command and its output or a box the real TUI would draw. Nothing is a picture of a terminal: there is no window chrome, no traffic-light dots and no fake title bar. The page ground *is* the terminal ground, and the content is set as text a visitor can select and copy.

The density is a terminal's: one monospace face, one body size, hierarchy carried by weight and colour rather than scale, and generous vertical gaps between sections, the way blank lines separate commands. Colour comes straight from the CLI's lipgloss styles: violet for prompts and focus, lilac for key hints, gray for rules and notes. The page follows the visitor's light or dark preference, with the light profile drawn like a light terminal theme.

**Key Characteristics:**
- One self-hosted monospace (Commit Mono) for every role, banner included.
- `$` prompts and `#` comment headings; the CLI's own copy (menu items, key hints) shown verbatim.
- Rounded 1px boxes only where the CLI draws them: violet around prompts, gray around notes.
- A `>` cursor row with a tinted band marks the current choice in every list.
- Keyboard parity with the TUI: arrows or j/k move, Enter acts, number keys jump between status-bar windows.

## Colors

The CLI's palette transplanted: a near-black violet ground, one violet voice and a six-step cyan-to-purple ramp that appears only in the banner.

### Primary
- **Prompt Violet** (`prompt-violet`; light: `light-prompt-violet`): the `$` prompt, section headings, prompt-box borders, the cursor row's text, links and the focus ring. Matches the CLI's `accent`.
- **Key-Hint Lilac** (`key-hint-lilac`; light: `light-key-hint`): key names in hint rows and the keys list, inline code. Matches the CLI's `hintKey`.

### Secondary
- **Banner Ramp** (`banner-cyan` through `banner-purple`): one colour per ASCII-art line, top to bottom, exactly as `banner.go` prints them. The light profile uses darker shades of the same hues (from `#0e7490` to `#7e22ce`) for legibility.

### Tertiary
- **Success Green** (`success-green`; light: `light-success-green`): confirmation only, such as the `◆ copied` state, mirroring the CLI's `okSym`.

### Neutral
- **Terminal Ground** (`terminal-ground`; light: `light-ground`): the page itself.
- **Status Bar** (`status-bar`): the tmux-style bar's band.
- **Command Well** (`command-well`): the faint field behind command blocks.
- **Cursor Row** (`cursor-row`): the band behind the selected list row and the current status window.
- **Terminal Ink** (`terminal-ink`) / **Dim Ink** (`dim-ink`): text and secondary text, at 4.5:1 or better on the ground in both profiles.
- **Box Rule** (`box-rule`): note-box borders, hint dividers, the status-bar edge. Matches the CLI's ANSI 8 gray.

### Named Rules
**The Banner-Only Ramp Rule.** The cyan-to-purple ramp lives in the banner and nowhere else. No gradient fills, gradient text or gradient buttons.

**The CLI Source Rule.** Every accent traces back to a lipgloss style in `main.go` or `banner.go`. A new colour needs a matching CLI use first.

## Typography

**Body Font:** Commit Mono (self-hosted from `docs/fonts/`, falling back to ui-monospace, SF Mono, Menlo, Consolas)

**Character:** a neutral, even monospace that holds ASCII art and column alignment, set without ligatures so `--` and `->` read as typed.

### Hierarchy
- **Banner** (700, fluid up to 17px, 1.2): the ASCII-art wordmark, sized so its 56 columns fit the viewport.
- **Lede** (400, 1.07rem, 1.6): the one-paragraph introduction, max 62ch.
- **Title** (700, 1rem): `#` section headings and box titles, set apart by weight and colour, not size.
- **Body** (400, 15px, 1.65; 14px under 640px): all reading text, max 72ch.
- **Label** (400, 13px): key hints, the status bar, file names on mobile, captions.

### Named Rules
**The Weight-Not-Size Rule.** Hierarchy comes from bold, violet and dim, as in a terminal. Only the banner and the lede step up in size.

## Layout

A single centred column of 92ch with 2ch side padding. Text inside it stays left-aligned, as in a terminal. Reading blocks cap at 72ch. Sections are separated by 4.5rem (3.5rem on mobile) with 1rem below each heading. Lists use 0.3rem row padding. The status bar is sticky, and anchors offset by 3.5rem to clear it. Under 640px, two-column rows (label and file, menu item and description) stack, and command blocks wrap rather than scroll, with the copy button moving below the commands.

## Elevation & Depth

Flat. No shadows anywhere. Depth comes only from tonal fields (the command well, the cursor-row band, the status bar) and 1px rules, as a terminal draws it.

### Named Rules
**The No-Shadow Rule.** Nothing floats. A layer is a tonal field or a 1px rule.

## Shapes

Softly rounded rectangles in three sizes: boxes at 10px (the CLI's RoundedBorder), command wells at 6px, and rows, buttons and status windows at 3–4px. Borders are always 1px. Every element is a rectangle, with no pills, circles or angled cuts.

## Components

### Prompt Box
The CLI's `promptBox`: a 1px violet rounded border, a bold violet title, a dim description, the content, then a hint row under a gray rule listing keys in lilac bold with dim actions, separated by ` · `. Used for anything the visitor acts on (the download picker, the main menu).

### Note Box
The CLI's `note`: the same shape with a gray border and a bold ink title. Used for reference content such as the keys list.

### Select List
Rows of `label | file | tag` in a fixed-column grid. The current row shows a violet `>`, bold violet label and ink file name on the cursor-row band. Arrow keys or j/k move the cursor (roving tabindex), Home/End jump, Enter activates, and hover moves the cursor too. A dim `detected` tag marks the visitor's platform.

### Menu Tabs
The CLI main menu as a vertical tab list with the same cursor treatment. Selecting an item swaps the output panel below it.

### Command Block
`$`-prefixed lines on the command well. The prompt is not selectable, and comments are dim. A `copy` button (13px, rule-bordered) copies the commands without prompts or comments and confirms with `◆ copied` in green.

### Navigation
A tmux-style sticky status bar: a violet session chip `git-aliases`, then numbered windows `1:download` and so on. The window in view gets the cursor-row band and a violet `*`. Number keys jump to windows. A GitHub link sits at the right.

### Inputs / Fields
None yet. Focus everywhere is a 2px violet outline with a 2px offset.

## Do's and Don'ts

### Do:
- **Do** write sections as commands and their output, or as boxes the CLI actually draws.
- **Do** use the CLI's real copy (menu labels, key names, flag syntax) verbatim.
- **Do** give every list a `>` cursor and TUI keys (arrows, j/k, Enter), plus normal pointer behaviour.
- **Do** keep both colour profiles: dark ground `#13111c`, light ground `#f6f5fa`.
- **Do** theme browser surfaces: violet selection, violet focus ring, a rule-coloured scrollbar.

### Don't:
- **Don't** draw a terminal window frame, traffic-light dots or a fake title bar. The page is the terminal.
- **Don't** use the banner ramp outside the banner.
- **Don't** add shadows or a second typeface.
- **Don't** show output Git or the CLI would not produce without labelling it, as with the aligned `--get-regexp` columns.
