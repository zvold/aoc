package instr

import (
	"fmt"
	"log"
)

// Immutable state representing the machine.
type State struct {
	a, b, c int    // Registers.
	pc      int    // Program counter.
	out     string // Program output.
}

func (s State) String() string {
	return fmt.Sprintf("{A=%032b B=%032b C=%032b, pc=%d}", s.a, s.b, s.c, s.pc)
}

func (s State) Print(v int) State {
	if len(s.out) == 0 {
		s.out = fmt.Sprintf("%d", v)
	} else {
		s.out = fmt.Sprintf("%s,%d", s.out, v)
	}
	return s
}

func (s State) Output() string {
	return s.out
}

func (s State) GetReg(cop ComboOp) int {
	switch cop.val() {
	case 4:
		return s.a
	case 5:
		return s.b
	case 6:
		return s.c
	default:
		log.Fatalf("Invalid combo op for registers: %v", cop)
	}
	log.Fatalf("Unreachable")
	return 0
}

func (s State) SetReg(cop ComboOp, v int) State {
	switch cop.val() {
	case 4:
		s.a = v
	case 5:
		s.b = v
	case 6:
		s.c = v
	default:
		log.Fatalf("Invalid combo op for registers: %v", cop)
	}
	return s
}

func (s State) GetPc() int {
	return s.pc
}

func (s State) Jump(v int) State {
	s2 := s
	s2.pc = v
	return s2
}

func (s State) Next() State {
	return s.Jump(s.pc + 1)
}
