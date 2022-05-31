package paymentmodel

import "time"

type Payment struct {
	PaymentId     uint      `json:"paymentid" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name          string    `json:"name"`
	Logo          string    `json:"logo"`
	PaymentTypeId int       `json:"paymenttypeid"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PaymentType   string    `json:"type"`
}
type PaymentTypes struct {
	Id        int       `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
