package struc

import "slices"

// ------------------------------- Structs -----------------------------

type AuthIds struct {
	ActorId  int64
	TenantId int64
}

// ------------------------------- Enums -----------------------------

type Provider string

const (
	Calendaria Provider = "CALENDARIA"
	Google     Provider = "GOOGLE"
	Outlook    Provider = "OUTLOOK"
	Apple      Provider = "APPLE"
)

func providerValues() []Provider {
	return []Provider{Calendaria, Google, Outlook, Apple}
}

func (Provider) Values() (kinds []string) {
	for _, s := range providerValues() {
		kinds = append(kinds, string(s))
	}
	return
}

func (p Provider) Value() string {
	return string(p)
}

func (p Provider) IsValid() bool {
	return slices.Contains(providerValues(), p)
}
