# vt100

### VT100 Terminal Package

* Supports colors and attributes.
* Can detect the terminal size.
* Can get key-preses, including arrow keys.
* Has a Canvas struct, for drawing to a buffer than only updating the characters that needs to be updated.

### Simple use

```go
fmt.Println(vt100.BrightColor("hi", "Blue"))
```

### Another example

See `cmd/move` for a more advanced example, where a character can be moved around with the arrow keys.

### General info

* Version: 0.2.0
* Licence: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
