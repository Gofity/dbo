package dbo

import "gorm.io/gorm"

func NewSession(db *gorm.DB) (tx *gorm.DB) {
	if db == nil {
		return
	}

	return db.Session(&gorm.Session{
		NewDB: true,
	})
}
