package main_test

import (
	"testing"
	"unsafe"
)

type struc1 struct {
	A string
	B string
}

type struc2 struct {
	C string
	D string
}

func BenchmarkUnsafe(b *testing.B) {
	var astruct1 []struc1 = []struc1{{"a", "b"}, {"a", "b"}}
	for i := 0; i < b.N; i++ {
		var _ []struc2 = *(*[]struc2)(unsafe.Pointer(&astruct1))
	}
}

func BenchmarkSafe(b *testing.B) {
	var astruct1 []struc1 = []struc1{{"a", "b"}, {"a", "b"}}
	for i := 0; i < b.N; i++ {
		res := make([]struc2, len(astruct1))
		for _, j := range astruct1 {
			res = append(res, struc2{j.A, j.B})
		}
	}
}

func BenchmarkInterfaceToUnsafe(b *testing.B) {
	var astruct1 []struc1 = []struc1{{"a", "b"}, {"a", "b"}}
	inter := (interface{})(astruct1)
	for i := 0; i < b.N; i++ {
		var _ []struc2 = *(*[]struc2)(unsafe.Pointer(&inter))
	}
}

func main() {}
