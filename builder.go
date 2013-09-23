// Copyright 2013 The llgo Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package llgo

import (
	"github.com/axw/gollvm/llvm"
)

type Builder struct {
	llvm.Builder
	types *TypeMap
}

func newBuilder(tm *TypeMap) *Builder {
	return &Builder{
		Builder: llvm.GlobalContext().NewBuilder(),
		types:   tm,
	}
}

func (b *Builder) CreateLoad(v llvm.Value, name string) llvm.Value {
	if v.Type().ElementType() == b.types.ptrstandin {
		// We represent recursive pointer types (T = *T)
		// in LLVM as a pointer to "ptrstdin", where
		// ptrstandin is a unique named struct.
		//
		// Cast the the pointer to a pointer to its own type first.
		v = b.CreateBitCast(v, llvm.PointerType(v.Type(), 0), "")
	}
	return b.Builder.CreateLoad(v, name)
}

func (b *Builder) CreateStore(v, ptr llvm.Value) {
	if !b.types.ptrstandin.IsNil() {
		vtyp, ptrtyp := v.Type(), ptr.Type()
		if vtyp == ptrtyp {
			// We must be dealing with a pointer to a recursive pointer
			// type, so bitcast the pointer to a pointer to its own
			// type first.
			ptr = b.CreateBitCast(ptr, llvm.PointerType(ptrtyp, 0), "")
		}
	}
	b.Builder.CreateStore(v, ptr)
}
