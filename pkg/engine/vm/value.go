package vm

import (
	"fmt"
)

type ValueType int

const (
	ValNil ValueType = iota
	ValBool
	ValNumber
	ValString
	ValObject // Interface{} / Map / Pointer
)

// Value represents any ZenoLang value in the VM.
// Designed to be compact and easily stored in an array (Stack).
type Value struct {
	Type  ValueType
	AsNum float64
	AsPtr interface{} // Used for strings and complex objects
}

func (v Value) String() string {
	switch v.Type {
	case ValNil:
		return "nil"
	case ValBool:
		if v.AsNum > 0 {
			return "true"
		}
		return "false"
	case ValNumber:
		return fmt.Sprintf("%g", v.AsNum)
	case ValString:
		return v.AsPtr.(string)
	case ValObject:
		return fmt.Sprintf("%v", v.AsPtr)
	default:
		return "unknown"
	}
}

// ToNative converts a VM Value back to a Go native type.
func (v Value) ToNative() interface{} {
	switch v.Type {
	case ValNil:
		return nil
	case ValBool:
		return v.AsNum > 0
	case ValNumber:
		return v.AsNum
	case ValString, ValObject:
		return v.AsPtr
	default:
		return nil
	}
}

// Helper constructors
func NewNumber(n float64) Value { return Value{Type: ValNumber, AsNum: n, AsPtr: n} }
func NewBool(b bool) Value {
	if b {
		return Value{Type: ValBool, AsNum: 1, AsPtr: true}
	}
	return Value{Type: ValBool, AsNum: 0, AsPtr: false}
}
func NewString(s string) Value      { return Value{Type: ValString, AsPtr: s} }
func NewNil() Value                 { return Value{Type: ValNil} }
func NewObject(o interface{}) Value { return Value{Type: ValObject, AsPtr: o} }
