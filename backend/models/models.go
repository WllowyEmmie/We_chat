package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type UserRoom struct {
	RoomID uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	UserID uuid.UUID `gorm:"primaryKey;type:uuid;not null"`

	Room *Room `gorm:"foreignKey:RoomID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type User struct {
	ID       uuid.UUID  `gorm:"type:uuid;primaryKey;" json:"id"`
	UserName string     `gorm:"type:varchar(100);not null" json:"username"`
	Email    string     `gorm:"type:varchar(100);not null;unique;" json:"email"`
	Password string     `gorm:"not null" json:"password"`
	Messages []*Message `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"messages"`
	Rooms    []*Room    `gorm:"many2many:user_rooms;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rooms"`
}

type Message struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	RoomID    *uuid.UUID `gorm:"type:uuid;" json:"room_id"`
	Room      *Room      `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"room"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:current_timestamp;" json:"created_at"`
}

type Room struct {
	ID       uuid.UUID                `gorm:"type:uuid;primaryKey" json:"id"`
	Name     string                   `gorm:"type:varchar(100);not null" json:"name"`
	Clients  map[*websocket.Conn]bool `gorm:"-" json:"-"`
	Members  []*User                  `gorm:"many2many:user_rooms;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Messages []*Message               `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"messages"`
}

// Auto-generate IDs
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}
