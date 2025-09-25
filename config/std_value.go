package config

import (
	"time"

	"github.com/spf13/cast"
)

var _ Value = (*stdValue)(nil)

type stdValue struct {
	val any
}

func (v *stdValue) Raw() any {
	return v.val
}

func (v *stdValue) String() string {
	return cast.ToString(v.val)
}

func (v *stdValue) Bool() bool {
	return cast.ToBool(v.val)
}

func (v *stdValue) Int() int {
	return cast.ToInt(v.val)
}

func (v *stdValue) Int8() int8 {
	return cast.ToInt8(v.val)
}

func (v *stdValue) Int16() int16 {
	return cast.ToInt16(v.val)
}

func (v *stdValue) Int32() int32 {
	return cast.ToInt32(v.val)
}

func (v *stdValue) Int64() int64 {
	return cast.ToInt64(v.val)
}

func (v *stdValue) Uint() uint {
	return cast.ToUint(v.val)
}

func (v *stdValue) Uint8() uint8 {
	return cast.ToUint8(v.val)
}

func (v *stdValue) Uint16() uint16 {
	return cast.ToUint16(v.val)
}

func (v *stdValue) Uint32() uint32 {
	return cast.ToUint32(v.val)
}

func (v *stdValue) Uint64() uint64 {
	return cast.ToUint64(v.val)
}

func (v *stdValue) Float32() float32 {
	return cast.ToFloat32(v.val)
}

func (v *stdValue) Float64() float64 {
	return cast.ToFloat64(v.val)
}

func (v *stdValue) Time() time.Time {
	return cast.ToTime(v.val)
}

func (v *stdValue) Duration() time.Duration {
	return cast.ToDuration(v.val)
}

func (v *stdValue) Slice() []any {
	return cast.ToSlice(v.val)
}

func (v *stdValue) IntSlice() []int {
	return cast.ToIntSlice(v.val)
}

func (v *stdValue) StringSlice() []string {
	return cast.ToStringSlice(v.val)
}

func (v *stdValue) Float64Slice() []float64 {
	return cast.ToFloat64Slice(v.val)
}

func (v *stdValue) Map() map[string]any {
	return cast.ToStringMap(v.val)
}
