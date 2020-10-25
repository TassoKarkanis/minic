package symbols

import (
	"github.com/TassoKarkanis/minic/types"
)

type scope struct {
	symbols map[string]types.Type
}

type Table struct {
	scopes []scope
}

func NewTable() *Table {
	t := &Table{}
	t.PushScope()
	return t
}

func (t *Table) AddSymbol(name string, type_ types.Type) {
	t.scopes[len(t.scopes)-1].symbols[name] = type_
}

func (t *Table) PushScope() {
	s := scope{
		symbols: make(map[string]types.Type),
	}
	t.scopes = append(t.scopes, s)
}

func (t *Table) PopScope() {
	t.scopes = t.scopes[0 : len(t.scopes)-1]
}
