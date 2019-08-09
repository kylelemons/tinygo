package main

import (
	"reflect"
	"unsafe"
)

type (
	myint    int
	myslice  []byte
	myslice2 []myint
	mychan   chan int
	myptr    *int
	point    struct {
		X int16
		Y int16
	}
	mystruct struct {
		n    int `foo:"bar"`
		some point
		zero struct{}
		buf  []byte
		Buf  []byte
	}
)

func main() {
	println("matching types")
	println(reflect.TypeOf(int(3)) == reflect.TypeOf(int(5)))
	println(reflect.TypeOf(int(3)) == reflect.TypeOf(uint(5)))
	println(reflect.TypeOf(myint(3)) == reflect.TypeOf(int(5)))
	println(reflect.TypeOf(myslice{}) == reflect.TypeOf([]byte{}))
	println(reflect.TypeOf(myslice2{}) == reflect.TypeOf([]myint{}))
	println(reflect.TypeOf(myslice2{}) == reflect.TypeOf([]int{}))

	println("\nvalues of interfaces")
	var zeroSlice []byte
	var zeroFunc func()
	var zeroMap map[string]int
	var zeroChan chan int
	n := 42
	for _, v := range []interface{}{
		// basic types
		true,
		false,
		int(2000),
		int(-2000),
		uint(2000),
		int8(-3),
		int8(3),
		uint8(200),
		int16(-300),
		int16(300),
		uint16(50000),
		int32(7 << 20),
		int32(-7 << 20),
		uint32(7 << 20),
		int64(9 << 40),
		int64(-9 << 40),
		uint64(9 << 40),
		uintptr(12345),
		float32(3.14),
		float64(3.14),
		complex64(1.2 + 0.3i),
		complex128(1.3 + 0.4i),
		myint(32),
		"foo",
		unsafe.Pointer(new(int)),
		// channels
		zeroChan,
		mychan(zeroChan),
		// pointers
		new(int),
		new(error),
		&n,
		myptr(new(int)),
		// slices
		[]byte{1, 2, 3},
		make([]uint8, 2, 5),
		[]rune{3, 5},
		[]string{"xyz", "Z"},
		zeroSlice,
		[]byte{},
		[]float32{1, 1.32},
		[]float64{1, 1.64},
		[]complex64{1, 1.64 + 0.3i},
		[]complex128{1, 1.128 + 0.4i},
		myslice{5, 3, 11},
		// array
		[4]int{1, 2, 3, 4},
		// functions
		zeroFunc,
		emptyFunc,
		// maps
		zeroMap,
		map[string]int{},
		// structs
		struct{}{},
		struct{ error }{},
		struct {
			a uint8
			b int16
			c int8
		}{42, 321, 123},
		mystruct{5, point{-5, 3}, struct{}{}, []byte{'G', 'o'}, []byte{'X'}},
	} {
		showValue(reflect.ValueOf(v), "")
	}

	// test sizes
	println("\nsizes:")
	for _, tc := range []struct {
		name string
		rt   reflect.Type
	}{
		{"int8", reflect.TypeOf(int8(0))},
		{"int16", reflect.TypeOf(int16(0))},
		{"int32", reflect.TypeOf(int32(0))},
		{"int64", reflect.TypeOf(int64(0))},
		{"uint8", reflect.TypeOf(uint8(0))},
		{"uint16", reflect.TypeOf(uint16(0))},
		{"uint32", reflect.TypeOf(uint32(0))},
		{"uint64", reflect.TypeOf(uint64(0))},
		{"float32", reflect.TypeOf(float32(0))},
		{"float64", reflect.TypeOf(float64(0))},
		{"complex64", reflect.TypeOf(complex64(0))},
		{"complex128", reflect.TypeOf(complex128(0))},
	} {
		println(tc.name, int(tc.rt.Size()), tc.rt.Bits())
	}
	assertSize(reflect.TypeOf(uintptr(0)).Size() == unsafe.Sizeof(uintptr(0)), "uintptr")
	assertSize(reflect.TypeOf("").Size() == unsafe.Sizeof(""), "string")
	assertSize(reflect.TypeOf(new(int)).Size() == unsafe.Sizeof(new(int)), "*int")

	// SetBool
	rv := reflect.ValueOf(new(bool)).Elem()
	rv.SetBool(true)
	if rv.Bool() != true {
		panic("could not set bool with SetBool()")
	}

	// SetInt
	for _, v := range []interface{}{
		new(int),
		new(int8),
		new(int16),
		new(int32),
		new(int64),
	} {
		rv := reflect.ValueOf(v).Elem()
		rv.SetInt(99)
		if rv.Int() != 99 {
			panic("could not set integer with SetInt()")
		}
	}

	// SetUint
	for _, v := range []interface{}{
		new(uint),
		new(uint8),
		new(uint16),
		new(uint32),
		new(uint64),
		new(uintptr),
	} {
		rv := reflect.ValueOf(v).Elem()
		rv.SetUint(99)
		if rv.Uint() != 99 {
			panic("could not set integer with SetUint()")
		}
	}

	// SetFloat
	for _, v := range []interface{}{
		new(float32),
		new(float64),
	} {
		rv := reflect.ValueOf(v).Elem()
		rv.SetFloat(2.25)
		if rv.Float() != 2.25 {
			panic("could not set float with SetFloat()")
		}
	}

	// SetComplex
	for _, v := range []interface{}{
		new(complex64),
		new(complex128),
	} {
		rv := reflect.ValueOf(v).Elem()
		rv.SetComplex(3 + 2i)
		if rv.Complex() != 3+2i {
			panic("could not set complex with SetComplex()")
		}
	}

	// SetString
	rv = reflect.ValueOf(new(string)).Elem()
	rv.SetString("foo")
	if rv.String() != "foo" {
		panic("could not set string with SetString()")
	}

	// Set int
	rv = reflect.ValueOf(new(int)).Elem()
	rv.SetInt(33)
	rv.Set(reflect.ValueOf(22))
	if rv.Int() != 22 {
		panic("could not set int with Set()")
	}

	// Set uint8
	rv = reflect.ValueOf(new(uint8)).Elem()
	rv.SetUint(33)
	rv.Set(reflect.ValueOf(uint8(22)))
	if rv.Uint() != 22 {
		panic("could not set uint8 with Set()")
	}

	// Set string
	rv = reflect.ValueOf(new(string)).Elem()
	rv.SetString("foo")
	rv.Set(reflect.ValueOf("bar"))
	if rv.String() != "bar" {
		panic("could not set string with Set()")
	}

	// Set complex128
	rv = reflect.ValueOf(new(complex128)).Elem()
	rv.SetComplex(3 + 2i)
	rv.Set(reflect.ValueOf(4 + 8i))
	if rv.Complex() != 4+8i {
		panic("could not set complex128 with Set()")
	}

	// Set to slice
	rv = reflect.ValueOf([]int{3, 5})
	rv.Index(1).SetInt(7)
	if rv.Index(1).Int() != 7 {
		panic("could not set int in slice")
	}
	rv.Index(1).Set(reflect.ValueOf(8))
	if rv.Index(1).Int() != 8 {
		panic("could not set int in slice")
	}
	if rv.Len() != 2 || rv.Index(0).Int() != 3 {
		panic("slice was changed while setting part of it")
	}
}

