package productMod

import "time"

type Product struct {
	Id         int    `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string `json:"sku"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryId int    `json:"categoryId"`
	DiscountId int    `json:"discountId"`
}
type ProductOrder struct {
	Id         int    `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string `json:"sku"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryId int    `json:"categoryId"`
	DiscountId int    `json:"discountId"`
}

type Discount struct {
	Id              int     `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Qty             int     `json:"qty"`
	Type            string  `json:"type"`
	Result          int     `json:"result"`
	ExpiredAt       float32 `json:"expiredAt"`
	ExpiredAtFormat string  `json:"expiredAtFormat"`
	StringFormat    string  `json:"stringFormat"`
}

type ProductResponse struct {
	Id       int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku      string   `json:"sku"`
	Name     string   `json:"name"`
	Stock    int      `json:"stock"`
	Price    int      `json:"price"`
	Image    string   `json:"image"`
	Category Category `json:"category"`
	Discount Discount `json:"discount"`
}
type Category struct {
	CategoryId uint      `json:"categoryid" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type ProductResponseOrder struct {
	ProductId int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku       string   `json:"sku"`
	Name      string   `json:"name"`
	Stock     int      `json:"stock"`
	Price     int      `json:"price"`
	Image     string   `json:"image"`
	Discount  Discount `json:"discount"`
}
