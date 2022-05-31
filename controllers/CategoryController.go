package controllers

import (
	"fmt"
	db "point/of/sale/db"
	model "point/of/sale/models"
	util "point/of/sale/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ListCategory(c *fiber.Ctx) error {

	auth := c.Get("Authorization")
	if len(auth) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}
	if err := util.AuthToken(util.SplitToken(auth)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Token expired or invalid",
		})
	}

	limit := c.Query("limit")
	skip := c.Query("skip")

	intLimit, _ := strconv.Atoi(limit)
	intSkip, _ := strconv.Atoi(skip)

	var catagory []model.Category
	res := db.DB.Select([]string{"category_id", "name"}).Limit(intLimit).Offset(intSkip).Find(&catagory)

	fmt.Println(res.RowsAffected)
	return c.Status(404).JSON(fiber.Map{
		"success": true,
		"Message": "Sucess",
		"data":    catagory,
		"meta": map[string]interface{}{
			"total": res.RowsAffected,
			"limit": limit,
			"skip":  skip,
		},
	})

}

func DetailCategory(c *fiber.Ctx) error {

	auth := c.Get("Authorization")
	if len(auth) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}
	if err := util.AuthToken(util.SplitToken(auth)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  "error",
			"message": "Token expired or invalid",
		})
	}

	param := c.Params("categoryId")

	var category model.Category
	db.DB.Where("category_Id=?", param).First(&category)
	if len(category.Name) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
		})
	}
	fmt.Println(param)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success!",
		"data":    category,
	})
}

func CreateCategory(c *fiber.Ctx) error {
	fmt.Println("Creating a New Catogory...")
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
			"Message": "Empty Body",
		})
	}

	catagory := model.Category{
		Name:      data.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	db.DB.Create(&catagory)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    catagory,
	})

}

func UpdateCategory(c *fiber.Ctx) error {

	catId := c.Params("categoryId")
	var category model.Category
	db.DB.Find(&category, "category_id = ?", catId)

	if category.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Category not exist against this id",
		})
	}

	var updateCategoryData model.Category
	c.BodyParser(&updateCategoryData)
	if updateCategoryData.Name == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "fields are required"})
	}

	category.Name = updateCategoryData.Name
	db.DB.Table("categories").Where("category_id = ?", category.CategoryId).Update("name", category.Name)
	db.DB.Table("categories").Where("category_id = ?", category.CategoryId).Update("updated_at", time.Now().UTC())

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    category,
	})

}

func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("categoryId")
	fmt.Println(id)
	catagory := model.Category{}
	db.DB.First(&catagory, id)

	if catagory.Name == "" {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"Message": "No catagory found against this Category id",
		})
	}

	result := db.DB.Delete(&catagory, "category_id=?", id)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Cashier removing failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success!",
	})

}
