<p align="right">
⭐ &nbsp;&nbsp;<strong>the project to show your appreciation.</strong> :arrow_upper_right:
</p>

<p align="right">
<a href="http://pkg.go.dev/github.com/romance-dev/speedboost"><img src="https://pkg.go.dev/badge/github.com/romance-dev/speedboost" /></a>
<a href="https://goreportcard.com/report/github.com/romance-dev/speedboost"><img src="https://goreportcard.com/badge/github.com/romance-dev/speedboost" /></a>
</p>

# SPEEDBOOST

Call any function from a shared library that is C ABI compatible (i.e. you can build with C, Zig, Rust etc).

It is CGO free and has less overhead than [Purego](https://github.com/ebitengine/purego).

After factoring in the overhead of calling a shared library, you may find that it is faster than
the equivalent Go implementation: https://niklas-heer.github.io/speed-comparison


See the example project for more details.

### C Library programmed in Zig (cross-platform)

```zig

const std = @import("std");

// multiply returns a ⨯ b, where a and b are float64.
export fn multiply(a: f64, b: f64) callconv(.c) f64 {
    return a * b;
}
```

### Go code consuming Library

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
    defer lib.Unload()

    addPtr := lib.GetSymbol("multiply")

    cifAdd := sb.SetFuncSignature(sb.DoubleTypeDescriptor, sb.DoubleTypeDescriptor, sb.DoubleTypeDescriptor)

    // Call Function
    var result float64
    err = sb.CallFunction(cifAdd, addPtr, Ptr(&result), Ptr(new(40.0))), Ptr(new(2.0)))
    if err != nil {
        panic(err)
    }

    fmt.Printf("multiply(40.0, 2.0) = %f\n", result) // 80.0    
}

```

## How to run example project

Install Zig 0.16 (later versions may break `build.zig`)

```bash 
$ go generate ./...
$ CGO_ENABLED=0 go run .

```

### For cross-platform compilation

```bash 
$ GOOS=linux go generate ./...
$ GOOS=linux CGO_ENABLED=0 go run .
```

## Special Thanks

[@kolkov](https://github.com/kolkov) of https://github.com/go-webgpu/goffi