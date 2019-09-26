## Simple text editor

A small and experimental editor. Don't use it on files you care about!

## Features and limitations

* The editor must be given a filename at launch-time.
* The editor is always in "insert mode".
* All trailing spaces are removed when saving, but a final newline is kept.
* Understands tab characters.
* The arrow keys works correctly.
* Can scroll, with `ctrl-n` and `ctrl-p`.
* `esc` can be used to toggle "end of line mode" where the cursor is limited to the end of lines and "free floating mode".
* The "free floating mode" is enabled at launch, and is not too bad for drawing ASCII graphics.
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
