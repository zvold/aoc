package instr

import (
	"fmt"
	"log"
)

type Opcode int

const (
	Adv Opcode = iota
	Bxl
	Bst
	Jnz
	Bxc
	Out
	Bdv
	Cdv
)

var opcodeStrings []string = []string{"adv", "bxl", "bst", "jnz", "bxc", "out", "bdv", "cdv"}

type Instr struct {
	opcode  Opcode
	operand int
}

func (i Instr) String() string {
	// TODO(zvold): operand should be something like Op which e.g. prints "A" for "4"...
	return fmt.Sprintf("%s(%d)", opcodeStrings[i.opcode], i.operand)
}

func NewInstr(i, j int) Instr {
	if i < 0 || i > 7 {
		log.Fatalf("Invalid opcode %d", i)
	}
	if j < 0 || j > 7 {
		log.Fatalf("Invalid operand %d", j)
	}
	return Instr{opcode: Opcode(i), operand: j}
}

func (i Instr) Opcode() Opcode {
	return i.opcode
}

func (i Instr) Operand() int {
	return i.operand
}

// Executes the instruction and returns the mutated State.
func (i Instr) Run(s State) State {
	// Below, 'x' is for combo op, 'y' is for literal op
	switch i.opcode {
	case Adv: // regA <- regA / 2 ^ x
		{
			nom := s.GetReg(RegA) // regA
			den := NewComboOp(i.operand).Deref(s)
			s = s.SetReg(RegA, nom/(1<<den))
		}
	case Bxl: // regB <- regB xor y
		{
			s = s.SetReg(RegB, s.GetReg(RegB)^i.operand)
		}
	case Bst: // regB <- x mod 8
		{
			x := NewComboOp(i.operand).Deref(s)
			s = s.SetReg(RegB, x&0x7) // Keep lowest 3 bits.
		}
	case Jnz: // if regA != 0, jump to y
		{
			if s.GetReg(RegA) != 0 {
				s = s.Jump(i.operand)
				return s // Don't advance PC after jump
			}
		}
	case Bxc: // regB <- regB ^ regC
		s = s.SetReg(RegB, s.GetReg(RegB)^s.GetReg(RegC))
	case Out: // print x mod 8
		s = s.Print(NewComboOp(i.operand).Deref(s) & 0x7)
	case Bdv: // regB <- regA / 2 ^ x
		{
			nom := s.GetReg(RegA) // regA
			den := NewComboOp(i.operand).Deref(s)
			s = s.SetReg(RegB, nom/(1<<den))
		}

	case Cdv: // regC <- regA / 2 ^ x
		{
			nom := s.GetReg(RegA) // regA
			den := NewComboOp(i.operand).Deref(s)
			s = s.SetReg(RegC, nom/(1<<den))
		}

	default:
		log.Fatalf("Invalid instruction %s", i)
	}

	s = s.Next()
	return s
}
