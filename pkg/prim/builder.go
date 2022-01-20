package prim

import (
	"math/big"
	"reflect"
)

type builder struct {
	args []string
}

func Builder() *builder {
	return &builder{}
}

func (b *builder) Int(i *big.Int) *builder {
	return b.addArg(intV(i))
}

func (b *builder) String(s string) *builder {
	return b.addArg(stringV(s))
}

func (b *builder) Bytes(v []byte) *builder {
	return b.addArg(bytesV(v))
}

// Address can be encoded as string or bytes; we use string for readability.
func (b *builder) Address(a string) *builder {
	return b.String(a)
}

func (b *builder) Primer(p Primer) *builder {
	return b.addArg(p.ToPrim())
}

func (b *builder) List(items interface{}) *builder {
	switch v := items.(type) {
	case []*big.Int:
		return b.ListI(v...)
	case []string:
		return b.ListS(v...)
	case [][]byte:
		return b.ListB(v...)
	}

	// I dream of generics

	itemsVal := reflect.ValueOf(items)
	if itemsVal.Kind() != reflect.Slice {
		panic("items must be a slice")
	}

	itemsLen := itemsVal.Len()
	strs := make([]string, itemsLen)
	for i := 0; i < itemsLen; i++ {
		primer, ok := itemsVal.Index(i).Interface().(Primer)
		if !ok {
			panic("items must be a slice of *big.Int, string, []byte or Primer")
		}
		strs[i] = primer.ToPrim()
	}

	return b.addArg(listV(strs...))
}

func (b *builder) ListI(items ...*big.Int) *builder {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = intV(item)
	}
	return b.addArg(listV(strs...))
}

func (b *builder) ListS(items ...string) *builder {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = stringV(item)
	}
	return b.addArg(listV(strs...))
}

func (b *builder) ListB(items ...[]byte) *builder {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = bytesV(item)
	}
	return b.addArg(listV(strs...))
}

func (b *builder) ListP(items ...Primer) *builder {
	strs := make([]string, len(items))
	for i, item := range items {
		strs[i] = item.ToPrim()
	}
	return b.addArg(listV(strs...))
}

func (b *builder) Union(val interface{}, branch UnionBranch) *builder {
	var inner string

	switch v := val.(type) {
	case *big.Int:
		inner = intV(v)
	case string:
		inner = stringV(v)
	case []byte:
		inner = bytesV(v)
	case Primer:
		inner = v.ToPrim()
	default:
		panic("union inner type should be either *big.Int, string, []byte or Primer")
	}

	return b.addArg(unionV(inner, branch))
}

func (b *builder) Some(val interface{}) *builder {
	var inner string

	switch v := val.(type) {
	case *big.Int:
		inner = intV(v)
	case string:
		inner = stringV(v)
	case []byte:
		inner = bytesV(v)
	case Primer:
		inner = v.ToPrim()
	default:
		panic("option inner type should be either *big.Int, string, []byte or Primer")
	}

	return b.addArg(someV(inner))
}

func (b *builder) None() *builder {
	return b.addArg(noneV())
}

func (b *builder) Finish() string {
	if len(b.args) == 0 {
		return unitV()
	}
	return recPairs(b.args)
}

func (b *builder) addArg(arg string) *builder {
	b.args = append(b.args, arg)
	return b
}

func recPairs(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	return pairV(args[0], recPairs(args[1:]))
}