package units

import (
	"github.com/coreos/go-systemd/unit"
)

var (
	Docker Unit = Unit{
		DropIn: []*unit.UnitOption{},
	}
)
