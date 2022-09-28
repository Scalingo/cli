package scalingo

import (
	"context"
	"strings"
	"time"

	"gopkg.in/errgo.v1"
)

type InvoicesService interface {
	InvoicesList(context.Context, PaginationOpts) (Invoices, PaginationMeta, error)
	InvoiceShow(context.Context, string) (*Invoice, error)
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
		PaginationMeta PaginationMeta `json:"pagination"`
	}
}

type InvoiceRes struct {
	Invoice *Invoice `json:"invoice"`
}

func (c *Client) InvoicesList(ctx context.Context, opts PaginationOpts) (Invoices, PaginationMeta, error) {
	var invoicesRes InvoicesRes
	err := c.ScalingoAPI().ResourceList(ctx, "account/invoices", opts.ToMap(), &invoicesRes)
	if err != nil {
		return nil, PaginationMeta{}, errgo.Mask(err)
	}
	return invoicesRes.Invoices, invoicesRes.Meta.PaginationMeta, nil
}

func (c *Client) InvoiceShow(ctx context.Context, id string) (*Invoice, error) {
	var invoiceRes InvoiceRes
	err := c.ScalingoAPI().ResourceGet(ctx, "account/invoices", id, nil, &invoiceRes)
	if err != nil {
		return nil, errgo.Mask(err)
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
