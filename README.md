# vt100

VT100 Terminal Package

* Everything is generated from the spec.
* Uses memoization of generated terminal codes.
* No external dependencies.
* Has a Canvas struct that can be used for drawing (see `cmd/canvas` for an example).

Simple use:

```go
fmt.Println(vt100.BrightColor("hi", "Blue"))
```

General info:

* Version: 0.2.0
* Licence: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
