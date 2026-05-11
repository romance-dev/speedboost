<p align="right">
⭐ &nbsp;&nbsp;<strong>the project to show your appreciation.</strong> :arrow_upper_right:
</p>

<p align="right">
<a href="http://pkg.go.dev/github.com/romance-dev/speedboost"><img src="https://pkg.go.dev/badge/github.com/romance-dev/speedboost" /></a>
<a href="https://goreportcard.com/report/github.com/romance-dev/speedboost"><img src="https://goreportcard.com/badge/github.com/romance-dev/speedboost" /></a>
</p>

# SPEEDBOOST

Call any function from a shared library that is C ABI compatible (i.e. you can build with C, Zig, Rust etc).

It is CGO free and less overheard than [Purego](https://github.com/ebitengine/purego).

After factoring in the overhead of calling a shared library, you may find that it is faster than
the equivalent Go implementation.

See the example project for more details.

```zig

const std = @import("std");

// multiply returns a ⨯ b, where a and b are float64.
export fn multiply(a: f64, b: f64) callconv(.c) f64 {
    return a * b;
}
```

```go
type Ptr = unsafe.Pointer

//go:generate zig build --build-file ./lib/build.zig go-build
//go:embed lib/zig-out/mymath.shared
var sharedLibrary []byte

func main() {
lib, err := sb.LoadLibrary(sharedLibrary)
if err != nil {
panic(err)
}

addPtr := lib.GetSymbol("multiply")

cifAdd := sb.SetFuncSignature(sb.DoubleTypeDescriptor, sb.DoubleTypeDescriptor, sb.DoubleTypeDescriptor)

// Call Function
var result float64

err = sb.CallFunction(cifAdd, addPtr, Ptr(&result), Ptr(new(float64(40))), Ptr(new(float64(2))))
if err != nil {
panic(err)
}

fmt.Printf("add(40, 2) = %f\n", result) // 80

lib.Unload()
}

```

## How to run example project

```bash 
$ go generate ./...
$ CGO_ENABLED=0 go run .

```

### For cross-platform compilation

```bash 
$ GOOS=linux go generate ./...
$ GOOS=linux CGO_ENABLED=0 go run .
```