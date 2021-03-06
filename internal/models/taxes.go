package models

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"strings"
)

// TaxOrigin is defines tax scope according to inventory item origin
type TaxOrigin string

const (
	TaxOriginAll    TaxOrigin = "ALL"
	TaxOriginLocal  TaxOrigin = "LOCAL"
	TaxOriginImport TaxOrigin = "IMPORT"
)

// TaxCondition is defines category given in tax subject or exempt from tax.
type TaxCondition string

const (
	UnknownTC    TaxCondition = "UNKNOWN"
	ExemptToTax  TaxCondition = "EXEMPT"  // refers to only tax types in context will be free from tax
	SubjectToTax TaxCondition = "SUBJECT" // refers to only tax types in context will be effected from tax
)

// Tax
type Tax struct {
	Id         uuid.UUID          `json:"id"`
	Name       string             `json:"name"`
	Rate       decimal.Decimal    `json:"rate"`
	Origin     TaxOrigin          `json:"origin"`
	Condition  TaxCondition       `json:"condition"`
	Categories map[uuid.UUID]bool `json:"categories"`
}

type SaleItem struct {
	*InventoryItem

	Taxes decimal.Decimal `json:"taxes"`
	Gross decimal.Decimal `json:"gross"`
}

func (t *Tax) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(b)
}

func (si *SaleItem) String() string {
	b, err := json.Marshal(si)
	if err != nil {
		return ""
	}
	return string(b)
}

func (tt *TaxOrigin) UnmarshalText(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch str {
	case "LOCAL", "IMPORT", "ALL":
		*tt = TaxOrigin(str)

	default:
		*tt = TaxOriginAll
	}

	return nil
}

func (tc *TaxCondition) UnmarshalText(b []byte) error {
	str := strings.Trim(string(b), `"`)

	switch str {
	case "EXEMPT", "SUBJECT":
		*tc = TaxCondition(str)

	default:
		*tc = UnknownTC
	}

	return nil
}
