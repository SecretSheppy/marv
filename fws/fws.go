package fws

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws/mutest_rs"
	"github.com/SecretSheppy/marv/fws/pitest"
)

func Frameworks() []fwlib.Framework {
	return []fwlib.Framework{
		mutest_rs.NewMutestRS(),
		pitest.NewPitest(),
	}
}

func FrameworksMap() map[string]fwlib.Framework {
	fws := make(map[string]fwlib.Framework)
	for _, fw := range Frameworks() {
		fws[fw.Meta().Name] = fw
	}
	return fws
}
