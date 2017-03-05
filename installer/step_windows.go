package installer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/deviceio/shared/install"
	"github.com/google/uuid"
)

type step struct {
}

func (t *step) Install(i *install.Installer, args install.Args) error {
	var err error

	orgdir := fmt.Sprintf("c:/PROGRA~1/deviceio/agent/%v", args["org"].(string))
	bindir := fmt.Sprintf("%v/bin", orgdir)
	exepath := fmt.Sprintf("%v/deviceio-agent.exe", bindir)
	confpath := fmt.Sprintf("%v/config.json", orgdir)

	if err = i.MakePaths([]string{
		orgdir,
		bindir,
	}); err != nil {
		return err
	}

	if !i.Exists(exepath) {
		exe, err := os.Executable()

		if err != nil {
			return err
		}

		if err = i.Copy(exepath, exe); err != nil {
			return err
		}
	}

	if !i.Exists(confpath) {
		type config struct {
			ID                         string
			Tags                       []string
			TransportHost              string
			TransportPort              int
			TransportPasscodeHash      string
			TransportPasscodeSalt      string
			TransportDisableKeyPinning bool
			TransportAllowSelfSigned   bool
		}

		f, err := os.OpenFile(confpath, os.O_RDWR|os.O_CREATE, 0700)
		defer f.Close()

		if err != nil {
			return err
		}

		uuid, err := uuid.NewRandom()

		if err != nil {
			return err
		}

		enc := json.NewEncoder(f)
		enc.SetIndent("", "    ")

		err = enc.Encode(&config{
			ID:                       uuid.String(),
			Tags:                     []string{},
			TransportHost:            args["transportHost"].(string),
			TransportPort:            args["transportPort"].(int),
			TransportAllowSelfSigned: args["transportSelfSigned"].(bool),
		})
	}

	return nil
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
