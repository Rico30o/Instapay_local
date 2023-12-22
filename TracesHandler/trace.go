package traceshandler

import (
	traceroute "instapay/TraceRoute"
	"instapay/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Limit(1), 10) // Allow 10 requests per second

func TxnID(c *fiber.Ctx) error {

	// Check rate limit
	if limiter.Allow() == false {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": []fiber.Map{
					{
						"Source":      "Gateway",
						"ReasonCode":  "RATE_LIMIT_EXCEEDED",
						"Description": "You have exceeded the service rate limit. Maximum allowed: ${rate_limit.output} TPS",
						"Recoverable": true,
						"Details":     nil,
					},
				},
			},
		})
	}
	// Check if the request method is not POST
	if c.Method() != fiber.MethodPost {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": []fiber.Map{
					{
						"Source":      "TRACE_FINANCIAL_CRIME",
						"ReasonCode":  "METHOD_NOT_ALLOWED",
						"Description": "Only POST method allowed",
						"Recoverable": false,
						"Details":     nil,
					},
				},
			},
		})
	}

	// Assuming you have a valid payload.Transaction struct
	Trace := &traceroute.Transaction{}

	// Parse the request body
	if parsErr := c.BodyParser(Trace); parsErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": []fiber.Map{
					{
						"Source":      "TRACE_FINANCIAL_CRIME",
						"ReasonCode":  "BAD_REQUEST",
						"Description": "We could not handle your request",
						"Recoverable": false,
						"Details":     "The request contains a bad payload",
					},
				},
			},
		})
	}

	// Check if the request body is empty
	if Trace == nil || (Trace.TxnID == "" && Trace.Type == "") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": []fiber.Map{
					{
						"Source":      "TRACE_FINANCIAL_CRIME",
						"ReasonCode":  "BAD_REQUEST",
						"Description": "We could not handle your request",
						"Recoverable": false,
						"Details":     "The request body is empty",
					},
				},
			},
		})
	}

	expectedAlertID := "2b588f38d1bc40bf85fc91397bc98465"

	// Handle permission denied error
	if Trace.TxnID != expectedAlertID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": []fiber.Map{
					{
						"Source":      "Gateway",
						"ReasonCode":  "PERMISSION_DENIED",
						"Description": "Invalid customer for third party",
						"Recoverable": false,
						"Details":     nil,
					},
				},
			},
		})
	}

	// Construct the response body
	response := &traceroute.TraceResponse{
		ID:        "e99a18c428cb38d5f260853678922e03",
		Time:      time.Now().UTC(),
		NetworkID: expectedAlertID,
		TransactionAlerts: []traceroute.TransactionAlert{
			{

				ID:             "e99a18c428cb38d5f260853678922e03",
				TxnID:          expectedAlertID,
				NetworkAlertID: "e99a18c428cb38d5f260853678922e03",
				NetworkID:      expectedAlertID,
				Time:           time.Now().UTC(),
				TxnTime:        time.Now().UTC(),
				SourceID:       "GB98MIDL07009312345678",
				DestID:         "GB98MIDL07009312345678",
				SourceBankID:   "DEUTDEFF",
				SourceBankName: "Barclays",
				DestBankID:     "DEUTDEFF",
				DestBankName:   "Lloyds",
				Value:          10034,
			},
		},
		AccountAlerts: []traceroute.AccountAlert{
			{
				ID:             "e99a18c428cb38d5f260853678922e03",
				NetworkAlertID: "e99a18c428cb38d5f260853678922e03",
				AccountID:      "GB98MIDL07009312345678",
				NetworkID:      "2b588f38-d1bc-40bf-85fc-91397bc98465",
				OwningBankID:   "DEUTDEFF",
				OwningBankName: "OwningBank",
				Time:           time.Now().UTC(),
			},
		},
		VizURL:             "https://api.fcs.uk.mastercard.com/trace/financialcrime/viz/d41d8cd98f00b204e9800998ecf8427e",
		SourceTxnID:        "2b588f38-d1bc-40bf-85fc-91397bc98465",
		SourceTxnType:      "FRAUD",
		Length:             20,
		Generations:        3,
		TotalValue:         10034,
		SourceValue:        10034,
		UniqueAccounts:     16,
		MeanDwellTime:      "P3Y6M4DT12H30M5S",
		MedianDwellTime:    "P3Y6M4DT12H30M5S",
		MeanMuleScore:      0.845,
		ElapsedTime:        "P3Y6M4DT12H30M5S",
		NumActionedMules:   2,
		NumLegitimate:      7,
		NumNotInvestigated: 3,
		ParentAlertID:      "e99a18c428cb38d5f260853678922e03",
		DecisionDate:       time.Now().UTC(),
		MostRecentFeedback: "ACTIONED_MULE",
	}

	// Return the constructed response
	return c.JSON(response)

}

