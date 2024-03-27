package struc

// ------------------------------- Structs -----------------------------

type FirebaseNotification struct {
	Type     NotificationType  `json:"type,omitempty"`
	UsersIds []int64           `json:"users_ids"`
	Title    string            `json:"title,omitempty"`
	Body     string            `json:"body,omitempty"`
	Image    string            `json:"image,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

// ------------------------------- Enums -----------------------------

type NotificationType string

const (
	Common  NotificationType = "COMMON"
	Event   NotificationType = "EVENT"
	Contact NotificationType = "CONTACT"
	Tasks   NotificationType = "TASKS"
)

func notificationTypeValues() []NotificationType {
	return []NotificationType{Common, Event, Contact, Tasks}
}

func (NotificationType) Values() (kinds []string) {
	for _, value := range notificationTypeValues() {
		kinds = append(kinds, string(value))
	}
	return
}

func (m NotificationType) Value() string {
	return string(m)
}

func (m NotificationType) IsValid() bool {
	for _, value := range notificationTypeValues() {
		if m == value {
			return true
		}
	}
	return false
}
