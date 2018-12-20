package transferwise

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

const tokenKey = "TRANSFERWISE_API_TOKEN"

var token string

func ensureToken(t *testing.T) {
	if token == "" {
		token = os.Getenv(tokenKey)
		if token == "" {
			t.Fatalf("%s not set", tokenKey)
		}
	}
}

func TestNew(t *testing.T) {
	ensureToken(t)

	t.Run("noOptions", func(t *testing.T) {
		api, err := New(token)
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if api.token != token {
			t.Errorf("expected another token than %v", api.token)
		}

		if api.url != url {
			t.Errorf("expected url to be %v, but got %v", url, api.url)
		}
	})

	t.Run("withStandbox", func(t *testing.T) {
		api, err := New(token, WithSandbox())
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if api.token != token {
			t.Errorf("expected another token than %v", api.token)
		}

		if api.url != sandboxURL {
			t.Errorf("expected url to be %v, but got %v", sandboxURL, api.url)
		}
	})
}

func TestTwDate(t *testing.T) {
	type profile struct {
		Name string `json:"name"`
		Dob  TwDate `json:"dob"`
	}

	t.Run("unmarshalling", func(t *testing.T) {
		j := `{
			"name":"Test person",
			"dob":"1977-12-18"
		}`

		p := profile{}
		err := json.Unmarshal([]byte(j), &p)
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if !strings.Contains(p.Dob.String(), "1977-12-18") {
			t.Errorf("expected dob to contain 1977-12-18, but is %v", p.Dob.String())
		}
	})

	t.Run("marshalling", func(t *testing.T) {
		dob, err := time.Parse("02-01-2006", "18-12-1977")
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		p := profile{
			Name: "Test Person",
			Dob:  TwDate{dob},
		}

		b, err := json.Marshal(&p)
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		j := string(b)

		if !strings.Contains(j, "1977-12-18") {
			t.Errorf("expected dob to contain 1977-12-18, but is %v", p.Dob.String())
		}
	})
}

// func TestQuote(t *testing.T) {
// 	ensureToken(t)

// 	api, err := New(token, WithSandbox())
// 	if err != nil {
// 		t.Fatalf("expected to pass, but got %v", err)
// 	}

// 	p, err := api.Profiles()
// 	if err != nil {
// 		t.Fatalf("expected to pass, but got %v", err)
// 	}

// 	if len(p) == 0 {
// 		t.Fatal("expected to get at least one profile")
// 	}

// 	qr, err := p[0].QuoteRequest("EUR", "GBP", 600.00, None, BalancePayout)
// 	if err != nil {
// 		t.Fatalf("expected to pass, but got: %v", err)
// 	}

// 	r, err := api.Quote(qr)
// 	if err != nil {
// 		t.Fatalf("expected to pass, but got: %v", err)
// 	}

// 	if r.TargetAmount <= None {
// 		t.Errorf("expected a targetAmount, but got %.2f", r.TargetAmount)
// 	}

// 	log.Printf("got quote: %v", r)
// }
