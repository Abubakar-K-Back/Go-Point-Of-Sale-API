package controllers

import (
	"fmt"
	db "point/of/sale/db"
	models "point/of/sale/models/OrderModel"
	model "point/of/sale/models/PaymentModel"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreatePayment(c *fiber.Ctx) error {
	fmt.Println("Creating a Payment Methods...")

	body := struct {
		Name          string `json:"name"`
		Type          string `json:"type"`
		Logo          string `json:"logo"`
		UpdatedAt     string `json:"updatedAt"`
		CreatedAt     string `json:"CreatedAt"`
		PaymentTypeId int    `json:"paymenttypeid"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(504).JSON(fiber.Map{
			"success": false,
			"Message": "UnprocessiableEntity",
		})
	}
	if len(body.Name) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
		})
	}

	var paymentTypes model.PaymentTypes
	db.DB.Where("name", body.Type).First(&paymentTypes)
	payment := models.Payment{
		Name:          body.Name,
		PaymentType:   body.Type,
		PaymentTypeId: int(paymentTypes.Id),
		Logo:          body.Logo,
	}
	db.DB.Create(&payment)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data": map[string]interface{}{
			"name":          payment.Name,
			"type":          payment.PaymentType,
			"logo":          payment.Logo,
			"updatedAt":     payment.UpdatedAt,
			"createdAt":     payment.CreatedAt,
			"paymentTypeId": payment.PaymentId,
		},
	})
}

func DeletePayment(c *fiber.Ctx) error {
	id := c.Params("paymentId")
	fmt.Println(id)
	payment := model.Payment{}
	db.DB.First(&payment, id)

	if payment.Name == "" {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"Message": "No Payment found against this payment id",
		})
	}

	result := db.DB.Delete(&payment, "payment_Id=?", id)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Payment removing failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success!",
	})
}

func UpdatePayment(c *fiber.Ctx) error {
	paymentId := c.Params("paymentId")

	var payment model.Payment

	db.DB.Find(&payment, "payment_id = ?", paymentId)

	if payment.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Payment not exist against this id",
		})
	}

	var updatePaymentData model.Payment

	c.BodyParser(&updatePaymentData)
	if updatePaymentData.Name == "" {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "fields are required"})
	}

	payment.Name = updatePaymentData.Name
	payment.Logo = updatePaymentData.Logo
	payment.PaymentType = updatePaymentData.PaymentType

	db.DB.Table("payments").Where("payment_id = ?", paymentId).Update("name", payment.Name)
	db.DB.Table("payments").Where("payment_id = ?", paymentId).Update("payment_type", payment.PaymentType)
	db.DB.Table("payments").Where("payment_id = ?", paymentId).Update("logo", payment.Logo)
	db.DB.Table("payments").Where("payment_id = ?", paymentId).Update("updated_at", time.Now().UTC())

	fmt.Println(payment.Name, payment.PaymentType)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    payment,
	})
}

func DetailPayment(c *fiber.Ctx) error {

	param := c.Params("paymentId")

	var payment model.Payment
	db.DB.Where("payment_id=?", param).First(&payment)
	if len(payment.Name) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success!",
		"data":    payment,
	})
}

func ListPayment(c *fiber.Ctx) error {

	//Token authenticate
	//subtotal,_ := strconv.Atoi(c.Query("subtotal"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var payment []model.Payment
	db.DB.Select("payment_id ,name,payment_type,logo,created_at,updated_at").Limit(limit).Offset(skip).Find(&payment).Count(&count)
	metaMap := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	categoriesData := map[string]interface{}{
		"payments": payment,
		"meta":     metaMap,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    categoriesData,
	})

}
