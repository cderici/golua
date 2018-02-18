package code

// Type1:  1XXXXabc AAAAAAAA BBBBBBBB CCCCCCCC
//
// Binary ops

// Type2:  0111Fabc AAAAAAAA BBBBBBBB CCCCCCCC
//
// Table lookup / setting

// Type3:  0110FaYY AAAAAAAA NNNNNNNN NNNNNNNN
//
// Setting reg from constant

// Type4a: 0101Fab1 AAAAAAAA BBBBBBBB CCCCCCCC
//
// Unary ops

// Type4b: 0101Fa00 AAAAAAAA BBBBBBBB CCCCCCCC
//
// Setting reg from constant (2)

// Type5:  0100FaYY AAAAAAAA NNNNNNNN NNNNNNNN
//
// Jump / call

// Type6:  00RRFabc AAAAAAAA BBBBBBBB CCCCCCCC
//
// Receiving args / Closure creation

// Opcode is the type of opcodes
type Opcode uint32

const (
	Type1Pfx uint32 = 1 << 31
	Type2Pfx uint32 = 7 << 28
	Type3Pfx uint32 = 6 << 28
	Type4Pfx uint32 = 5 << 28
	Type5Pfx uint32 = 4 << 28
)

func MkType1(op BinOp, rA, rB, rC Reg) Opcode {
	return Opcode(1<<31 | rA.ToA() | rB.ToB() | rC.ToC() | op.ToX())
}

func MkType2(f Flag, rA, rB, rC Reg) Opcode {
	return Opcode(0x7<<28 | rA.ToA() | rB.ToB() | rC.ToC() | f.ToF())
}

func MkType3(f Flag, op UnOpK16, rA Reg, k Lit16) Opcode {
	return Opcode(0x6<<28 | f.ToF() | op.ToY() | rA.ToA() | k.ToN())
}

func MkType4a(f Flag, op UnOp, rA, rB Reg) Opcode {
	return Opcode(0x5<<28 | 1<<24 | f.ToF() | op.ToC() | rA.ToA() | rB.ToB())
}

func MkType4b(f Flag, op UnOpK, rA Reg, k Lit8) Opcode {
	return Opcode(0x5<<28 | f.ToF() | rA.ToA() | k.ToB() | op.ToC())
}

func MkType5(f Flag, op JumpOp, rA Reg, k Lit16) Opcode {
	return Opcode(Type5Pfx | f.ToF() | op.ToY() | rA.ToA() | k.ToN())
}

func MkType6(f Flag, n uint8, rA, rB, rC Reg) Opcode {
	return Opcode(f.ToF() | uint32(n)<<28 | rA.ToA() | rB.ToB() | rC.ToC())
}

func MkType0(rA, rB, rC Reg) Opcode {
	return Opcode(rA.ToA() | rB.ToB() | rC.ToC())
}

func (c Opcode) GetA() Reg {
	return Reg((c >> 18 & 1) | (c >> 16 & 0xff))
}

func (c Opcode) GetB() Reg {
	return Reg((c >> 17 & 1) | (c >> 8 & 0xff))
}

func (c Opcode) GetC() Reg {
	return Reg((c >> 16 & 1) | (c & 0xff))
}

func (c Opcode) GetN() uint16 {
	return uint16(c)
}

func (c Opcode) GetX() BinOp {
	return BinOp((c >> 27) & 0xf)
}

func (c Opcode) GetY() uint8 {
	return uint8((c >> 24) & 3)
}

func (c Opcode) GetF() bool {
	return c&(1<<27) != 0
}

func (c Opcode) GetType() uint8 {
	return uint8(c >> 28)
}

func (c Opcode) HasType1() bool {
	return c&(1<<31) != 0
}

func (c Opcode) HasType2or4() bool {
	return c&(1<<28) != 0
}

func (c Opcode) HasSubtypeFlagSet() bool {
	return c&(1<<29) != 0
}

func (c Opcode) HasType4a() bool {
	return c&(1<<24) != 0
}

type BinOp uint8

const (
	OpAdd BinOp = iota
	OpSub
	OpMul
	OpDiv
	OpFloorDiv
	OpMod
	OpPow
	OpBitAnd
	OpBitOr
	OpBitXor
	OpShiftL
	OpShiftR
	OpEq
	OpLt
	OpLeq
	OpConcat
)

func (op BinOp) ToX() uint32 {
	return uint32(op) << 27
}

type Flag uint8

const (
	On  Flag = 1
	Off Flag = 0
)

func (f Flag) ToF() uint32 {
	return uint32(f) << 27
}

type UnOpK16 uint8

const (
	OpInt16 UnOpK16 = iota
	OpK
	OpClosureK
	OpStr2
)

func (op UnOpK16) ToY() uint32 {
	return uint32(op) << 24
}

type Lit16 uint16

func (l Lit16) ToN() uint32 {
	return uint32(l)
}

type UnOp uint8

const (
	OpNeg UnOp = iota
	OpBitNot
	OpLen
	OpClosure
	OpCont
	OpId
	OpTruth // Turn operand to boolean
	OpCell  // ?
)

func (op UnOp) ToC() uint32 {
	return uint32(op)
}

type UnOpK uint8

const (
	OpNil UnOpK = iota
	OpStr0
	OpTable
	OpStr1
	OpBool
	OpCC
	OpInt   // Extra 64 bits (2 opcodes)
	OpFloat // Extra 64 bits (2 opcodes)
	OpStrN  // Extra [n / 4] opcodes
)

func (op UnOpK) ToC() uint32 {
	return uint32(op)
}

type Lit8 uint8

func (l Lit8) ToB() uint32 {
	return uint32(l) << 8
}

type JumpOp uint8

const (
	OpCall JumpOp = iota
	OpJump
	OpJumpIf
	OpJumpIfForLoopDone // Extra opcode (3 registers needed)
)

func (op JumpOp) ToY() uint32 {
	return uint32(op) << 24
}
