package routes

import (
	CashierController "point/of/sale/controllers"

	"github.com/gofiber/fiber/v2"
)

func Routings(app *fiber.App) {

	app.Post("/cashiers", CashierController.CreateCashiar)
	app.Get("/cashiers/:cashierId", CashierController.CashierDetails)
	app.Put("/cashiers/:cashierId", CashierController.UpdateCashier)
	app.Get("/cashiers", CashierController.GetCashiers)
	app.Delete("/cashiers/:cashierId", CashierController.DeleteCashier)

	app.Get("/cashiers/:cashierId/passcode", CashierController.GetPasscode)
	app.Post("/cashiers/:cashierId/login", CashierController.VerifyPasscode)
	app.Post("cashiers/:cashierId/logout", CashierController.VerifyLogoutPasscode)

	app.Get("/categories", CashierController.ListCategory)
	app.Get("/categories/:categoryId", CashierController.DetailCategory)
	app.Post("/categories", CashierController.CreateCategory)
	app.Put("/categories/:categoryId", CashierController.UpdateCategory)
	app.Delete("/categories/:categoryId", CashierController.DeleteCategory)

	app.Get("/products", CashierController.GetProducts)
	app.Get("/products/:productId", CashierController.GetProductsById)
	app.Post("/products", CashierController.AddProducts)
	app.Put("/products/:productsId", CashierController.UpdateProduct)
	app.Delete("/products/:productId", CashierController.DeleteProduct)

	app.Get("/payments", CashierController.ListPayment)
	app.Get("/payments/:paymentId", CashierController.DetailPayment)
	app.Post("/payments", CashierController.CreatePayment)
	app.Put("/payments/:paymentId", CashierController.UpdatePayment)
	app.Delete("/payments/:paymentId", CashierController.DeletePayment)

	app.Get("/orders", CashierController.ListOrder)
	app.Get("/orders/:orderId", CashierController.DetailOrder)
	app.Post("/orders", CashierController.AddOrder)
	app.Post("/orders/subtotal", CashierController.SubtotalOrder)
	app.Get("/orders/:orderId/download", CashierController.DownloadOrder)
	app.Get("orders/:orderId/check-download", CashierController.CheckOrderDownload)

	app.Get("/revenues", CashierController.GetRevenues)
	app.Get("/solds", CashierController.GetSolds)

}
