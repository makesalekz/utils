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

func (p ApplicationID) CompanyFullName() string {
	name := os.Getenv("COMPANY_FULL_NAME")
	if name != "" {
		return name
	}

	if p == AppCalendaria || p == AppPMS {
		return "TOO \"AXIO\""
	}

	return "TOO \"AXIO\""
}
