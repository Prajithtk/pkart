package controller

import (
	"fmt"
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

func GetReportData(c *gin.Context) {

	var orders []model.OrderItem
	var today time.Time
	var sales, salesreturn = 0, 0

	Filter := c.Query("filter")
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

	switch Filter {
	case "Today":
		today = time.Now().Truncate(24 * time.Hour)
	case "This week":
		today = time.Now().Truncate(168 * time.Hour)
	case "This month":
		today = time.Now().Truncate(730 * time.Hour)
	default:
		today = time.Now()
	}
	marginX := 10.0
	marginY := 20.0

	lineHt := 10.0
	const colNumber = 7

	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.SetMargins(marginX, marginY, marginX)
	pdf.AddPage()

	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 25)
	pdf.CellFormat(0, 0, "SALES REPORT", "1", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 0, today.String()[:10]+" to "+time.Now().String()[:10], "1", 0, "R", false, 0, "")
	pdf.Ln(5)

	header := [colNumber]string{"Order Id", "Product Name", "Price/Unit", "Quantity", "Amount", "Date", "Status"}
	colWidth := [colNumber]float64{20.0, 60.0, 25.0, 20.0, 20.0, 25.0, 20.0}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)
	pdf.SetFont("Arial", "", 9)

	for _, v := range orders {
		if v.Status == "pending" || v.Status == "delivered" ||v.Status =="shipped" {
			sales++
		} else if v.Status == "cancelled" || v.Status == "returned" {
			salesreturn++
		}
		if today == time.Now() {
			amount := v.Quantity * (v.Product.Price - v.Product.Offer)
			pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", v.OrderId), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[1], lineHt, v.Product.Name, "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[2], lineHt, strconv.Itoa(v.Product.Price-v.Product.Offer), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(v.Quantity), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[4], lineHt, strconv.Itoa(amount), "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[5], lineHt, v.Order.CreatedAt.String()[:10], "1", 0, "C", false, 0, "")
			pdf.CellFormat(colWidth[6], lineHt, v.Status, "1", 0, "C", false, 0, "")
			pdf.Ln(-1)

		} else {
			if time.Now().After(today) {
				amount := v.Quantity * (v.Product.Price-v.Product.Offer)
				pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", v.OrderId), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[1], lineHt, v.Product.Name, "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[2], lineHt, strconv.Itoa(v.Product.Price-v.Product.Offer), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[3], lineHt, strconv.Itoa(v.Quantity), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[4], lineHt, strconv.Itoa(amount), "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[5], lineHt, v.Order.CreatedAt.String()[:10], "1", 0, "C", false, 0, "")
				pdf.CellFormat(colWidth[6], lineHt, v.Status, "1", 0, "C", false, 0, "")
				pdf.Ln(-1)
			}
		}
	}
	pdf.SetFont("Arial", "B", 10)
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprint("Sales : ", sales), "1", 0, "R", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(0, 0, fmt.Sprint("Sales return : ", salesreturn), "1", 0, "R", false, 0, "")

	path := "/home/prajith/Desktop/Bttp/salesReport_" + time.Now().String()[:10] + "_" + Filter + ".pdf"
	if err := pdf.OutputFileAndClose(path); err != nil {
		c.JSON(401, gin.H{
			"Code":    401,
			"Message": "Failed to generate PDF file",
			"Status":  "Error!",
			"Error":   err.Error(),
		})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Pdf downloaded successfully!",
		"Data":    gin.H{},
	})
}
