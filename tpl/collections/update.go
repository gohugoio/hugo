package collections

import (
	"fmt"
	"github.com/gohugoio/hugo/common/types"
	"github.com/spf13/cast"
	"reflect"
)

func (ns *Namespace) Update(index any, value any, mapOrSlice any) (string, error) {
	vmapOrSlice := reflect.ValueOf(mapOrSlice)

	vmapOrSlice, isNil := indirect(vmapOrSlice)
	if isNil {
		return "", fmt.Errorf("can't update a nil value")
	}

	var err error
	if vmapOrSlice.Kind() == reflect.Map {
		err = ns.updateMap(index, value, vmapOrSlice)
	} else if vmapOrSlice.Kind() == reflect.Slice {
		err = ns.updateSlice(index, value, vmapOrSlice)
	} else {
		return "", fmt.Errorf("target must be a map or a slice, got %T", mapOrSlice)
	}
	// This is used in templates, we need to return something.
	return "", err
}

func (ns *Namespace) updateMap(index, value any, mapp reflect.Value) error {
	vindex := reflect.ValueOf(index)
	if !vindex.Type().AssignableTo(mapp.Type().Key()) {
		return fmt.Errorf("cannot assign %v as a key in map %v", index, mapp.Type())
	}

	if types.IsNil(value) {
		// special case, delete the key
		mapp.SetMapIndex(vindex, reflect.Value{})
		return nil
	}

	vvalue := reflect.ValueOf(value)
	if !vvalue.Type().AssignableTo(mapp.Type().Elem()) {
		return fmt.Errorf("cannot assign %v as a value in map %v", value, mapp.Type())
	}

	mapp.SetMapIndex(vindex, vvalue)
	return nil
}

func (ns *Namespace) updateSlice(index, value any, slice reflect.Value) error {
	vvalue := reflect.ValueOf(value)
	if !vvalue.Type().AssignableTo(slice.Type().Elem()) {
		return fmt.Errorf("cannot assign %v as a value in slice %v", value, slice.Type())
	}

	indexv, err := cast.ToIntE(index)
	if err != nil {
		return err
	}

	slice.Index(indexv).Set(vvalue)
	return nil
}
