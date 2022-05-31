package controllers

import (
	"fmt"
	"point/of/sale/db"
	"point/of/sale/models"
	auth "point/of/sale/utils"

	"github.com/gofiber/fiber/v2"
)

func GetPasscode(c *fiber.Ctx) error {
	params := c.Params("cashierId")
	cashier := models.Cashier{}
	db.DB.Where("cashier_id=?", params).First(&cashier)

	if len(cashier.Name) == 0 {

		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
			"error":   map[string]interface{}{},
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data": map[string]interface{}{
			"passcode": cashier.Passcode,
		},
	})

}

func VerifyPasscode(c *fiber.Ctx) error {

	id := c.Params("cashierId")
	data := struct {
		Passcode string `json:"passcode"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"Message": "Passcode is required",
		})
	}
	fmt.Println("Auth in Progress...")

	cashier := models.Cashier{}

	db.DB.Where("passcode = ? AND cashier_Id= ?", data.Passcode, id).Find(&cashier)

	if cashier.Name == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"Message": "Not Match",
		})
	}
	return c.Status(504).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data": map[string]interface{}{
			"token": auth.AuthenticateJWT(data.Passcode),
		},
	})
}

func VerifyLogoutPasscode(c *fiber.Ctx) error {
	id := c.Params("cashierId")
	data := struct {
		Passcode string `json:"passcode"`
	}{}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"Message": "Empty Body",
		})
	}
	cashier := models.Cashier{}

	db.DB.Where("passcode = ? AND cashier_Id= ?", data.Passcode, id).Find(&cashier)

	if cashier.Name == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"Message": "Not Match",
		})
	}
	return c.Status(504).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
	})
}
