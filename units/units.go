package units

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/coreos/go-systemd/dbus"
	"github.com/coreos/go-systemd/unit"
)

const (
	mCConfigDir     = "/etc/mc/system/"
	sytemdConfigDir = "/etc/systemd/system"
)

type Unit struct {
	Name     string
	Enable   bool
	Contents []*unit.UnitOption
	DropIns  []DropIn
}

type DropIn struct {
	Name     string
	Contents []*unit.UnitOption
}

func Enable(units []Unit) error {
	conn, err := dbus.New()
	if err != nil {
		return fmt.Errorf("failed to connect to systemd D-Bus: %s", err)
	}

	if err = os.RemoveAll(mCConfigDir); err != nil {
		return fmt.Errorf("failed to clear config dir: %s got: %s", mCConfigDir, err)
	}

	for _, u := range units {
		if len(u.Contents) != 0 {
			path := unitPath(mCConfigDir, u)
			err := writeUnit(path, u.Contents)
			if err != nil {
				return fmt.Errorf("failed to write systemd unit file for %s got: %s", u.Name, err)
			}
			_, _, err = conn.EnableUnitFiles([]string{path}, false, true)
			if err != nil {
				return fmt.Errorf("failed to enable systemd unit file for %s got: %s", u.Name, err)
			}
		}
		for _, d := range u.DropIns {
			spath := dropInPath(mCConfigDir, u, d)
			dpath := dropInPath(sytemdConfigDir, u, d)
			if err = writeUnit(spath, d.Contents); err != nil {
				return fmt.Errorf("failed to write systemd dropin %s got: %s", d.Name, err)
			}
			if err = os.MkdirAll(filepath.Dir(dpath), 0755); err != nil {
				return fmt.Errorf("failed to link systemd dropin %s got: %s", d.Name, err)
			}
			if err = os.RemoveAll(dpath); err != nil {
				return fmt.Errorf("failed to remove symlink target: %s", err)
			}
			if err = os.Symlink(spath, dpath); err != nil {
				return fmt.Errorf("failed to link systemd dropin %s got: %s", d.Name, err)
			}
		}
	}

	if err = conn.Reload(); err != nil {
		return fmt.Errorf("failed to reload systemd: %s", err)
	}
	return nil
}

func unitPath(base string, u Unit) string {
	return path.Join(base, u.Name)
}

func dropInPath(base string, u Unit, d DropIn) string {
	return path.Join(base, fmt.Sprintf("%s.d", u.Name), d.Name)
}

func writeUnit(path string, u []*unit.UnitOption) error {
	if len(u) == 0 {
		return nil
	}
	os.MkdirAll(filepath.Dir(path), 0755)
	b, err := ioutil.ReadAll(unit.Serialize(u))
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)
}
