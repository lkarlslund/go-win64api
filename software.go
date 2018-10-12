// +build windows,amd64
package winapi

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

type Software struct {
	DisplayName    string `json:"displayName"`
	DisplayVersion string `json:"displayVersion"`
	Arch           string `json:"arch"`
}

func (s *Software) Name() string {
	return s.DisplayName
}

func (s *Software) Version() string {
	return s.DisplayVersion
}

func (s *Software) Architecture() string {
	return s.Arch
}

func InstalledSoftwareList() ([]Software, error) {
	sw64, err := getSoftwareList(`SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, "X64")
	if err != nil {
		return nil, err
	}
	sw32, err := getSoftwareList(`SOFTWARE\Wow6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "X32")
	if err != nil {
		return nil, err
	}

	return append(sw64, sw32...), nil
}

func getSoftwareList(baseKey string, arch string) ([]Software, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, baseKey, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, fmt.Errorf("Error reading from registry: %s", err.Error())
	}
	defer k.Close()

	swList := make([]Software, 0)

	subkeys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return nil, fmt.Errorf("Error reading subkey list from registry: %s", err.Error())
	}
	for _, sw := range subkeys {
		sk, err := registry.OpenKey(registry.LOCAL_MACHINE, baseKey+`\`+sw, registry.QUERY_VALUE)
		if err != nil {
			return nil, fmt.Errorf("Error reading from registry (subkey %s): %s", sw, err.Error())
		}

		dn, _, err := sk.GetStringValue("DisplayName")
		if err == nil {
			swv := Software{DisplayName: dn, Arch: arch}

			dv, _, err := sk.GetStringValue("DisplayVersion")
			if err == nil {
				swv.DisplayVersion = dv
			}

			swList = append(swList, swv)
		}
	}

	return swList, nil
}
