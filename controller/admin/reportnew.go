package controller

import (
	"fmt"
	"net/http"
	"pkart/database"
	"pkart/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
)

func SalesReport(c *gin.Context) {
	var sales []model.Orders
	var totalamount int
	database.DB.Find(&sales)
	for _, val := range sales {
		totalamount += val.Amount
	}
	var salesItems []model.OrderItem
	var cancelCount int
	var totalSales int
	database.DB.Find(&salesItems)
	for _, val := range salesItems {
		if val.Status == "cancelled" {
			cancelCount++
		} else {
			totalSales++
		}
	}
	c.JSON(200, gin.H{
		"TotalSalesAmount": totalamount,
		"TotalSalesCount":  totalSales,
		"TotalOrderCancel": cancelCount,
	})
}

func GetReportDataa(c *gin.Context) {
	var orders []model.OrderItem
	var today time.Time
	var sales, salesreturn, totalAmount, totalDiscount int

	// Retrieve the filter query parameter
	filter := c.Query("filter")

	// Set the filter date based on the query parameter
	switch filter {
	case "Today":
		today = time.Now().Truncate(24 * time.Hour)
	case "This week":
		today = time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	case "This month":
		today = time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	default:
		today = time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour) // Default to last 30 days
	}

	// Fetch orders from the database with necessary associations and sorting
	if err := database.DB.Preload("Order").Preload("Order.User").Preload("Product").Find(&orders).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"Status":  "Error!",
			"Code":    http.StatusNotFound,
			"Error":   err.Error(),
			"Message": "No orders found!",
			"Data":    gin.H{},
		})
		return
	}

	// PDF generation settings
	marginX := 10.0
	marginY := 20.0
	lineHeight := 10.0
	const columnCount = 10

	pdf := fpdf.New("P", "mm", "A3", "")
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 25)
	pdf.CellFormat(0, 0, "SALES REPORT", "1", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 0, today.Format("2006-01-02")+" to "+time.Now().Format("2006-01-02"), "1", 0, "R", false, 0, "")
	pdf.Ln(5)

	// Table header
	headers := [columnCount]string{"Order Id", "Product Name", "Price/Unit", "Quantity", "Poffer", "Total", "Cou-Disc", "Amount", "Date", "Status"}
	colWidths := [columnCount]float64{20.0, 80.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	for _, header := range headers {
		pdf.CellFormat(colWidths[0], lineHeight, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)

	// Process orders and fill table data
	for _, order := range orders {
		if order.Order.CreatedAt.After(today) {
			totalAmount += int(order.Amount)
			totalDiscount += int(order.SubTotal) - int(order.Amount)
			if order.Status == "pending" || order.Status == "delivered" || order.Status == "shipped" {
				sales++
			} else if order.Status == "cancelled" || order.Status == "returned" {
				salesreturn++
			}

			amount := order.Quantity * (order.Product.Price - order.Product.Offer)
			pdf.CellFormat(colWidths[0], lineHeight, fmt.Sprintf("%d", order.OrderId), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[1], lineHeight, order.Product.Name, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[2], lineHeight, strconv.Itoa(order.Product.Price), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[3], lineHeight, strconv.Itoa(order.Quantity), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[4], lineHeight, strconv.Itoa(order.Product.Offer*order.Quantity), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[5], lineHeight, strconv.Itoa(amount), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[6], lineHeight, strconv.Itoa(int(order.SubTotal)-int(order.Amount)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[7], lineHeight, strconv.Itoa(int(order.Amount)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[8], lineHeight, order.Order.CreatedAt.Format("2006-01-02"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidths[9], lineHeight, order.Status, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
		}
	}

	// Summary section
	pdf.SetFont("Arial", "B", 10)
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprintf("Sales: %d", sales), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprintf("Sales return: %d", salesreturn), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprintf("Total amount: %d", totalAmount), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprintf("Total discount: %d", totalDiscount), "1", 0, "R", false, 0, "")

	// Generate PDF file
	path := fmt.Sprintf("/home/prajith/Desktop/Bttp/salesReport_%s_%s.pdf", time.Now().Format("20060102_150405"), filter)
	if err := pdf.OutputFileAndClose(path); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Code":    http.StatusUnauthorized,
			"Message": "Failed to generate PDF file",
			"Status":  "Error!",
			"Error":   err.Error(),
		})
		return
	}

	// Set response headers and send the PDF file
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	c.JSON(http.StatusOK, gin.H{
		"Status":  "Success!",
		"Code":    http.StatusOK,
		"Message": "PDF downloaded successfully!",
		"Data":    gin.H{},
	})
}

func GetReportData(c *gin.Context) {
	var orders []model.OrderItem
	var today time.Time
	var sales, salesreturn, totalAmount, totalDiscount int

	Filter := c.Query("filter")

	switch Filter {
	case "Today":
		today = time.Now().Truncate(24 * time.Hour)
	case "This week":
		today = time.Now().AddDate(0, 0, -7).Truncate(24 * time.Hour)
	case "This month":
		today = time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	default:
		today = time.Now().AddDate(0, 0, -30).Truncate(24 * time.Hour) // Default to last 30 days
	}

	if err := database.DB.Preload("Order").Preload("Order.User").Preload("Product").Find(&orders).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No orders found!",
			"Data":    gin.H{},
		})
		return
	}

	marginX := 10.0
	marginY := 20.0
	lineHt := 10.0
	const colNumber = 11

	pdf := fpdf.New("P", "mm", "A3", "")
	pdf.SetMargins(marginX, marginY, marginX)
	pdf.AddPage()

	// Title
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 25)
	pdf.CellFormat(0, 0, "SALES REPORT", "1", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 0, today.Format("2006-01-02")+" to "+time.Now().Format("2006-01-02"), "1", 0, "R", false, 0, "")
	pdf.Ln(5)

	// Table header
	header := [colNumber]string{"Sl No", "Order Id", "Product Name", "Price/Unit", "Quantity", "Poffer", "Total", "Cou-Disc", "Amount", "Date", "Status"}
	colWidth := [colNumber]float64{15.0, 20.0, 80.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)
	serialNumber := 1

	// Process orders and fill table data
	for _, order := range orders {
		if order.Status == "pending" || order.Status == "delivered" || order.Status == "shipped" {
			totalAmount += int(order.Amount)
			totalDiscount += int(order.SubTotal) - int(order.Amount)
		}
	}

	for _, v := range orders {
		if v.Status == "pending" || v.Status == "delivered" || v.Status == "shipped" {
			sales++
		} else if v.Status == "cancelled" || v.Status == "returned" {
			salesreturn++
		}
		if v.Order.CreatedAt.After(today) {
			totamount := v.Quantity * (v.Product.Price - v.Product.Offer)
			pdf.CellFormat(colWidth[0], lineHt, strconv.Itoa(serialNumber), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[1], lineHt, fmt.Sprintf("%d", v.OrderId), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[2], lineHt, v.Product.Name, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(v.Product.Price), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[4], lineHt, strconv.Itoa(v.Quantity), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[5], lineHt, strconv.Itoa(v.Product.Offer*v.Quantity), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[6], lineHt, strconv.Itoa(totamount), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[7], lineHt, strconv.Itoa(int(v.SubTotal)-int(v.Amount)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[8], lineHt, strconv.Itoa(int(v.Amount)), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[9], lineHt, v.Order.CreatedAt.String()[:10], "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[10], lineHt, v.Status, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)

			// Increment the serial number
			serialNumber++
		}
	}

	pdf.SetFont("Arial", "B", 10)
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprint("Sales : ", sales), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprint("Sales return : ", salesreturn), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprint("Total amount : ", totalAmount), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprint("Total discount : ", totalDiscount), "1", 0, "R", false, 0, "")

	// Generate PDF file

	path := fmt.Sprintf("/home/ubuntu/pkart_1/salesReport_%s_%s.pdf", time.Now().Format("20060102_150405"), Filter)
	if err := pdf.OutputFileAndClose(path); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Status":  "error",
			"Code":    401,
			"Message": "failed to generate PDF file",
			"Error":   err.Error(),
			"Data": gin.H{},
		})
		return
	}

	// Set response headers and send the PDF file
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"Message": "PDF downloaded successfully",
		"Data":    gin.H{},
	})
}
