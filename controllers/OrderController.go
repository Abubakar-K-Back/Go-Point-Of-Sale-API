package controllers

import (
	"fmt"
	"os"
	"point/of/sale/db"
	models "point/of/sale/models/OrderModel"
	util "point/of/sale/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func ListOrder(c *fiber.Ctx) error {
	limit := c.Query("limit")
	skip := c.Query("skip")

	intLimit, _ := strconv.Atoi(limit)
	intSkip, _ := strconv.Atoi(skip)

	var order []models.Order

	res := db.DB.Select([]string{"*"}).Limit(intLimit).Offset(intSkip).Find(&order)
	fmt.Println(order)

	type OrderList struct {
		OrderId        int                 `json:"orderId"`
		CashierID      int                 `json:"cashiersId"`
		PaymentTypesId int                 `json:"paymentTypesId"`
		TotalPrice     int                 `json:"totalPrice"`
		TotalPaid      int                 `json:"totalPaid"`
		TotalReturn    int                 `json:"totalReturn"`
		ReceiptId      string              `json:"receiptId"`
		CreatedAt      time.Time           `json:"createdAt"`
		Payments       models.PaymentTypes `json:"payment_type"`
		Cashiers       models.Cashier      `json:"cashier"`
	}
	OrderResponse := make([]*OrderList, 0)

	for _, v := range order {
		cashier := models.Cashier{}
		res := db.DB.Where("cashier_id = ?", v.CashierID).Find(&cashier)
		fmt.Println(res)

		paymentType := models.PaymentTypes{}

		res2 := db.DB.Where("id = ?", v.PaymentTypesId).Find(&paymentType)
		fmt.Println(res2)

		fmt.Println(v)

		OrderResponse = append(OrderResponse, &OrderList{
			OrderId:        v.Id,
			CashierID:      v.CashierID,
			PaymentTypesId: v.PaymentTypesId,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptId:      v.ReceiptId,
			CreatedAt:      v.CreatedAt,
			Payments:       paymentType,
			Cashiers:       cashier,
		})

	}

	return c.Status(404).JSON(fiber.Map{
		"success": true,
		"Message": "Sucess",
		"data":    OrderResponse,
		"meta": map[string]interface{}{
			"total": res.RowsAffected,
			"limit": limit,
			"skip":  skip,
		},
	})
}

func DetailOrder(c *fiber.Ctx) error {
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

	param := c.Params("orderId")

	var order models.Order
	db.DB.Where("id=?", param).First(&order)
	// db.DB.Preload("orders").First(&order, "id = ?", param)

	fmt.Println(order)
	if order.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
			"error":   map[string]interface{}{},
		})
	}
	productIds := strings.Split(order.ProductId, ",")
	TotalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := models.Product{}
		fmt.Println(productIds[i])
		res := db.DB.Where("id = ?", productIds[i]).Find(&prods)
		fmt.Println(res)
		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := models.Cashier{}
	res := db.DB.Where("cashier_id = ?", order.CashierID).Find(&cashier)
	fmt.Println(res)

	paymentType := models.PaymentTypes{}

	res2 := db.DB.Where("id = ?", order.PaymentTypesId).Find(&paymentType)
	fmt.Println(res2)

	orderTable := models.Order{}

	res3 := db.DB.Where("id = ?", order.Id).Find(&orderTable)
	fmt.Println(res3)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": map[string]interface{}{
			"order": map[string]interface{}{
				"orderId":        order.Id,
				"cashiersId":     order.CashierID,
				"paymentTypesId": order.PaymentTypesId,
				"totalPrice":     order.TotalPrice,
				"totalPaid":      order.TotalPaid,
				"totalReturn":    order.TotalReturn,
				"receiptId":      order.ReceiptId,
				"createdAt":      order.CreatedAt,
				"cashier":        cashier,
				"payment_type":   paymentType,
			},
			"products": TotalProducts,
		},
		"Message": "Success",
	})

}

func SubtotalOrder(c *fiber.Ctx) error {
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

	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		Products []products `json:"products"`
	}{}

	fmt.Println(body.Products)

	if err := c.BodyParser(&body.Products); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"Message": "Empty Body",
		})
	}

	Prodresponse := make([]*models.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	for _, v := range body.Products {
		totalPrice := 0

		prods := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductId).First(&prods)
		res := db.DB.Where("id = ?", prods.DiscountId).Find(&discount)
		fmt.Println(res)

		disc := 0
		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			percentage := totalPrice * discount.Result / 100
			disc = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc
		}
		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			fmt.Println(totalPrice)
			disc = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc

		}

		Prodresponse = append(Prodresponse,
			&models.ProductResponseOrder{
				ProductId:        prods.Id,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  disc,
			},
		)

	}
	return c.Status(200).JSON(fiber.Map{

		"message": "success",
		"success": true,
		"data": map[string]interface{}{
			"Subtotal": TotalInvoicePrice.ttprice,
			"products": Prodresponse,
		},
	})

}

