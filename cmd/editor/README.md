# ved

A simple and limited `vt100` text editor.

Don't use it on files you care about, yet!

## Features and limitations

* The editor must be given a filename when starting.
* The editor is always in "insert mode". Characters are never moved around.
* All trailing spaces are removed when saving, but a final newline is kept.
* Can handle text that contains the tab character (`\t`).
* `Esc` can be used to toggle "end of line mode" where the cursor is limited to the end of lines and "ASCII drawing mode".
* There is no undo!

## Hotkeys

* `ctrl-q` to quit
* `ctrl-a` go to start of line
* `ctrl-e` go to end of line
* `ctrl-p` scroll up
* `ctrl-n` scroll down
* `ctrl-l` to redraw
* `ctrl-k` to delete characters to the end of the line
* `ctrl-s` to save (don't use this on files you care about!)
* `esc` to toggle "text edit mode" and "ASCII graphics mode"
* `ctrl-g` to show cursor positions and word count