// func NetworkAlertID(c *fiber.Ctx) error {
// 	networkAlertID := c.Params("networkAlertID") // Assuming "networkAlertID" is the parameter name
// 	Trace := &traceroute.Transaction{}

// 	NetworkAlertID := &traceroute.NetworkAlertID{}
// 	if parseErr := c.BodyParser(NetworkAlertID); parseErr != nil {
// 		// Return a bad request response with JSON body
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"Errors": fiber.Map{
// 				"Error": []fiber.Map{
// 					{
// 						"Source":      "TRACE_FINANCIAL_CRIME",
// 						"ReasonCode":  "BAD_REQUEST",
// 						"Description": "We could not handle your request",
// 						"Recoverable": false,
// 						"Details":     "The request contains a bad payload",
// 					},
// 				},
// 			},
// 		})
// 	}

// 	if networkAlertID == Trace.TxnID {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"message": "Visualization has been created.",
// 		})
// 	}

// 	// Check for third-party error
// 	if thirdPartyErr := checkThirdPartyError(c, *NetworkAlertID); thirdPartyErr == nil {
// 		// Handle the third-party error
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"Errors": fiber.Map{
// 				"Error": []fiber.Map{
// 					{
// 						"Source":      "Gateway",
// 						"ReasonCode":  "PERMISSION_DENIED",
// 						"Description": "Invalid customer for third party",
// 						"Recoverable": false,
// 						"Details":     nil,
// 					},
// 				},
// 			},
// 		})
// 	}

// 	// Handle other cases if needed
// 	return nil
// }

// // Function to simulate third-party error checking
// func checkThirdPartyError(c *fiber.Ctx, alert traceroute.NetworkAlertID) error {
// 	// Assuming traceroute is imported correctly

// 	result := database.DB.Debug().Exec(`
//         INSERT INTO third_party (
//             network_alert_id, format, width, height, legend, type, colour_Mode
//         ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
// 		alert.Network_alert_id,
// 		alert.Format,
// 		alert.Width,
// 		alert.Height,
// 		alert.Legend,
// 		alert.Type,
// 		alert.Colour_Mode,
// 	)

// 	if result.Error != nil {
// 		// Handle the error if the database query fails
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": result.Error.Error(),
// 		})
// 	}

// 	// Check if the number of affected rows is zero
// 	if result.RowsAffected == 0 {
// 		// Handle the error if the database query fails
// 		return errors.New("Invalid customer for third party")
// 	}

// 	return nil
// }

// func NetworkAlertID(c *fiber.Ctx) error {
// 	networkAlertID := c.Params("network_alert_id") // Use the correct parameter name
// 	expectedAlertID := "e99a18c428cb38d5f260853678922e03"

// 	networkAlert := &traceroute.NetworkAlertID{}
// 	if parseErr := c.BodyParser(networkAlert); parseErr != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(createError("BAD_REQUEST", "The request contains a bad payload"))
// 	}

// 	if networkAlert.Network_alert_id == "" && networkAlert.Format == "" && networkAlert.Type == "" && networkAlert.Colour_Mode == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(createError("BAD_REQUEST", "The request contains a bad payload"))
// 	}

// 	if c.Method() != fiber.MethodGet {
// 		return c.Status(fiber.StatusMethodNotAllowed).JSON(createError("METHOD_NOT_ALLOWED", "Only GET method allowed"))
// 	}

// 	alert := &traceroute.NetworkAlertID{}
// 	result := database.DB.Debug().Model(&traceroute.NetworkAlertID{}).Where("network_alert_id = ?", networkAlertID).First(&alert)

// 	if result.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(createError("DATABASE_ERROR", fmt.Sprintf("Error checking Network_alert_id: %s", result.Error)))
// 	}

// 	if networkAlertID != expectedAlertID {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"message": "Visualization has been created.",
// 		})
// 	} else {
// 		return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 			"message": createError("INVALID_ALERT_ID", "Invalid network alert ID"),
// 		})
// 	}

// }

// func createError(reasonCode, details string) fiber.Map {
// 	return fiber.Map{
// 		"Errors": fiber.Map{
// 			"Error": []fiber.Map{
// 				{
// 					"Source":      "Gateway",
// 					"ReasonCode":  reasonCode,
// 					"Description": "Invalid customer for third party",
// 					"Recoverable": false,
// 					"Details":     details,
// 				},
// 			},
// 		},
// 	}

// }

// func checkIfAlertIDExistsInDatabase(alert *traceroute.NetworkAlertID) (bool, error) {
// 	// Perform a database query to check if networkAlertID exists
// 	result := database.DB.Debug().Exec(`
// 	SELECT id, network_alert_id, status
//     FROM public.third_party
// 	WHERE network_alert_id = ?`,
// 		alert.Network_alert_id,
// 	)

// 	if result.Error != nil {
// 		// Return the error if the database query fails
// 		return false, result.Error
// 	} else {

// 	}

// 	// Check if any rows are returned (indicating that the alert ID exists)
// 	return result.RowsAffected > 0, nil
// }

// func NetworkAlertID(c *fiber.Ctx) error {
// 	alert := &traceroute.NetworkAlertID{}
// 	expectedAlertID := "2b588f38d1bc40bf85fc91397bc98465"

// 	// Parse the request body
// 	if parseErr := c.BodyParser(alert); parseErr != nil {
// 		// Return a bad request response in case of parsing error
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"Errors": fiber.Map{
// 				"Error": fiber.Map{
// 					"Source":      "TRACE_FINANCIAL_CRIME",
// 					"ReasonCode":  "BAD_REQUEST",
// 					"Description": "We could not handle your request",
// 					"Recoverable": false,
// 					"Details":     "The request contains a bad payload",
// 				},
// 			},
// 		})
// 	}

// 	// Check if the network_alert_id matches the expectedAlertID
// 	if alert.Network_alert_id != expectedAlertID {
// 		// Check if the network_alert_id exists in the database
// 		status, err := checkIfAlertIDExistsInDatabase(alert)
// 		if err != nil {
// 			// Handle the error if the database query fails
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}

// 		if status {
// 			// Respond with a third-party error if the status is true
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 				"Errors": fiber.Map{
// 					"Error": fiber.Map{
// 						"Source":      "Gateway",
// 						"ReasonCode":  "PERMISSION_DENIED",
// 						"Description": "Invalid customer for third party",
// 						"Recoverable": false,
// 						"Details":     nil,
// 					},
// 				},
// 			})
// 		}

// Return a success response if the status is false
// return c.JSON(fiber.Map{
// 	"Response": "Visualization has been created",
// })
// 	}

// 	return nil
// }

func NetworkAlertID(c *fiber.Ctx) error {
	alert := &traceroute.NetworkAlertID{}
	// expectedAlertID := "2b588f38d1bc40bf85fc91397bc98465"

	// Parse the request body
	if parseErr := c.BodyParser(alert); parseErr != nil {
		// Return a bad request response in case of parsing error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": fiber.Map{
					"Source":      "TRACE_FINANCIAL_CRIME",
					"ReasonCode":  "BAD_REQUEST",
					"Description": "We could not handle your request",
					"Recoverable": false,
					"Details":     "The request contains a bad payload",
				},
			},
		})
	}

	// Call a function to retrieve status from the database
	if err := database.DB.Debug().Raw(`SELECT * FROM third_party WHERE status = ?`, alert.Status).Error; err != nil {
		// Handle database error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": fiber.Map{
					"Source":      "DATABASE",
					"ReasonCode":  "INTERNAL_SERVER_ERROR",
					"Description": "Error retrieving status from the database",
					"Recoverable": false,
					"Details":     err.Error(),
				},
			},
		})
	}
	if alert.Status == true {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Errors": fiber.Map{
				"Error": fiber.Map{
					"Source":      "THIRD_PARTY",
					"ReasonCode":  "INTERNAL_SERVER_ERROR",
					"Description": "Third-party error",
					"Recoverable": false,
					"Details":     "An error occurred while processing the request",
				},
			},
		})
	} else {
		return c.JSON(fiber.Map{
			"Response": "Visualization has been created",
		})
	}
}

// Alert
// func AlertAccounts(c *fiber.Ctx) error {
// 	// Parse query parameters
// 	since := c.Params("since")
// 	limitStr := c.Params("limit")
// 	filter := c.Params("filter")
// 	paginationToken := c.Params("pagination_token")

// }