func AddOrder(c *fiber.Ctx) error {
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

	//add Authorization Bearer Token...

	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		PaymentId int        `json:"paymentId"`
		TotalPaid int        `json:"totalPaid"`
		Products  []products `json:"products"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"Message": "Empty Body",
		})
	}

	Prodresponse := make([]*models.ProductResponseOrder, 0)

	var TotalInvoicePrice = struct {
		ttprice int
	}{}

	productsIds := ""
	quantities := ""
	for _, v := range body.Products {
		totalPrice := 0

		prods := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductId).First(&prods)
		res := db.DB.Where("id = ?", prods.DiscountId).Find(&discount)
		fmt.Println(res)

		productsIds = productsIds + "," + strconv.Itoa(v.ProductId)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		disc := 0
		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			percentage := totalPrice * discount.Result / 100
			disc = totalPrice - percentage
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc
		}
		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity //5000*3=15000
			fmt.Println(totalPrice)
			disc = totalPrice - discount.Result
			TotalInvoicePrice.ttprice = TotalInvoicePrice.ttprice + disc

		}

		Prodresponse = append(Prodresponse,
			&models.ProductResponseOrder{
				ProductId:        prods.Id,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  disc,
			},
		)

	}
	orderResp := models.Order{
		CashierID:      13,
		PaymentTypesId: body.PaymentId,
		TotalPrice:     TotalInvoicePrice.ttprice,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - TotalInvoicePrice.ttprice,
		ReceiptId:      "SS001",
		UpdatedAt:      time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
		ProductId:      productsIds,
		IsDownloaded:   0,
		Quantities:     quantities,
	}
	db.DB.Create(&orderResp)

	return c.Status(200).JSON(fiber.Map{

		"message": "success",
		"success": true,
		"data": map[string]interface{}{
			"order":    orderResp,
			"products": Prodresponse,
		},
	})

}

func DownloadOrder(c *fiber.Ctx) error {
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

	param := c.Params("orderId")

	var order models.Order
	db.DB.Where("id=?", param).First(&order)
	// db.DB.Preload("orders").First(&order, "id = ?", param)

	fmt.Println(order)
	if order.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
			"error":   map[string]interface{}{},
		})
	}
	productIds := strings.Split(order.ProductId, ",")

	TotalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := models.Product{}
		fmt.Println(productIds[i])
		db.DB.Where("id = ?", productIds[i]).Find(&prods)
		fmt.Println(prods.ProductId, "--------")

		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := models.Cashier{}
	res := db.DB.Where("cashier_id = ?", order.CashierID).Find(&cashier)
	fmt.Println(res)

	paymentType := models.PaymentTypes{}

	db.DB.Where("id = ?", order.PaymentTypesId).Find(&paymentType)

	orderTable := models.Order{}

	db.DB.Where("id = ?", order.Id).Find(&orderTable)
	quantities := strings.Split(order.Quantities, ",")
	quantities = quantities[1:]

	///pdf thing
	darray := [][]string{{}}
	for i := 0; i < len(TotalProducts); i++ {

		singlearray := []string{}

		singlearray = append(singlearray, TotalProducts[i].Sku)
		singlearray = append(singlearray, TotalProducts[i].Name)
		singlearray = append(singlearray, quantities[i])
		singlearray = append(singlearray, strconv.Itoa(TotalProducts[i].Price))
		darray = append(darray, singlearray)

	}

	begin := time.Now()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := []string{"Product SKU", "Product Items", "Quantities", "Price"}
	contents := darray

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	//Top Heading
	m.SetBackgroundColor(grayColor)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Order Invoice #"+strconv.Itoa(order.Id), props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.SetBackgroundColor(whiteColor)

	//Table setting
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
	//Total price
	m.Row(20, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})
	m.Row(21, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Paid:", props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPaid), props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	m.Row(22, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Return", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalReturn), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	//Invoice creation
	err := m.OutputFileAndClose("order.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	//pdf thing
	db.DB.Table("orders").Where("id=?", order.Id).Update("is_downloaded", 1)
	return c.Status(200).JSON(fiber.Map{
		"Success": true,
		"Message": "Success",
	})

}

func CheckOrderDownload(c *fiber.Ctx) error {
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

	param := c.Params("orderId")

	var order models.Order
	db.DB.Where("id=?", param).First(&order)
	if order.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
			"error":   map[string]interface{}{},
		})
	}

	fmt.Println(order)
	if order.IsDownloaded == 1 {
		return c.Status(200).JSON(fiber.Map{
			"Success": true,
			"Message": "Success",
			"data": map[string]interface{}{
				"isDownloaded": true,
			},
		})
	}

	if order.IsDownloaded == 0 {
		return c.Status(200).JSON(fiber.Map{
			"Success": true,
			"Message": "Success",
			"data": map[string]interface{}{
				"isDownloaded": false,
			},
		})
	}
	return nil
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
