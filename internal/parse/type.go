package parse

import (
	"blockwatch.cc/tzgo/micheline"
	"github.com/jeanschmitt/tzgen/pkg/types"
)

func (p *Parser) parseType(t *micheline.Typedef) (types.Type, error) {
	// Unwrap optional
	if t.Optional {
		typ, err := p.parseType(&micheline.Typedef{Name: t.Name, Type: t.Type, Args: t.Args})
		if err != nil {
			return nil, err
		}
		return &types.Option{Type: typ}, err
	}

	// Builtin types

	switch t.Type {
	case types.TypeNat:
		return types.Nat{}, nil
	case types.TypeInt:
		return types.Int{}, nil
	case types.TypeString:
		return types.String{}, nil
	case types.TypeBool:
		return types.Bool{}, nil
	case types.TypeBytes:
		return types.Bytes{}, nil
	case types.TypeUnit:
		return types.Unit{}, nil
	case types.TypeTimestamp:
		return types.Timestamp{}, nil
	case types.TypeAddress:
		return types.Address{}, nil
	case types.TypeMutez:
		return types.Mutez{}, nil
	case types.TypeKey:
		return types.Key{}, nil
	case types.TypeKeyHash:
		return types.KeyHash{}, nil
	case types.TypeSignature:
		return types.Signature{}, nil
	case types.TypeChainID:
		return types.ChainID{}, nil
	case types.TypeOperation:
		return types.Operation{}, nil
	case types.TypeContract:
		return types.Contract{}, nil
	}

	// Container types

	if t.Type == types.TypeList || t.Type == types.TypeSet {
		itemType, err := p.parseType(&t.Args[0])
		if err != nil {
			return nil, err
		}
		switch t.Type {
		case types.TypeList:
			return &types.List{Type: itemType}, nil
		case types.TypeSet:
			return &types.Set{Type: itemType}, nil
		}
	}
	if t.Type == types.TypeUnion || t.Type == types.TypeMap || t.Type == types.TypeBigmap {
		type1, err := p.parseType(&t.Args[0])
		if err != nil {
			return nil, err
		}
		type2, err := p.parseType(&t.Args[1])
		if err != nil {
			return nil, err
		}
		switch t.Type {
		case types.TypeUnion:
			return p.registerUnion(&types.Union{Left: type1, Right: type2}), nil
		case types.TypeMap:
			return &types.Map{Key: type1, Value: type2}, nil
		case types.TypeBigmap:
			return &types.BigMap{Key: type1, Value: type2}, nil
		}
	}
	if t.Type == types.TypeStruct {
		var fieldTypes []types.Param
		for _, a := range t.Args {
			typ, err := p.parseType(&a)
			if err != nil {
				return nil, err
			}
			fieldTypes = append(fieldTypes, types.Param{Name: a.Name, Type: typ})
		}
		return p.registerStruct(&types.Struct{Fields: fieldTypes}), nil
	}

	return nil, nil
}
