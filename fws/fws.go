package fws

import (
	"github.com/SecretSheppy/marv/fwlib"
	"github.com/SecretSheppy/marv/fws/mull"
	"github.com/SecretSheppy/marv/fws/mutest_rs"
	"github.com/SecretSheppy/marv/fws/pitest"
	"github.com/SecretSheppy/marv/fws/stryker4s"
	"github.com/SecretSheppy/marv/fws/stryker_js"
	"github.com/SecretSheppy/marv/fws/stryker_net"
)

func Frameworks() []fwlib.Framework {
	return []fwlib.Framework{
		mull.NewMull(),
		mutest_rs.NewMutestRS(),
		pitest.NewPitest(),
		stryker4s.NewStryker4s(),
		stryker_js.NewStrykerJS(),
		stryker_net.NewStrykerNet(),
	}
}

func FrameworksMap() map[string]fwlib.Framework {
	fws := make(map[string]fwlib.Framework)
	for _, fw := range Frameworks() {
		fws[fw.Meta().Name] = fw
	}
	return fws
}

func ActiveFrameworks(yml []byte) ([]fwlib.Framework, error) {
	active := make([]fwlib.Framework, 0)
	for _, fw := range Frameworks() {
		loaded, err := fw.Yaml().Load(yml)
		if err != nil {
			return nil, err
		}
		if !loaded {
			continue
		}
		if err := fw.LoadResults(); err != nil {
			return nil, err
		}
		active = append(active, fw)
	}
	return active, nil
}
