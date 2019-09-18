package units

import (
	"github.com/coreos/go-systemd/unit"
)

type Unit struct {
	DropIn []*unit.UnitOption
	// TODO
}

func Update(unit Unit) error {
	return nil
}
