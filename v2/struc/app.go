package struc

import (
	"os"
)

type ApplicationID string

const (
	AppCalendaria ApplicationID = "calendaria"
	AppPMS        ApplicationID = "pms"
)

func (p ApplicationID) Value() string {
	return string(p)
}

func (p ApplicationID) IsValid() bool {
	switch p {
	case AppCalendaria:
		return true
	case AppPMS:
		return true
	}
	return false
}

func (p ApplicationID) BrandName() string {
	name := os.Getenv("BRAND_NAME")
	if name != "" {
		return name
	}

	switch p {
	case AppCalendaria:
		return "AIgenda"
	case AppPMS:
		return "BasQaru"
	}

	return "AXIO"
}
