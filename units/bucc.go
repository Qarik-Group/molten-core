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
			unit.NewUnitOption("Service", "ExecStart", "/opt/bin/mc bucc-up"),
			unit.NewUnitOption("Service", "RemainAfterExit", "true"),
			unit.NewUnitOption("Service", "StandardOutput", "journal"),

			unit.NewUnitOption("Install", "WantedBy", "multi-user.target"),
		},
	}
)
