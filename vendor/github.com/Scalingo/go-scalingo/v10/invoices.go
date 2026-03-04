package scalingo

import (
	"context"
	"strings"
	"time"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/pagination"
)

type InvoicesService interface {
	InvoicesList(ctx context.Context, paginationReq pagination.Request) (Invoices, pagination.Meta, error)
	InvoiceShow(ctx context.Context, id string) (*Invoice, error)
}

var _ InvoicesService = (*Client)(nil)

const BillingMonthDateFormat = "2006-01-02"

type billingMonthDate time.Time

type InvoiceItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Price int    `json:"price"`
}

type InvoiceDetailedItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Price int    `json:"price"`
	App   string `json:"app"`
}

type Invoice struct {
	ID                string                `json:"id"`
	TotalPrice        int                   `json:"total_price"`
	TotalPriceWithVat int                   `json:"total_price_with_vat"`
	BillingMonth      billingMonthDate      `json:"billing_month"`
	PdfURL            string                `json:"pdf_url"`
	InvoiceNumber     string                `json:"invoice_number"`
	State             string                `json:"state"`
	VatRate           int                   `json:"vat_rate"`
	Items             []InvoiceItem         `json:"items"`
	DetailedItems     []InvoiceDetailedItem `json:"detailed_items"`
}

type Invoices []*Invoice

type InvoicesRes struct {
	Invoices Invoices `json:"invoices"`
	Meta     struct {
		Pagination pagination.Meta `json:"pagination"`
	}
}

type InvoiceRes struct {
	Invoice *Invoice `json:"invoice"`
}

func (c *Client) InvoicesList(ctx context.Context, paginationReq pagination.Request) (Invoices, pagination.Meta, error) {
	var invoicesRes InvoicesRes
	err := c.ScalingoAPI().ResourceList(ctx, "account/invoices", paginationReq.ToURLValues(), &invoicesRes)
	if err != nil {
		return nil, pagination.Meta{}, errors.Wrap(ctx, err, "list invoices")
	}
	return invoicesRes.Invoices, invoicesRes.Meta.Pagination, nil
}

func (c *Client) InvoiceShow(ctx context.Context, id string) (*Invoice, error) {
	var invoiceRes InvoiceRes
	err := c.ScalingoAPI().ResourceGet(ctx, "account/invoices", id, nil, &invoiceRes)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "show invoice")
	}
	return invoiceRes.Invoice, nil
}

func (b *billingMonthDate) UnmarshalJSON(data []byte) error {
	value := strings.Trim(string(data), `"`)
	if value == "" || value == "null" {
		return nil
	}
	t, err := time.Parse(BillingMonthDateFormat, value)
	if err != nil {
		return err
	}
	*b = billingMonthDate(t)
	return nil
}
