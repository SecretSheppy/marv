package extlib

import (
	"plugin"
)

func Load(p string) ([]*Extension, error) {
	exts := make([]*Extension, 0)

	// TODO: loop
	plg, err := plugin.Open("")
	if err != nil {
		return nil, err
	}

	v, err := plg.Lookup("Ext")
	if err != nil {
		return nil, err
	}

	ext, ok := v.(*Extension)
	if !ok {
		return nil, ErrNoExtDeclaration
	}

	exts = append(exts, ext)
	// TODO: end loop

	return exts, nil
}
