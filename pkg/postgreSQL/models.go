package postgreSQL

import (
	"gorm.io/gorm"
)

// @Schema
// @Description User structure for registration and authentication
// @type User struct {
//   ID           uint   `json:"id"`
//   Username     string `json:"username"`
//   PasswordHash string `json:"password_hash"`
//   Email        string `json:"email"`
//   CreatedAt    time.Time `json:"created_at"`
//   UpdatedAt    time.Time `json:"updated_at"`
// }

// Dialog представляет диалог пользователя с AI.
// @Description  Структура диалога с AI
type Dialog struct {
	gorm.Model
	// ID пользователя, которому принадлежит диалог
	UserID uint `gorm:"not null"`
	// Сообщение пользователя
	Message string `gorm:"not null"`
	// Ответ AI
	Response string `gorm:"not null"`
}

// User представляет пользователя системы.
// @Description  Структура пользователя
type User struct {
	gorm.Model
	// Имя пользователя (уникальное)
	Username string `gorm:"unique;not null" json:"username"`
	// Хэш пароля (скрыт в JSON)
	PasswordHash string `gorm:"not null" json:"password"`
	// Список диалогов пользователя
	Dialogs []Dialog `gorm:"foreignKey:UserID" json:"dialogs"`
}
