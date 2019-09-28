# rED

A simple and limited `vt100` text editor.

For a more modern editor, also written in Go, look into [micro](https://github.com/zyedidia/micro).

## Features and limitations

* Syntax highlighting for Go code!
* Can be used for drawing "ASCII graphics".
* The editor must be given a filename at start.
* The editor is always in "overwrite mode". Characters are never moved around. Not even with `ctrl-k`.
* All trailing spaces are removed when saving, but a final newline is kept.
* `Esc` can be used to toggle "writing mode" where the cursor is limited to the end of lines and "ASCII drawing mode".
* Can handle text that contains the tab character (`\t`).
* Keys like `Home` and `End` are not even registered by the key handler.
* There is no undo.

## Known bugs

* After scrolling, the data editor cursor is misaligned. Don't save after scrolling! Press `ctrl-g` to check if the screen and data cursor coordinates look correct.
* Letters that are not a-z, A-Z or simple punctuation may not be possible to type in.
* Lines longer than the terminal width are not handled correctly.
* Characters may appear on the screen when keys are pressed. Clear them with `ctrl-l`.
* Unicode characters are not displayed correctly when loading a file.

## Hotkeys

* `ctrl-q` to quit
* `ctrl-a` go to start of line
* `ctrl-e` go to end of line
* `ctrl-p` scroll up
* `ctrl-n` scroll down
* `ctrl-l` to redraw the screen
* `ctrl-k` to delete characters to the end of the line
* `ctrl-s` to save (don't use this on files you care about!)
* `ctrl-g` to show cursor positions, current letter and word count
* `esc` to toggle "text edit mode" and "ASCII graphics mode"
