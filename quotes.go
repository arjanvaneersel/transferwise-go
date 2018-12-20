package transferwise

import (
	"fmt"
	"net/http"
	"time"
)

type QuoteRequestType string

var (
	BalancePayout     QuoteRequestType = "BALANCE_PAYOUT"
	BalanceConversion QuoteRequestType = "BALANCE_CONVERSION"
)

type quoteRequest struct {
	Profile      int              `json:"profile"`
	Source       string           `json:"source"`
	Target       string           `json:"target"`
	RateType     string           `json:"rateType"`
	TargetAmount float64          `json:"targetAmount,omitempty"`
	SourceAmount float64          `json:"sourceAmount,omitempty"`
	Type         QuoteRequestType `json:"type"`
}

var None float64 = 0.0

type QuoteResponse struct {
	Profile                int              `json:"profile"`
	ID                     int              `json:"id"`
	Source                 string           `json:"source"`
	Target                 string           `json:"target"`
	RateType               string           `json:"rateType"`
	TargetAmount           float64          `json:"targetAmount,omitempty"`
	SourceAmount           float64          `json:"sourceAmount,omitempty"`
	Type                   QuoteRequestType `json:"type"`
	Rate                   float64          `json:"rate"`
	Created                time.Time        `json:"createdTime"`
	UserID                 int              `json:"createdByUserId"`
	DeliveryEstimate       time.Time        `json:"deliveryEstimate"`
	Fee                    float64          `json:"fee"`
	AllowedProfileTypes    []string         `json:"allowedProfileTypes"`
	GuaranteedTargetAmount bool             `json:"guaranteedTargetAmount,omitempty"`
	OfSourceAmount         bool             `json:"ofSourceAmount,omitempty"`
}

func (a *API) Quote(r quoteRequest) (*QuoteResponse, error) {
	d := QuoteResponse{}
	if err := a.do("v1/quotes", http.MethodPost, r, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (a *API) QuoteByID(id int) (*QuoteResponse, error) {
	d := QuoteResponse{}
	url := fmt.Sprintf("v1/quotes/%d", id)
	if err := a.do(url, http.MethodGet, nil, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

type PayInMethods struct {
	Type    string `json:"type"`
	Details struct {
		PayInReference string `json:"payInReference"`
	} `json:"details"`
}

func (a *API) PayInMethods(id int) (*PayInMethods, error) {
	d := PayInMethods{}
	url := fmt.Sprintf("v1/quotes/%d/pay-in-methods", id)
	if err := a.do(url, http.MethodGet, nil, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

type temporaryQuoteRequest struct {
	Source       string  `json:"source"`
	Target       string  `json:"target"`
	RateType     string  `json:"rateType"`
	TargetAmount float64 `json:"targetAmount,omitempty"`
	SourceAmount float64 `json:"sourceAmount,omitempty"`
}

func (a *API) TemoraryQuote(source, target string, targetAmount, sourceAmount float64) (*QuoteResponse, error) {
	if (targetAmount <= None && sourceAmount <= None) || (targetAmount > None && sourceAmount > None) {
		return nil, fmt.Errorf("specify either a target or source amount ")
	}

	req := temporaryQuoteRequest{
		Source:   source,
		Target:   target,
		RateType: "FIXED",
	}

	if targetAmount > None {
		req.TargetAmount = targetAmount
	}

	if sourceAmount > None {
		req.SourceAmount = sourceAmount
	}

	r := QuoteResponse{}
	if err := a.do("v1/quotes", http.MethodGet, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
