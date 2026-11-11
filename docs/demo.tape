# Records docs/demo.gif with VHS (https://github.com/charmbracelet/vhs).
# Run `make demo`, which builds dist/git-aliases first.
# The recording uses a throwaway global Git config, so your real
# ~/.gitconfig is never read or changed.

Output docs/demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1200
Set Height 820
Set Padding 24
Set Theme "Catppuccin Mocha"
Set TypingSpeed 60ms

Hide
Type `export GIT_CONFIG_GLOBAL="$(mktemp)" PATH="$PWD/dist:$PATH" PS1="$ " && clear`
Enter
Show

Type "git aliases"
Sleep 400ms
Enter
Sleep 2s

# Install a few built-in aliases.
Enter
Sleep 1s
Space
Sleep 300ms
Down
Space
Sleep 300ms
Down
Space
Sleep 300ms
Down 3
Space
Sleep 300ms
Down 2
Space
Sleep 1s
Enter
Sleep 1.5s

# Add a custom alias with a command of our own.
Down
Sleep 400ms
Enter
Sleep 1s
Type "undo"
Sleep 400ms
Enter
Sleep 800ms
Enter
Sleep 800ms
Type "reset --soft HEAD~1"
Sleep 800ms
Enter
Sleep 2s

# View the global aliases.
Down 2
Sleep 400ms
Enter
Sleep 1s
Enter
Sleep 3s
Enter
Sleep 1s

# Quit.
Escape
Sleep 1.5s
