package units

import (
	"github.com/coreos/go-systemd/unit"
)

var (
	BUCC Unit = Unit{
		Name: "bucc.service",
		Contents: []*unit.UnitOption{
			unit.NewUnitOption("Unit", "Description", "BUCC - BOSH UAA Credhub and Concourse"),
			unit.NewUnitOption("Unit", "After", "docker.service"),
			unit.NewUnitOption("Unit", "Wants", "docker.service"),
			unit.NewUnitOption("Unit", "Requires", "docker.service"),

			unit.NewUnitOption("Service", "Type", "oneshot"),
			unit.NewUnitOption("Service", "ExecStartPre", "mkdir -p /var/lib/bucc"),
		},
	}
)
