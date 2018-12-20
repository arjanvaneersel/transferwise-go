package transferwise

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type CompanyType string

var (
	Limited                     CompanyType = "LIMITED"
	Partnership                 CompanyType = "PARTNETSHIP"
	SoleTrader                  CompanyType = "SOLE_TRADER"
	LimitedByGuarantee          CompanyType = "LIMITED_BY_GUARANTEE"
	LimitedLiabilityCompany     CompanyType = "LIMITED_LIABILITY_COMPANY"
	ForProfitCorporation        CompanyType = "FOR_PROFIT_CORPORATION"
	NonProfitCorporation        CompanyType = "NON_PROFIT_COROPORATION"
	LimitedPartnership          CompanyType = "LIMITED_PARTNERSHIP"
	LimitedLiabilityPartnership CompanyType = "LIMITED_LIABILITY_PARTNERSHIP"
	GeneralPartnership          CompanyType = "GENERAL_PARTNERSHIP"
	SoleProprietorship          CompanyType = "SOLE_PROPRIETORSHIP"
	PrivateLimitedCompany       CompanyType = "PRIVATE_LIMITED_COMPANY"
	PublicLimitedCompany        CompanyType = "PUBLIC_LIMITED_COMPANY"
	Trust                       CompanyType = "TRUST"
	OtherType                   CompanyType = "OTHER"
)

type CompanyRole string

var (
	Owner     CompanyRole = "OWNER"
	Director  CompanyRole = "DIRECTOR"
	OtherRole CompanyRole = "OTHER"
)

type Person struct {
	ID      int `json:"id"`
	Details struct {
		FirstName      string    `json:"firstName"`
		LastName       string    `json:"lastName"`
		DateOfBirth    time.Time `json:"dateOfBirth"`
		PhoneNumber    string    `json:"phoneNumber"`
		Avatar         string    `json:"avatar"`
		Occupation     string    `json:"occupation"`
		PrimaryAddress int       `json:"primaryAddress"`
	} `json:"details"`
}

type Business struct {
	ID      int `json:"id"`
	Details struct {
		Name               string      `json:"name"`
		RegistrationNumber string      `json:"registrationNumber"`
		ACN                string      `json:"acn,omitempty"`
		ABN                string      `json:"abn,omitempty"`
		ARBN               string      `json:"arbn,omitempty"`
		CompanyType        CompanyType `json:"companyType"`
		CompanyRole        CompanyRole `json:"companyRole"`
		Description        string      `json:"description"`
		Webpage            string      `json:"webpage"`
		PrimaryAddress     int         `json:"primaryAddress"`
	} `json:"details"`
}

type profile struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Details struct {
		// Personal profile fields
		FirstName   string `json:"firstName,omitempty"`
		LastName    string `json:"lastname,omitempty"`
		DateOfBirth TwDate `json:"dateOfBirth,omitempty"`
		PhoneNumber string `json:"phoneNUmber,omitempty"`
		Avatar      string `json:"avatar,omitempty"`
		Occupation  string `json:"occupation,omitempty"`
		// Business profile fields
		Name               string      `json:"name,omitempty"`
		RegistrationNumber string      `json:"registrationNumber,omitempty"`
		ACN                string      `json:"acn,omitempty"`
		ABN                string      `json:"abn,omitempty"`
		ARBN               string      `json:"arbn,omitempty"`
		CompanyType        CompanyType `json:"companyType,omitempty"`
		CompanyRole        CompanyRole `json:"companyRole,omitempty"`
		Description        string      `json:"descriptionOfBusiness,omitempty"`
		Webpage            string      `json:"webpage,omitempty"`
		// Common fields
		PrimaryAddress int `json:"primaryAddress"`
	} `json:"details"`
}

func (p profile) IsPerson() bool {
	return p.Type == "personal"
}

func (p profile) IsBusiness() bool {
	return p.Type == "business"
}

func (p profile) GetPerson() *Person {
	if !p.IsPerson() {
		return nil
	}

	o := Person{}
	deepCopy(p, &o)
	return &o
}

func (p profile) GetBusiness() *Business {
	if !p.IsBusiness() {
		return nil
	}

	o := Business{}
	deepCopy(p, &o)
	return &o
}

func (p profile) QuoteRequest(source, target string, targetAmount, sourceAmount float64, t QuoteRequestType) (quoteRequest, error) {
	if (targetAmount <= None && sourceAmount <= None) || (targetAmount > None && sourceAmount > None) {
		return quoteRequest{}, fmt.Errorf("specify either a target or source amount ")
	}

	q := quoteRequest{
		Profile:  p.ID,
		Source:   source,
		Target:   target,
		RateType: "FIXED",
		Type:     t,
	}

	if targetAmount > None {
		q.TargetAmount = targetAmount
	}

	if sourceAmount > None {
		q.SourceAmount = sourceAmount
	}

	return q, nil
}

func (a *API) Profiles() ([]profile, error) {
	res := []profile{}
	if err := a.do("v1/profiles", http.MethodGet, nil, &res); err != nil {
		return nil, err
	}

	return res, nil
}

type personalProfileRequest struct {
	Type    string                 `json:"type"`
	Details PersonalProfileRequest `json:"details"`
}

type PersonalProfileRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DateOfBirth TwDate `json:"dateOfBirth"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type businessProfileRequest struct {
	Type    string                 `json:"type"`
	Details BusinessProfileRequest `json:"details"`
}

type BusinessProfileRequest struct {
	Name               string      `json:"name"`
	RegistrationNumber string      `json:"registrationNumber"`
	ACN                string      `json:"acn,omitempty"`
	ABN                string      `json:"abn,omitempty"`
	ARBN               string      `json:"arbn,omitempty"`
	CompanyType        CompanyType `json:"companyType"`
	CompanyRole        CompanyRole `json:"companyRole"`
	Description        string      `json:"descriptionOfBusiness"`
	Webpage            string      `json:"webpage"`
}

func (a *API) CreateProfile(r interface{}) (*profile, error) {
	p := profile{}

	reqOk := false
	if pr, ok := r.(PersonalProfileRequest); ok {
		reqOk = true
		req := personalProfileRequest{
			Type:    "personal",
			Details: pr,
		}

		if err := a.do("v1/profiles", http.MethodPost, &req, &p); err != nil {
			return nil, err
		}
	}

	if br, ok := r.(BusinessProfileRequest); ok {
		reqOk = true
		req := businessProfileRequest{
			Type:    "business",
			Details: br,
		}

		if err := a.do("v1/profiles", http.MethodPost, &req, &p); err != nil {
			return nil, err
		}
	}

	if !reqOk {
		return nil, fmt.Errorf("r should be of type PersonalProfileRequest or BusinessProfileRequest")
	}

	return &p, nil
}

func (a *API) CreatePersonalProfile(r PersonalProfileRequest) (*Person, error) {
	p, err := a.CreateProfile(r)
	if err != nil {
		return nil, err
	}

	if !p.IsPerson() {
		return nil, fmt.Errorf("expected person as response, but got %v", p.Type)
	}

	return p.GetPerson(), nil
}

func (a *API) CreateBusinessProfile(r BusinessProfileRequest) (*Business, error) {
	b, err := a.CreateProfile(r)
	if err != nil {
		return nil, err
	}

	if !b.IsBusiness() {
		return nil, fmt.Errorf("expected business as response, but got %v", b.Type)
	}

	return b.GetBusiness(), nil
}

func (a *API) UpdateProfile(r interface{}) (*profile, error) {
	p := profile{}

	reqOk := false
	if pr, ok := r.(PersonalProfileRequest); ok {
		reqOk = true
		req := personalProfileRequest{
			Type:    "personal",
			Details: pr,
		}

		if err := a.do("v1/profiles", http.MethodPut, &req, &p); err != nil {
			return nil, err
		}
	}

	if br, ok := r.(BusinessProfileRequest); ok {
		reqOk = true
		req := businessProfileRequest{
			Type:    "business",
			Details: br,
		}

		if err := a.do("v1/profiles", http.MethodPut, &req, &p); err != nil {
			return nil, err
		}
	}

	if !reqOk {
		return nil, fmt.Errorf("r should be of type PersonalProfileRequest or BusinessProfileRequest")
	}

	return &p, nil
}

func (a *API) UpdatePersonalProfile(r PersonalProfileRequest) (*Person, error) {
	p, err := a.UpdateProfile(r)
	if err != nil {
		return nil, err
	}

	if !p.IsPerson() {
		return nil, fmt.Errorf("expected person as response, but got %v", p.Type)
	}

	return p.GetPerson(), nil
}

func (a *API) UpdateBusinessProfile(r BusinessProfileRequest) (*Business, error) {
	b, err := a.UpdateProfile(r)
	if err != nil {
		return nil, err
	}

	if !b.IsBusiness() {
		return nil, fmt.Errorf("expected business as response, but got %v", b.Type)
	}

	return b.GetBusiness(), nil
}

func (a *API) GetProfile(id int) (*profile, error) {
	p := profile{}
	if err := a.do("v1/profiles"+strconv.Itoa(id), http.MethodGet, nil, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (a *API) GetPerson(id int) (*Person, error) {
	p, err := a.GetProfile(id)
	if err != nil {
		return nil, err
	}

	if !p.IsPerson() {
		return nil, fmt.Errorf("expected person as response, but got %v", p.Type)
	}

	return p.GetPerson(), nil
}

func (a *API) GetBusiness(id int) (*Business, error) {
	b, err := a.GetProfile(id)
	if err != nil {
		return nil, err
	}

	if !b.IsBusiness() {
		return nil, fmt.Errorf("expected business as response, but got %v", b.Type)
	}

	return b.GetBusiness(), nil
}

type DocumentType string

var (
	DriversLicence DocumentType = "DRIVERS_LICENSE"
	IdentityCard                = "IDENTITY_CARD"
	GreenCard                   = "GREEN_CARD"
	MyNumber                    = "MY_NUMBER"
	Passport                    = "PASSPORT"
	Other                       = "OTHER"
)

type verificationDocument struct {
	FirstName        string       `json:"firstName"`
	LastName         string       `json:"lastName"`
	Type             DocumentType `json:"type"`
	UniqueIdentifier string       `json:"uniqueIdentifier"`
	IssueDate        TwDate       `json:"issueDate"`
	IssuerCountry    string       `json:"issuerCountry"`
	IssuerState      string       `json:"issuerState"`
	ExpiryDate       TwDate       `json:"expiryDate,omitempty"`
}

var NoExpiry time.Time

func (a *API) VerificationDocument(p *Person, t DocumentType, id string, issued time.Time, country string, state string, expires time.Time) error {
	d := verificationDocument{
		FirstName:        p.Details.FirstName,
		LastName:         p.Details.LastName,
		Type:             t,
		UniqueIdentifier: id,
		IssueDate:        TwDate{issued},
		IssuerCountry:    country,
		IssuerState:      state,
	}

	if !expires.IsZero() {
		d.ExpiryDate = TwDate{expires}
	}

	url := fmt.Sprintf("v1/profiles/%d/verification-documents", p.ID)
	if err := a.do(url, http.MethodPost, d, nil); err != nil {
		return err
	}

	return nil
}
