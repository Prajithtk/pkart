package controller

import (
	"fmt"
	"pkart/database"
	"pkart/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func Invoice(c *gin.Context, OrderId int) error {

	var payment model.Payment
	var Order []model.OrderItem
	var address model.Address

	if err := database.DB.Preload("Prdct").Preload("Order").Preload("Order").Find(&Order, "Order_Id=?", OrderId).Error; err != nil {
		return err
	}

	if err := database.DB.Preload("Order").Preload("User").First(&payment, "Order_Id=?", OrderId).Error; err != nil {
		return err
	}

	if err := database.DB.First(&address, "Id=?", payment.Order.AddressId).Error; err != nil {
		return err
	}

	marginX := 10.0
	marginY := 20.0

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetMargins(marginX, marginY, marginX)

	pdf.AddPage()

	// pdf.ImageOptions("assets/logo.gif", 5, 0, 75, 25, false, gofpdf.ImageOptions{ImageType: "GIF", ReadDpi: true}, 0, "")

	pdf.AddUTF8Font("PlayfairDisplayMedium-9YKZK", "", "./assets/fonts/PlayfairDisplayMedium-9YKZK.ttf")
	pdf.SetFont("PlayfairDisplayMedium-9YKZK", "", 20)

	pdf.Ln(5)

	pdf.CellFormat(0, 0, "INVOICE", "0", 0, "C", false, 0, "")

	pdf.SetFont("Arial", "", 13)

	pdf.Ln(5)
	pdf.Cell(0, 0, "")
	pdf.Ln(5)
	pdf.Cell(0, 0, "Dear "+payment.User.Name+",")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 16)

	_, lineHeight := pdf.GetFontSize()
	currentY := pdf.GetY() + lineHeight

	pdf.SetY(currentY)

	pdf.Cell(150, 10, "PKART LTD")
	pdf.Cell(0, 10, "To")

	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 11)

	pdf.Cell(150, 5, "Kakkanchery")
	pdf.Cell(0, 5, address.BuildingName)

	pdf.Ln(-1)

	pdf.Cell(150, 5, "673620")
	pdf.Cell(0, 5, address.PinCode)

	pdf.Ln(-1)

	pdf.Cell(150, 5, "Kinfra")
	pdf.Cell(0, 5, address.City)

	pdf.Ln(-1)

	pdf.Cell(150, 5, "Kerala, India")
	pdf.Cell(0, 5, "India")

	pdf.Ln(5)
	pdf.Cell(0, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 16)

	pdf.Cell(0, 0, "Ref#"+strconv.Itoa(int(payment.OrderId)))
	pdf.Ln(3)

	lineHt := 10.0
	const colNumber = 5
	header := [colNumber]string{"No", "Product Name", "Quantity", "Price", "Total"}
	colWidth := [colNumber]float64{10.0, 55.0, 25.0, 30.0, 30.0}

	// Headers
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	for colJ := 0; colJ < colNumber; colJ++ {
		pdf.CellFormat(colWidth[colJ], lineHt, header[colJ], "1", 0, "C", true, 0, "")
	}

	pdf.Ln(-1)

	// Table data
	pdf.SetFont("Arial", "", 11)

	for i := 0; i < len(Order); i++ {
		// Column 1: Unit
		// Column 2: Description
		// Column 3: Price per unit
		amount := Order[i].Quantity * Order[i].Product.Offer

		pdf.CellFormat(colWidth[0], lineHt, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[1], lineHt, Order[i].Product.Name, "1", 0, "LM", false, 0, "")
		pdf.CellFormat(colWidth[2], lineHt, fmt.Sprintf("%d", Order[i].Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%d", Order[i].Product.Offer), "1", 0, "C", false, 0, "")
		pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%d", amount), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Calculate the subtotal
	pdf.SetFontStyle("B")
	leftIndent := 0.0
	for i := 0; i < 3; i++ {
		leftIndent += colWidth[i]
	}
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Subtotal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidth[4], lineHt, fmt.Sprintf("%d", Order[0].Order.Amount), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	DeliveryCharge := 0
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Delivery Charge", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%d", DeliveryCharge), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	grandTotal := float64(Order[0].Order.Amount) + float64(DeliveryCharge)
	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Grand total", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidth[3], lineHt, fmt.Sprintf("%.2f", grandTotal), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetX(marginX + leftIndent)
	pdf.CellFormat(colWidth[3], lineHt, "Payment status", "1", 0, "C", false, 0, "")
	pdf.CellFormat(colWidth[3], lineHt, payment.Status, "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	path := "/home/ubuntu/My-Ecommerce/invoice" + strconv.Itoa(int(payment.OrderId)) + ".pdf"
	if err := pdf.OutputFileAndClose(path); err != nil {
		return err
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))
	c.Writer.Header().Set("Content-Type", "application/pdf")
	c.File(path)

	return nil
}
