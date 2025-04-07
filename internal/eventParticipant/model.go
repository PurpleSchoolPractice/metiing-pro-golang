package eventParticipant

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
)

// EventParticipant представляет связь "многие ко многим" между событиями и пользователями.
type EventParticipant struct {
	EventID uint `gorm:"primaryKey"`
	UserID  uint `gorm:"primaryKey"`

	Event event.Event `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`
	User  user.User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
