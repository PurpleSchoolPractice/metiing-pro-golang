package eventParticipant

// EventParticipant представляет связь "многие ко многим" между событиями и пользователями.
type EventParticipant struct {
	EventID uint `gorm:"primaryKey"`
	UserID  uint `gorm:"primaryKey"`

	Event interface{} `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`
	User  interface{} `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// NewEventParticipant создает новую связь между событием и участником
func NewEventParticipant(eventID, userID uint) *EventParticipant {
	return &EventParticipant{
		EventID: eventID,
		UserID:  userID,
	}
}
