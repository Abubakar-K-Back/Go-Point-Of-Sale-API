package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	db "point/of/sale/db"
	"point/of/sale/models"
	model "point/of/sale/models"

	"github.com/gofiber/fiber/v2"
)

func CashierDetails(c *fiber.Ctx) error {

	param := c.Params("cashierId")

	var cashier model.Cashier
	db.DB.Where("cashier_id=?", param).First(&cashier)
	if len(cashier.Name) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
			"error":   map[string]interface{}{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success!",
		"data":    cashier,
	})
}

func CreateCashiar(c *fiber.Ctx) error {

	fmt.Println("Creating a New Cashier...")
	data := struct {
		Name string `json:"name"`
	}{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(504).JSON(fiber.Map{
			"success": false,
			"Message": "UnprocessiableEntity",
		})
	}
	if len(data.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"Message": "Bad Request",
			"error":   map[string]interface{}{},
		})
	}

	cashier := model.Cashier{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      data.Name,
		Passcode:  uint(rand.Intn(1000000)),
	}

	db.DB.Create(&cashier)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    cashier,
	})
}

func UpdateCashier(c *fiber.Ctx) error {
	cashierId := c.Params("cashierId")
	var cashier models.Cashier
	db.DB.Find(&cashier, "cashier_id = ?", cashierId)

	if cashier.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Cashier Not Found",
		})
	}

	var updateCashierData models.Cashier
	c.BodyParser(&updateCashierData)
	if updateCashierData.Name == "" {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Empty Body"})
	}
	cashier.Name = updateCashierData.Name
	db.DB.Table("cashiers").Where("cashier_Id = ?", cashierId).Update("name", cashier.Name)
	db.DB.Table("cashiers").Where("cashier_Id = ?", cashierId).Update("updated_at", time.Now().UTC())

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    cashier,
	})
}

func GetCashiers(c *fiber.Ctx) error {
	limit := c.Query("limit")
	skip := c.Query("skip")

	intLimit, _ := strconv.Atoi(limit)
	intSkip, _ := strconv.Atoi(skip)

	var cashier []model.Cashier
	res := db.DB.Select([]string{"cashier_id", "name"}).Limit(intLimit).Offset(intSkip).Find(&cashier)

	fmt.Println(res.RowsAffected)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Sucess",
		"data":    cashier,
		"meta": map[string]interface{}{
			"total": res.RowsAffected,
			"limit": intLimit,
			"skip":  skip,
		},
	})
}

func DeleteCashier(c *fiber.Ctx) error {
	id := c.Params("cashierId")
	fmt.Println(id)
	cashier := model.Cashier{}
	db.DB.First(&cashier, id)

	if cashier.Name == "" {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"Message": "No cashier found against this cashier id",
		})
	}

	result := db.DB.Delete(&cashier, "cashier_Id=?", id)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Cashier Not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success!",
	})
}
