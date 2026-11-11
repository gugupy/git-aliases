package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// version is set at build time with -ldflags "-X main.version=...".
var version = "dev"

var bannerArt = []string{
	`         _  _              _   _`,
	`  __ _  (_)| |_      __ _ | | (_)  __ _  ___   ___  ___`,
	` / _` + "`" + ` | | || __|    / _` + "`" + ` || | | | / _` + "`" + ` |/ __| / _ \/ __|`,
	`| (_| | | || |_    | (_| || | | || (_| |\__ \|  __/\__ \`,
	` \__, | |_| \__|    \__,_||_| |_| \__,_||___/ \___||___/`,
	` |___/`,
}

// bannerColors fade from cyan to purple, one per art line.
var bannerColors = []string{"#22D3EE", "#38BDF8", "#60A5FA", "#818CF8", "#A78BFA", "#C084FC"}

// printBanner shows the startup banner. gitVersion is the output of
// "git --version". Narrow terminals get a one-line title instead.
func printBanner(gitVersion string) {
	subtitle := fmt.Sprintf("Manage your global Git aliases · v%s · %s",
		strings.TrimPrefix(version, "v"), gitVersion)

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width < len(bannerArt[3])+2 {
		fmt.Println(title.Render("Git aliases"))
		fmt.Println(dim.Render(subtitle))
		fmt.Println()
		return
	}

	fmt.Println()
	for i, line := range bannerArt {
		color := lipgloss.Color(bannerColors[i%len(bannerColors)])
		fmt.Println(lipgloss.NewStyle().Foreground(color).Bold(true).Render(line))
	}
	fmt.Println(dim.Render(subtitle))
	fmt.Println()
}
