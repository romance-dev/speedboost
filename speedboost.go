package speedboost

import (
	"context"
	"os"
	"unsafe"

	"github.com/go-webgpu/goffi/ffi"
	"github.com/go-webgpu/goffi/types"
)

var (
	VoidTypeDescriptor    = types.VoidTypeDescriptor
	IntTypeDescriptor     = types.IntTypeDescriptor
	FloatTypeDescriptor   = types.FloatTypeDescriptor
	DoubleTypeDescriptor  = types.DoubleTypeDescriptor
	UInt8TypeDescriptor   = types.UInt8TypeDescriptor
	SInt8TypeDescriptor   = types.SInt8TypeDescriptor
	UInt16TypeDescriptor  = types.UInt16TypeDescriptor
	SInt16TypeDescriptor  = types.SInt16TypeDescriptor
	UInt32TypeDescriptor  = types.UInt32TypeDescriptor
	SInt32TypeDescriptor  = types.SInt32TypeDescriptor
	UInt64TypeDescriptor  = types.UInt64TypeDescriptor
	SInt64TypeDescriptor  = types.SInt64TypeDescriptor
	PointerTypeDescriptor = types.PointerTypeDescriptor
)

// SharedLibrary represents a shared library that has been loaded and is ready to use.
type SharedLibrary struct {
	ptr  unsafe.Pointer
	path string
}

func (l *SharedLibrary) handle() unsafe.Pointer {
	return l.ptr
}

// Unload can be called after the already loaded library is no longer needed.
func (l *SharedLibrary) Unload() {
	ffi.FreeLibrary(l.ptr)
	os.Remove(l.path)

	l.ptr = nil
	l.path = ""
}

// GetSymbol returns a pointer to an exported function in the shared library.
func (l *SharedLibrary) GetSymbol(symbol string) unsafe.Pointer {
	fnPtr, err := ffi.GetSymbol(l.ptr, symbol)
	if err != nil {
		panic(err)
	}
	return fnPtr
}

// LoadLibrary loads a shared library after being supplied with the binary for the library.
func LoadLibrary(libBinary []byte) (*SharedLibrary, error) {
	f, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Copy data to file
	_, err = f.Write(libBinary)
	if err != nil {
		return nil, err
	}

	lib, err := ffi.LoadLibrary(f.Name())
	if err != nil {
		return nil, err
	}

	return &SharedLibrary{
		ptr:  lib,
		path: f.Name(),
	}, nil
}

// CallFunctionContext calls fn (obtained using GetSymbol). It accepts a context for the first argument.
func CallFunctionContext(ctx context.Context, cif *types.CallInterface, fn unsafe.Pointer, rvalue unsafe.Pointer, avalue ...unsafe.Pointer) error {
	return ffi.CallFunctionContext(ctx, cif, fn, rvalue, avalue)
}

// CallFunction calls fn (obtained using GetSymbol).
func CallFunction(cif *types.CallInterface, fn unsafe.Pointer, rvalue unsafe.Pointer, avalue ...unsafe.Pointer) error {
	return CallFunctionContext(context.Background(), cif, fn, rvalue, avalue...)
}

// SetFuncSignature configures the C function's signature (i.e. return value argument types).
func SetFuncSignature(returnType *types.TypeDescriptor, argTypes ...*types.TypeDescriptor) *types.CallInterface {
	var cif types.CallInterface
	err := ffi.PrepareCallInterface(&cif, types.DefaultConvention(), returnType, argTypes)
	if err != nil {
		panic(err)
	}
	return &cif
}
