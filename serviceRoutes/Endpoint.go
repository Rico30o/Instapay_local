package serviceroutes

import (
	traceshandler "instapay/TracesHandler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Post("/Trace", traceshandler.TxnID)           //Txn Id
	app.Post("/Txn-Id", traceshandler.NetworkAlertID) //Network Alert Id

	// app.Get("/alerts", traceshandler.AlertAccounts) //AlertAccounts
	// Trace.Get("/alert/:account", traceshandler.AlertAccounts)
}
