package models

import (
	"github.com/goravel/framework/database/orm"
)

type Resi struct {
	orm.Model
	TrackingNum string
	Expedition  string
	Status      string
	Details     string
	Email       string
}

func (r *Resi) TableName() string {
	return "resi"
}
