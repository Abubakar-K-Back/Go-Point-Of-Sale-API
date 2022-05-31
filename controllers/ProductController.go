package controllers

import (
	"fmt"
	"log"
	config "point/of/sale/db"
	models "point/of/sale/models/productMod"
	util "point/of/sale/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProducts(c *fiber.Ctx) error {

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
			"success": false,
			"message": "Unauthorized",
			"error":   map[string]interface{}{},
		})
	}

	limit := c.Query("limit")
	skip := c.Query("skip")
	categoryId := c.Query("categoryId")
	productName := c.Query("q")
	intLimit, _ := strconv.Atoi(limit)
	intSkip, _ := strconv.Atoi(skip)
	var products []models.Product
	var count int64
	productsRes := make([]*models.ProductResponse, 0)
	fmt.Println(len(productName))
	if len(productName) <= 0 {
		fmt.Println("isme agya")
		result := config.DB.Where("category_id = ?", categoryId).Limit(intLimit).Offset(intSkip).Find(&products).Count(&count)
		fmt.Println(products)
		if result.Error != nil {
			Response := map[string]interface{}{
				"success": false,
				"message": "Error",
			}
			return (c.JSON(Response))
		}
		var category models.Category
		var discount models.Discount
		fmt.Println(len(products))
		for i := 0; i < len(products); i++ {
			fmt.Println("----")
			fmt.Println(products[i])
			fmt.Println(products[i].CategoryId)
			result1 := config.DB.Table("categories").Where("category_id = ?", products[i].CategoryId).Find(&category)
			fmt.Println(category)

			if result1.Error != nil {
				Response := map[string]interface{}{
					"success": false,
					"message": "Error",
				}
				return (c.JSON(Response))
			}
			fmt.Println(category.Name)

			result2 := config.DB.Where("id = ?", products[i].DiscountId).Limit(intLimit).Offset(intSkip).Find(&discount).Count(&count)
			fmt.Println("---discount---", discount.Id)

			if result2.Error != nil {
				Response := map[string]interface{}{
					"success": false,
					"message": "Error",
				}
				return (c.JSON(Response))
			}
			//productsRes =
			productsRes = append(productsRes,
				&models.ProductResponse{
					Id:       products[i].Id,
					Sku:      products[i].Sku,
					Name:     products[i].Name,
					Stock:    products[i].Stock,
					Price:    products[i].Price,
					Image:    products[i].Image,
					Category: category,
					Discount: discount,
				},
			)
		}

		meta := map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		}
		Response := map[string]interface{}{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"products": productsRes,
			},
			"meta": meta,
		}
		return (c.JSON(Response))
	}
	// result := config.DB.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&pr)
	result := config.DB.Where("category_Id = ? AND name= ?", categoryId, productName).Limit(intLimit).Offset(intSkip).Find(&products).Count(&count)
	if result.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
		}
		return (c.JSON(Response))
	}
	var category models.Category
	var discount models.Discount
	for i := 0; i < len(products); i++ {
		fmt.Println("----")
		fmt.Println(products[i])
		fmt.Println(products[i].CategoryId)
		result1 := config.DB.Where("id = ?", products[i].CategoryId).Find(&category)
		if result1.Error != nil {
			Response := map[string]interface{}{
				"success": false,
				"message": "Error",
			}
			return (c.JSON(Response))
		}
		fmt.Println(category.Name)
		result2 := config.DB.Where("id = ?", products[i].DiscountId).Limit(intLimit).Offset(intSkip).Find(&discount).Count(&count)

		if result2.Error != nil {
			Response := map[string]interface{}{
				"success": false,
				"message": "Error",
			}
			return (c.JSON(Response))
		}
		//productsRes =
		productsRes = append(productsRes,
			&models.ProductResponse{
				Id:       products[i].Id,
				Sku:      products[i].Sku,
				Name:     products[i].Name,
				Stock:    products[i].Stock,
				Price:    products[i].Price,
				Image:    products[i].Image,
				Category: category,
				Discount: discount,
			},
		)
	}

	meta := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
		"data":    productsRes,
		"meta":    meta,
	}
	return (c.JSON(Response))
}

