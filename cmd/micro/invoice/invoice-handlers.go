package main

import (
	"fmt"
	"github.com/go-pdf/fpdf"
	"net/http"
	"time"
)

type Order struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	Amount    int       `json:"amount"`
	Product   string    `json:"product"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (app *application) CreateAndSendInvoice(w http.ResponseWriter, r *http.Request) {
	// receive json
	var order Order

	err := app.readJSON(w, r, &order)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// generate pdf invoice
	err = app.createInvoicePDF(order)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// create mail
	attachments := []string{
		fmt.Sprintf("./invoices/%d.pdf", order.ID),
	}

	// send mail with attachment
	err = app.SendMail("info@widgets.com", order.Email, "Your invoice", "invoice", attachments, nil)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	resp.Error = false
	resp.Message = fmt.Sprintf("Invoice %d.pdf created and sent to %s", order.ID, order.Email)
	app.writeJSON(w, http.StatusCreated, resp)

}

func (app *application) createInvoicePDF(order Order) error {
	// generate pdf invoice
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 13, 10)
	pdf.SetAutoPageBreak(true, 0)

	pdf.AddPage()

	// write info
	pdf.SetY(50)
	pdf.SetX(10)
	// Font: Case-insensitive: "Courier" for fixed-width, "Helvetica" or "Arial" for sans serif, "Times" for serif, "Symbol" or "ZapfDingbats" for symbolic.
	// StyleStr can be "B" (bold), "I" (italic), "U" (underscore), "S" (strike-out) or any combination.
	pdf.SetFont("Times", "", 11)
	//CellFormat(w float64, h float64, txtStr string, borderStr string, ln int, alignStr string, fill bool, link int, linkStr string)
	pdf.CellFormat(97, 8, fmt.Sprintf("Attention: %s %s", order.FirstName, order.LastName), "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(97, 8, order.Email, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(97, 8, order.CreatedAt.Format("2006-01-01"), "", 0, "L", false, 0, "")

	pdf.SetX(58)
	pdf.SetY(93)
	pdf.CellFormat(155, 8, order.Product, "", 0, "L", false, -1, "")
	pdf.SetX(166)
	pdf.CellFormat(20, 8, fmt.Sprintf("%d", order.Quantity), "", 0, "C", false, 0, "")
	pdf.SetX(185)
	pdf.CellFormat(20, 8, fmt.Sprintf("$%.2f", float32(order.Amount/100.0)), "", 0, "R", false, 0, "")

	invoicePath := fmt.Sprintf("./invoices/%d.pdf", order.ID)
	err := pdf.OutputFileAndClose(invoicePath)
	if err != nil {
		return err
	}
	return nil
}
