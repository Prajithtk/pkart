package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func CreateInvoice(c *gin.Context) {
	userID := c.GetUint("userid")
	orderId := c.Param("ID")
	var user model.Users
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "failed",
			"Code":    404,
			"Message": "user not found",
			"Data":    gin.H{},
		})
		return
	}
	var orderItem []model.OrderItem
	if err := database.DB.Where("order_id = ? AND status!=?", orderId, "cancelled").Preload("Product").Preload("Order.Address").Find(&orderItem).Error; err != nil {
		c.JSON(503, gin.H{
			"Status":  "failed",
			"Code":    503,
			"Message": "failed to fetch orders",
			"Data":    gin.H{},
		})
		return
	}
	for _, order := range orderItem {
		if order.Status != "delivered" {
			c.JSON(202, gin.H{
				"Status":  "failed",
				"Code":    202,
				"Message": "order not delivered",
				"Data":    gin.H{},
			})
			return
		}
	}
	var order model.Orders
	var Discount float64
	database.DB.First(&order, orderId)

	pdf := gofpdf.New("P", "mm", "A3", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.Ln(5)
	pdf.CellFormat(0, 0, "INVOICE", "", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(30)
	pdf.Cell(10, -32, "Invoice No: "+orderId)
	pdf.Ln(5)
	pdf.Cell(10, -32, "Invoice Date: "+order.CreatedAt.Format("2006-01-02"))
	pdf.Ln(15)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(10, -32, "Bill To: ")
	pdf.Ln(5)
	pdf.Cell(10, -32, "Customer Name: "+user.Name)
	pdf.SetFont("Arial", "", 12)
	pdf.Ln(5)
	for _, val := range orderItem {
		pdf.Cell(10, -32, "Address: "+val.Order.Address.City+", "+val.Order.Address.State)
		pdf.Ln(5)
		pdf.Cell(10, -32, (val.Order.Address.PinCode))
		pdf.Ln(5)
		pdf.Cell(10, -32, "Phone no : "+(user.Phone))
		pdf.Ln(5)
		pdf.SetFont("Arial", "", 12)
		pdf.Ln(10)
		break
	}

	pdf.SetXY(10, 20)
	pdf.CellFormat(250, 30, "pkart", "", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(12, 40, "Crystal Plaza , Calicut road", "", 0, "R", false, 0, "")
	pdf.CellFormat(12, 50, "15th floor ,Ph: +95 365452", "", 0, "R", false, 0, "")
	pdf.Ln(60)

	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(20, 10, "No.", "1", 0, "C", true, 0, "")
	pdf.CellFormat(75, 10, "Item Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Quantity", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Product Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Offer Price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Coup-Disc", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 10, "Net Price", "1", 0, "C", true, 0, "")

	pdf.Ln(10)

	totalAmount := 0.0
	for i, order := range orderItem {
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(75, 10, order.Product.Name, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", order.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%d", order.Product.Price), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", order.SubTotal), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", order.SubTotal-order.Amount), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", order.Amount), "1", 0, "R", false, 0, "")
		pdf.Ln(10)
		totalAmount += float64(order.SubTotal)
	}
	if order.ShippingCharge > 0 {
		order.Amount -= (order.ShippingCharge)
	}
	Discount = totalAmount - float64(order.Amount)
	totalAmount -= float64(Discount)
	if Discount > 0 {
		pdf.CellFormat(235, 10, "Discount:", "1", 0, "R", true, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%2.f", Discount), "1", 0, "R", true, 0, "")
		pdf.Ln(10)
	}
	if order.ShippingCharge > 0 {
		totalAmount += float64(order.ShippingCharge)
		pdf.CellFormat(235, 10, "Shipping charge:", "1", 0, "R", true, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d", order.ShippingCharge), "1", 0, "R", true, 0, "")
		pdf.Ln(10)
	}
	Discount = 0
	pdf.CellFormat(235, 10, "Total Amount: ", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", totalAmount), "1", 0, "R", true, 0, "")

	pdfPath := "/home/ubuntu/pkart_1/Invoice.pdf"
	// pdfPath := "/home/prajith/Desktop/Bttp/Invoice.pdf"
	if err := pdf.OutputFileAndClose(pdfPath); err != nil {
		c.JSON(500, gin.H{
			"Status":  "failed",
			"Code":    500,
			"Message": "failed to generate PDF file",
			"Data":    gin.H{},
		})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", pdfPath))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(pdfPath)

	c.JSON(200, gin.H{
		"Status":  "success",
		"Code":    200,
		"message": "PDF file generated and sent successfully",
		"Data":    gin.H{},
	})
}
