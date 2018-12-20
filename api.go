package transferwise

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	url        = "https://api.transferwise.com/"
	sandboxURL = "https://api.sandbox.transferwise.tech/"
)

func deepCopy(s interface{}, d interface{}) error {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(s); err != nil {
		return err
	}

	if err := gob.NewDecoder(b).Decode(d); err != nil {
		return err
	}

	return nil
}

type TwDate struct {
	time.Time
}

func (d *TwDate) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02", strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}

	(*d) = TwDate{t}
	return nil
}

func (d *TwDate) MarshalJSON() ([]byte, error) {
	// Using explicit functions to construct the date, because d.Format("2006-01-02") isn't working,
	// it returns YYYY-MM-DD 00:00:00 +0000 UTC instead of just YYYY-MM-DD
	t := d.Format("2006-01-02") //fmt.Sprintf("\"%d-%d-%d\"", d.Year(), d.Month(), d.Day())
	if t == "" {
		return nil, fmt.Errorf("conversion error")
	}

	return []byte(fmt.Sprintf("\"%s\"", t)), nil
}

type Language string

var (
	AmericanEnglish Language = "en_US"
	BritishEnglish           = "en"
	Dutch                    = "nl"
	French                   = "fr"
	German                   = "de"
	Hungarian                = "hu"
	Italian                  = "it"
	Japanese                 = "ja"
	Korean                   = "ko"
	Polish                   = "po"
	Portugese                = "pt"
	Romanian                 = "ro"
	Russian                  = "ru"
	Spanish                  = "es"
)

type APIError struct {
	Errors []struct {
		Code      string        `json:"code"`
		Message   string        `json:"message"`
		Path      string        `json:"path"`
		Arguments []interface{} `json:"arguments"`
	} `json:"errors"`
}

func (a APIError) Error() string {
	s := ""
	l := len(a.Errors)
	for _, err := range a.Errors {
		s += fmt.Sprintf("%s: %s, path: %s, arguments: %v", err.Code, err.Message, err.Path, err.Arguments)
		if l > 1 {
			s += fmt.Sprintf("\n")
		}
	}

	return s
}

type API struct {
	url   string
	token string
	lang  Language
}

type ReqOption func(*http.Request) error

func WithReqLanguage(l Language) ReqOption {
	return func(r *http.Request) error {
		r.Header.Set("Accept-Language", string(l))
		return nil
	}
}

func (a *API) do(url string, method string, body interface{}, d interface{}, options ...ReqOption) error {
	b := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return fmt.Errorf("json encoding error: %v", err)
		}
	}

	req, err := http.NewRequest(method, a.url+url, b)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	req.Header.Set("Content-Type", "application/json")
	if a.lang != "" {
		req.Header.Set("Accept-Language", string(a.lang))
	}

	for _, opt := range options {
		if err := opt(req); err != nil {
			return fmt.Errorf("option error: %v", err)
		}
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusIMUsed {
		err := APIError{}
		json.NewDecoder(res.Body).Decode(&err)
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(d); err != nil {
		return fmt.Errorf("json decoding error: %v", err)
	}

	return nil
}

type APIOption func(*API) error

func WithURL(url string) APIOption {
	return func(a *API) error {
		a.url = url
		return nil
	}
}

func WithSandbox() APIOption {
	return WithURL(sandboxURL)
}

func WithLanguage(l Language) APIOption {
	return func(a *API) error {
		a.lang = l
		return nil
	}
}

func New(token string, options ...APIOption) (*API, error) {
	api := API{
		url:   url,
		token: token,
	}

	for _, opt := range options {
		if err := opt(&api); err != nil {
			return nil, fmt.Errorf("option error: %v", err)
		}
	}

	return &api, nil
}