func emptyFunc() {
}

func showValue(rv reflect.Value, indent string) {
	rt := rv.Type()
	if rt.Kind() != rv.Kind() {
		panic("type kind is different from value kind")
	}
	print(indent+"reflect type: ", rt.Kind().String())
	if rv.CanSet() {
		print(" settable=", rv.CanSet())
	}
	println()
	switch rt.Kind() {
	case reflect.Bool:
		println(indent+"  bool:", rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		println(indent+"  int:", rv.Int())
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		println(indent+"  uint:", rv.Uint())
	case reflect.Float32, reflect.Float64:
		println(indent+"  float:", rv.Float())
	case reflect.Complex64, reflect.Complex128:
		println(indent+"  complex:", rv.Complex())
	case reflect.String:
		println(indent+"  string:", rv.String(), rv.Len())
		for i := 0; i < rv.Len(); i++ {
			showValue(rv.Index(i), indent+"  ")
		}
	case reflect.UnsafePointer:
		println(indent+"  pointer:", rv.Pointer() != 0)
	case reflect.Array:
		println(indent + "  array")
	case reflect.Chan:
		println(indent+"  chan:", rt.Elem().Kind().String())
		println(indent+"  nil:", rv.IsNil())
	case reflect.Func:
		println(indent + "  func")
		println(indent+"  nil:", rv.IsNil())
	case reflect.Interface:
		println(indent + "  interface")
		println(indent+"  nil:", rv.IsNil())
	case reflect.Map:
		println(indent + "  map")
		println(indent+"  nil:", rv.IsNil())
	case reflect.Ptr:
		println(indent+"  pointer:", rv.Pointer() != 0, rt.Elem().Kind().String())
		println(indent+"  nil:", rv.IsNil())
		if !rv.IsNil() {
			showValue(rv.Elem(), indent+"  ")
		}
	case reflect.Slice:
		println(indent+"  slice:", rt.Elem().Kind().String(), rv.Len(), rv.Cap())
		println(indent+"  pointer:", rv.Pointer() != 0)
		println(indent+"  nil:", rv.IsNil())
		for i := 0; i < rv.Len(); i++ {
			println(indent+"  indexing:", i)
			showValue(rv.Index(i), indent+"  ")
		}
	case reflect.Struct:
		println(indent+"  struct:", rt.NumField())
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			println(indent+"  field:", i, field.Name)
			println(indent+"  tag:", field.Tag)
			println(indent+"  embedded:", field.Anonymous)
			showValue(rv.Field(i), indent+"  ")
		}
	default:
		println(indent + "  unknown type kind!")
	}
}

func assertSize(ok bool, typ string) {
	if !ok {
		panic("size mismatch for type " + typ)
	}
}
