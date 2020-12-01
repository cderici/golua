package ast

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/arnodel/golua/token"
)

// NewNumber returns an ExpNode that represents the numeric literal in the given
// token, or an error if it's not a valid numeric literal (the error should not
// happen, perhaps panic instead?).
func NewNumber(id *token.Token) (ExpNode, error) {
	loc := LocFromToken(id)
	nstring := string(id.Lit)
	if isFloatToken(id) {
		f, err := strconv.ParseFloat(nstring, 64)
		if err != nil {
			return nil, err
		}
		return Float{Location: loc, Val: f}, nil
	}
	var n uint64
	var err error
	if strings.HasPrefix(nstring, "0x") || strings.HasPrefix(nstring, "0X") {
		nstring = nstring[2:]
		if len(nstring) > 16 {
			// A hex integral constant is "truncated" if too long (more than 8 bytes)
			nstring = nstring[len(nstring)-16:]
		}
		n, err = strconv.ParseUint(nstring, 16, 64)
	} else {
		n, err = strconv.ParseUint(nstring, 10, 64)
	}
	if err != nil {
		return nil, err
	}
	return Int{Location: loc, Val: n}, nil
}

// IsNumber returns true if the given expression node is a numerical value.
func IsNumber(e ExpNode) bool {
	n, ok := e.(numberOracle)
	return ok && n.isNumber()
}

//
// Int
//

// Int is an expression node representing a non-negative integer literal.
type Int struct {
	Location
	Val uint64
}

var _ ExpNode = Int{}

// NewInt returns an Int instance with the given value.
func NewInt(val uint64) Int {
	return Int{Val: val}
}

// HWrite prints a tree representation of the node.
func (n Int) HWrite(w HWriter) {
	w.Writef("%d", n.Val)
}

// ProcessExp uses the given ExpProcessor to process the receiver.
func (n Int) ProcessExp(p ExpProcessor) {
	p.ProcessIntExp(n)
}

func (n Int) isNumber() bool {
	return true
}

//
// Float
//

// Float is an expression node representing a floating point numeric literal.
type Float struct {
	Location
	Val float64
}

var _ ExpNode = Float{}

// NewFloat returns a Float instance with the given value.
func NewFloat(x float64) Float {
	return Float{Val: x}
}

// ProcessExp uses the given ExpProcessor to process the receiver.
func (f Float) ProcessExp(p ExpProcessor) {
	p.ProcessFloatExp(f)
}

// HWrite prints a tree representation of the node.
func (f Float) HWrite(w HWriter) {
	w.Writef("%f", f.Val)
}

func (f Float) isNumber() bool {
	return true
}

type numberOracle interface {
	isNumber() bool
}

func isFloatToken(tok *token.Token) bool {
	switch tok.Type {
	case token.NUMDEC:
		return bytes.ContainsAny(tok.Lit, ".eE")
	case token.NUMHEX:
		return bytes.ContainsAny(tok.Lit, ".pP")
	default:
		return false
	}
}
