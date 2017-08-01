# Gosp -- Go struct parser
Gosp is a simple helper package to make binary file format parser easier, it uses reflection in Golang to parse a struct format in a binary file.

# Installation
```
go get github.com/Escapingbug/gosp
```

# Example
Actually the code in `parser_test.go` is kind of an example.

import look like this
```
import "github.com/Escapingbug/gosp"
```
And we use functions as parser.xxxx

First you need to have a struct, let's say:
```
type MyStruct {
    Num uint64 // beware that field must be exported
}
```

And this struct represents a file format header or something.
It means that this struct is saved in some file.

File should be like this:
```
0xde 0xad 0xbe 0xef 0xde 0xad 0xbe 0xef
```
This 8 bytes, if we use big endian, this should be number 0xdeadbeefdeadbeef, which is 16045690984833335023 decimal.

So to parse this file, we use ParseStructFromBinaryFile function, name is a little bit long to specify its usage.

Like this:
```
parsed, ok := parser.Parse(file, reflect.TypeOf(MyStruct{}), parser.BIG_ENDIAN)
```

But we haven't type asserted yet, we should do that first
```
parsedVal _ := parsed.(MyStruct)
```

And now `parsedVal` is our struct, with value like this:
```
MyStruct:
    Num: 16045690984833335023
```
