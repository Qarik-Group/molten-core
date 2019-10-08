package units

import (
	"github.com/coreos/go-systemd/unit"
)

var (
	BUCC []Unit = []Unit{{
		Name: "bucc.service",
		Contents: []*unit.UnitOption{
			unit.NewUnitOption("Unit", "Description", "BUCC - BOSH UAA Credhub and Concourse"),
			unit.NewUnitOption("Unit", "After", "docker.service"),
			unit.NewUnitOption("Unit", "Requires", "docker.service"),

			unit.NewUnitOption("Service", "Type", "oneshot"),
			unit.NewUnitOption("Service", "ExecStart", "/opt/bin/mc bucc-up"),
			unit.NewUnitOption("Service", "RemainAfterExit", "true"),
			unit.NewUnitOption("Service", "StandardOutput", "journal"),

			unit.NewUnitOption("Install", "WantedBy", "multi-user.target"),
		},
	},
		{
			Name: "bucc-configs.service",
			Contents: []*unit.UnitOption{
				unit.NewUnitOption("Unit", "Description", "Updates BOSH {cloud,cpi,runtime}-configs"),
				unit.NewUnitOption("Unit", "After", "bucc.service"),
				unit.NewUnitOption("Unit", "Requires", "bucc.service"),

				unit.NewUnitOption("Service", "Type", "oneshot"),
				unit.NewUnitOption("Service", "ExecStart", "/opt/bin/mc update-bucc-configs"),
				unit.NewUnitOption("Service", "RemainAfterExit", "true"),
				unit.NewUnitOption("Service", "StandardOutput", "journal"),

				unit.NewUnitOption("Install", "WantedBy", "multi-user.target"),
			},
		},
		// bucc-sync-dns is a workaround for https://github.com/cloudfoundry/bosh/issues/2103
		{
			Name: "bucc-sync-dns.service",
			Contents: []*unit.UnitOption{
				unit.NewUnitOption("Unit", "Description", "Forcefully sync bosh-dns"),
				unit.NewUnitOption("Unit", "After", "bucc.service"),
				unit.NewUnitOption("Unit", "Requires", "bucc.service"),

				unit.NewUnitOption("Service", "Type", "oneshot"),
				unit.NewUnitOption("Service", "ExecStart", "/bin/bash -c \"docker exec $(jq -r '.current_vm_cid' /var/lib/moltencore/bucc/state.json) /var/vcap/jobs/director/bin/trigger-one-time-sync-dns\""),
				unit.NewUnitOption("Service", "StandardOutput", "journal"),
			},
		},
		{
			Name: "bucc-sync-dns.timer",
			Contents: []*unit.UnitOption{
				unit.NewUnitOption("Unit", "Description", "Run sync-dns.service every 30 seconds"),

				unit.NewUnitOption("Timer", "OnCalendar", "*:*:0,30"),
			},
		},
	}
)
