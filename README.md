# go-mirror
A tool for performing reflective actions on objects

Made for devil worshipers, and people who like to thrash hard (thrash the garbage collector that is).

Examples:
```go
package main

import (
	"github.com/narcolepticsnowman/go-mirror/mirror"
)

//get
println(
	mirror.Reflect(test).
        GetPath("/Foo/Bar").
        Value().
        Interface().
        String()
	)
//set
mirror.Reflect(test).SetPath("/Foo/Bar", "NewValue")

//exec
res = mirror.Reflect(test).
	GetPath("/favoriteMethod").exec("arg1", 2, &myStruct{}).
	UnwrapResult()[0].(String)
```
