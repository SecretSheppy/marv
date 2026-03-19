package fws

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws/mutest_rs"
	"github.com/SecretSheppy/marv/fws/pitest"
)

func Frameworks() []fwlib.Framework {
	return []fwlib.Framework{
		&mutest_rs.MutestRS{},
		&pitest.Pitest{},
	}
}
