package main

import (
	_ "embed"
	"fmt"
	"unsafe"

	sb "github.com/romance-dev/speedboost"
)

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
