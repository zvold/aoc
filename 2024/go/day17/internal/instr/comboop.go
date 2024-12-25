package instr

import "log"

var (
	RegA ComboOp = NewComboOp(4)
	RegB ComboOp = NewComboOp(5)
	RegC ComboOp = NewComboOp(6)
)

type ComboOp struct {
	v int
}

func NewComboOp(v int) ComboOp {
	if v < 0 || v >= 7 {
		log.Fatalf("Unexpected combo op value %v", v)
	}
	return ComboOp{v: v}
}

func (cop ComboOp) val() int {
	return cop.v
}

// Dereferences combo op:
// for values < 4, returns the value
// otherwise, reads from corresponding reg
func (cop ComboOp) Deref(s State) int {
	if cop.v < 4 {
		return cop.v
	} else {
		return s.GetReg(cop)
	}
}
