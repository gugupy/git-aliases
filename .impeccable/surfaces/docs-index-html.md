---
version: 1
slug: "docs-index-html"
primary_target: "docs/index.html"
related_targets: []
---

# Surface: download page (docs/index.html)

Mode: Persuade. Visitor: a developer (new to aliases or experienced) deciding to install. Action: download the binary for their system. Proof on hand: the real CLI copy (main menu, -h output, built-in alias list) and docs/demo.gif. Constraints: plain static HTML/CSS/JS served by GitHub Pages, links always point at the latest release, no invented social proof, keyboard operable, follows system light/dark.

## Direction contract

THESIS: the page is a git aliases terminal session; it refuses the hero-plus-sections product page by letting the CLI's own output and prompts carry every section.

OWN-WORLD: terminal ground (dark #13111c / light paper-white profile), one self-hosted monospace (Commit Mono) for everything, the banner's cyan-to-purple ramp, purple #A78BFA accent, #C4B5FD key hints, rounded 1px boxes exactly where the CLI draws them (prompt boxes purple, note boxes gray), `$` prompts and `#` comment headings.

STORY: the visitor sees `$ git aliases` run, recognises the tool, picks their binary in a TUI select box, watches the demo, copies install commands, sees the real menu and alias list, ends at a prompt.

FIRST VIEWPORT: top tab bar with section links; prompt line typing `git aliases`; full-width gradient ASCII banner; one-line lede; a purple prompt box "Download git aliases" holding the platform select list with the detected system under the `>` cursor and key hints at its foot.

FORM: pinned. seed: none (pinned): "Can we keep the index.html in the terminal look", answers "Session + TUI boxes", "Follow the system", "Yes, TUI-style keys". Signature interaction: arrow/j/k cursor through the download list, Enter downloads. Motion: one typed command and banner print on load.

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
