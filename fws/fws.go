package fws

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws/mutest_rs"
)

func Frameworks() []fwlib.Framework {
	return []fwlib.Framework{
		&mutest_rs.MutestRS{},
	}
}
