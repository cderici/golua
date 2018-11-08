package ast

import (
	"github.com/arnodel/golua/ir"
	"github.com/arnodel/golua/token"
)

type BreakStat struct {
	Location
}

func NewBreakStat(tok *token.Token) BreakStat {
	return BreakStat{Location: LocFromToken(tok)}
}

func (s BreakStat) HWrite(w HWriter) {
	w.Writef("break")
}

func (s BreakStat) CompileStat(c *ir.Compiler) {
	EmitJump(c, s, breakLblName)
}

var breakLblName = ir.Name("<break>")
