package db_helpers

import (
	"github.com/jinzhu/gorm"
)

func GetOrCreate(db *gorm.DB, model interface{}) (bool, error) {
	if res := db.Where(model).Take(model); res.Error != nil {
		if res.RecordNotFound() {
			db.NewRecord(model)
			db.Create(model)
			return true, nil
		} else {
			return false, res.Error
		}
	}

	return false, nil
}
