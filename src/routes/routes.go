package routes

import (
	"pg_bridge_go/controllers"
	"pg_bridge_go/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// SetupRouter sets up the Fiber router.
func SetupRouter() *fiber.App {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(middleware.CORSMiddleware())

	v1 := app.Group("/v1")

	v1.Get("/ping", controllers.Ping)

	admin := v1.Group("/admin", middleware.BasicAuthMiddlewareAdmin())
	admin.Get("/ping", controllers.Ping)
	admin.Post("/register", controllers.RegisterHandler)

	cb := v1.Group("/callback/:vendorcode")
	cb.Get("/payment", controllers.PaymentCallback)
	cb.Post("/notification", controllers.HandlePostNotificationFromPG)

	pg := v1.Group("/pg", middleware.BasicAuthMiddleware())
	pg.Get("/ping", controllers.Ping)
	pg.Post("/create-pg-vendor", controllers.CreatePaymentGatewayCredential)
	pg.Get("/get-pg-vendor/:code", controllers.GetPaymentGatewayCredential)
	pg.Get("/get-all-pg-vendor", controllers.GetAllPaymentGatewayCredential)
	pg.Put("/update-pg-vendor/:code", controllers.UpdatePaymentGatewayCredential)
	pg.Delete("/delete-pg-vendor/:code", controllers.DeletePaymentGatewayCredential)

	pgVendor := pg.Group("/vendor/:vendorcode")
	pgVendor.Post("/create-payment-request", controllers.HandleCreatePayment)
	pgVendor.Get("/get-payment-status", controllers.HandleGetPaymentStatus)

	return app
}