func GetProductsById(c *fiber.Ctx) error {

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

	productId := c.Params("productId")
	productsRes := make([]*models.ProductResponse, 0)
	var products []models.Product
	// result := config.DB.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&pr)
	result := config.DB.Where("id = ? ", productId).Find(&products)
	if result.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
		}
		return (c.JSON(Response))
	}
	var category models.Category
	var discount models.Discount
	for i := 0; i < len(products); i++ {
		fmt.Println("----")
		fmt.Println(products[i])
		fmt.Println(products[i].CategoryId)
		result1 := config.DB.Where("category_id = ?", products[i].CategoryId).Find(&category)
		if result1.Error != nil {
			Response := map[string]interface{}{
				"success": false,
				"message": "Error",
			}
			return (c.JSON(Response))
		}
		fmt.Println(category.Name)
		result2 := config.DB.Where("id = ?", products[i].DiscountId).Find(&discount)

		if result2.Error != nil {
			Response := map[string]interface{}{
				"success": false,
				"message": "Error",
			}
			return (c.JSON(Response))
		}
		//productsRes =
		productsRes = append(productsRes,
			&models.ProductResponse{
				Id:       products[i].Id,
				Sku:      products[i].Sku,
				Name:     products[i].Name,
				Stock:    products[i].Stock,
				Price:    products[i].Price,
				Image:    products[i].Image,
				Category: category,
				Discount: discount,
			},
		)
	}
	fmt.Println(productsRes)
	if len(productsRes) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   map[string]interface{}{},
		})
	}
	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
		"data":    productsRes,
	}
	return (c.JSON(Response))
}

type Products struct {
	Products     models.Product
	CategoriesId string `json:"categories_Id"`
}
type ProdDiscount struct {
	ProductId  int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string   `json:"sku"`
	Name       string   `json:"name"`
	Stock      int      `json:"stock"`
	Price      int      `json:"price"`
	Image      string   `json:"image"`
	CategoryId int      `json:"categoryId"`
	Discount   Discount `json:"discount"`
}
type Discount struct {
	Qty       int    `json:"qty"`
	Types     string `json:"type"`
	Result    int    `json:"result"`
	ExpiredAt int    `json:"expiredAt"`
}

func AddProducts(c *fiber.Ctx) error {
	var data ProdDiscount
	err := c.BodyParser(&data)
	fmt.Println(err)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Body required",
		})
	}
	fmt.Println(data)
	discount := models.Discount{
		Qty:       data.Discount.Qty,
		Type:      data.Discount.Types,
		Result:    data.Discount.Result,
		ExpiredAt: float32(data.Discount.ExpiredAt),
	}
	resultDiscount := config.DB.Create(&discount)
	if resultDiscount.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
			"error":   map[string]interface{}{},
		}
		return (c.JSON(Response))
	}
	fmt.Println(discount.Id)
	product := models.Product{
		Name:       data.Name,
		Image:      data.Image,
		CategoryId: data.CategoryId,
		DiscountId: discount.Id,
		Price:      data.Price,
		Stock:      data.Stock,
	}
	fmt.Print(product)
	result := config.DB.Create(&product)
	if result.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
		}
		return (c.JSON(Response))
	}
	fmt.Println(product.Id)
	Id := fmt.Sprintf("%04d", product.Id)
	result1 := config.DB.Table("products").Where("id = ?", product.Id).Update("sku", Id)
	if result1.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
		}
		return (c.JSON(Response))
	}
	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
		"data":    product,
	}
	return (c.JSON(Response))
}

func UpdateProduct(c *fiber.Ctx) error {
	var data models.Product
	productId := c.Params("productsId")
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatalf("Product error in post request %v", err)
	}

	//Validation rules
	if data.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Product Name is required",
		})
	}

	if data.Price <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Price field is required",
		})
	}

	if data.CategoryId <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Id field is required",
		})
	}
	if data.Image == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Image field is required",
		})
	}
	if data.Stock <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": " Stock field is required",
		})
	}

	product := models.Product{
		Name:       data.Name,
		Image:      data.Image,
		CategoryId: data.CategoryId,
		Price:      data.Price,
		Stock:      data.Stock,
	}
	var products models.Product
	fmt.Print(product)
	fmt.Println("------>", productId)
	result := config.DB.Model(products).Where("id= ?", productId).Updates(&product)
	if result.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
		}
		return (c.JSON(Response))
	}
	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
		"data":    product,
	}
	return (c.JSON(Response))

}

func DeleteProduct(c *fiber.Ctx) error {
	productId := c.Params("productId")
	var product models.Product
	result := config.DB.Where("id = ?", productId).Delete(&product)
	if result.Error != nil {
		Response := map[string]interface{}{
			"success": false,
			"message": "Error",
		}
		return (c.JSON(Response))
	}

	if (result.RowsAffected) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Token expired or invalid",
			"error":   map[string]interface{}{},
		})
	}
	Response := map[string]interface{}{
		"success": true,
		"message": "Success",
	}
	return (c.JSON(Response))
}
