# Red

A simple and limited `vt100` text editor.

Don't use it on files you care about, yet!

## Features and limitations

* Syntax highlighting for Go code.
* The editor must be given a filename when starting.
* The editor is always in "insert mode". Characters are never moved around.
* All trailing spaces are removed when saving, but a final newline is kept.
* Can handle text that contains the tab character (`\t`).
* `Esc` can be used to toggle "writing mode" where the cursor is limited to the end of lines and "ASCII drawing mode".
* There is no undo!
* Lines longer than the terminal width are not handled correctly.

## Hotkeys

* `ctrl-q` to quit
* `ctrl-a` go to start of line
* `ctrl-e` go to end of line
* `ctrl-p` scroll up
* `ctrl-n` scroll down
* `ctrl-l` to redraw the screen
* `ctrl-k` to delete characters to the end of the line
* `ctrl-s` to save (don't use this on files you care about!)
* `ctrl-g` to show cursor positions and word count
* `esc` to toggle "text edit mode" and "ASCII graphics mode"
