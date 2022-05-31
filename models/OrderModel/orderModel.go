package ordermodel

import "time"

type ProductResponseOrder struct {
	ProductId        int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name             string   `json:"name"`
	Price            int      `json:"price"`
	Qty              int      `json:"qty"`
	Discount         Discount `json:"discount"`
	TotalNormalPrice int      `json:"totalNormalPrice"`
	TotalFinalPrice  int      `json:"totalFinalPrice"`
}

type RevenueResponse struct {
	PaymentTypeId int    `json:"paymentTypeId"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	TotalAmount   int    `json:"totalAmount"`
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

type Product struct {
	ProductId  int    `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
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

type Order struct {
	Id             int       `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	CashierID      int       `json:"cashierId"`
	PaymentTypesId int       `json:"paymentTypesId"`
	TotalPrice     int       `json:"totalPrice"`
	TotalPaid      int       `json:"totalPaid"`
	TotalReturn    int       `json:"totalReturn"`
	ReceiptId      string    `json:"receiptId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	ProductId      string    `json:"productId"`
	IsDownloaded   int       `json:"isdownloaded"`
	Quantities     string    `json:"quantities"`
}

type SubtotalResponse struct {
	SubTotal int                  `json:"subtotal"`
	Products ProductResponseOrder `json:"products"`
}
type Cashier struct {
	CashierId uint      `json:"cashierid" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	Passcode  uint      `json:"passcode"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type PaymentTypes struct {
	Id        int       `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Payment struct {
	PaymentId     uint      `json:"paymentid" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name          string    `json:"name"`
	Logo          string    `json:"logo"`
	PaymentTypeId int       `json:"paymenttypeid"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PaymentType   string    `json:"type"`
}

type SoldResponse struct {
	ProductId   int    `json:"productId"`
	Name        string `json:"name"`
	TotalQty    int    `json:"totalQty"`
	TotalAmount int    `json:"totalAmount"`
}
