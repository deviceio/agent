package installer

import (
	"errors"

	"github.com/deviceio/shared/install"
)

type step struct {
}

func (t *step) Install(i *install.Installer, args install.Args) error {
	return errors.New("Darwin not yet suported for installer")
}

func (t *step) Uninstall(i *install.Installer, args install.Args) error {
	return nil
}

func (t *step) Upgrade(i *install.Installer, args install.Args) error {
	return nil
}

func (t *step) Repair(i *install.Installer, args install.Args) error {
	return nil
}
