package symbols

import (
	"github.com/TassoKarkanis/minic/codegen"
	"github.com/TassoKarkanis/minic/types"
)

type symbolData struct {
	typ   types.Type
	value *codegen.Value
}

type scope struct {
	symbols map[string]symbolData
	types   map[string]types.Type
}

type Table struct {
	scopes []scope
}

func NewTable() *Table {
	t := &Table{}
	t.PushScope()
	return t
}

func (t *Table) AddSymbol(name string, typ types.Type, val *codegen.Value) {
	scope := t.scopes[len(t.scopes)-1]
	scope.symbols[name] = symbolData{
		typ:   typ,
		value: val,
	}
}

func (t *Table) FindSymbol(name string) (types.Type, *codegen.Value, bool) {
	for i := len(t.scopes) - 1; i >= 0; i-- {
		scope := t.scopes[i]
		data, ok := scope.symbols[name]
		if ok {
			return data.typ, data.value, true
		}
	}

	return nil, nil, false
}

func (t *Table) PushScope() {
	s := scope{
		symbols: make(map[string]symbolData),
		types:   make(map[string]types.Type),
	}
	t.scopes = append(t.scopes, s)
}

func (t *Table) PopScope() {
	t.scopes = t.scopes[0 : len(t.scopes)-1]
}
