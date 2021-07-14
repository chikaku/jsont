package jsont

import "strconv"

type JSON interface {
	Type() Type
}

type Type uint8

const (
	TObject Type = 0
	TArray
	TString
	TNumber
	TBool
	TNull
)

const (
	True  = true
	False = false
)

type Object map[string]JSON

type Array []JSON

type String string

type Number string

type Bollean bool

type Null struct{}

func (Object) Type() Type  { return TObject }
func (Array) Type() Type   { return TArray }
func (String) Type() Type  { return TString }
func (Number) Type() Type  { return TNumber }
func (Bollean) Type() Type { return TBool }
func (Null) Type() Type    { return TNull }

func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}
