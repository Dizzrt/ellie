package config

import "time"

type Value interface {
	Raw() any
	String() string
	Bool() bool
	Int() int
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
	Uint() uint
	Uint8() uint8
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Float32() float32
	Float64() float64
	Time() time.Time
	Duration() time.Duration
	Slice() []any
	IntSlice() []int
	StringSlice() []string
	Float64Slice() []float64
	Map() map[string]any
}
