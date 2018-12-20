package transferwise

import (
	"testing"
	"time"
)

func TestProfiles(t *testing.T) {
	ensureToken(t)

	api, err := New(token, WithSandbox())
	if err != nil {
		t.Fatalf("expected to pass, but got %v", err)
	}

	p, err := api.Profiles()
	if err != nil {
		t.Fatalf("expected to pass, but got %v", err)
	}

	if l := len(p); l == 0 {
		t.Errorf("expected at least one profile, but got %d profiles", l)
	}

	t.Logf("profiles: %#v", p)
}

func TestCreateProfile(t *testing.T) {
	ensureToken(t)

	api, err := New(token, WithSandbox())
	if err != nil {
		t.Fatalf("expected to pass, but got %v", err)
	}

	t.Run("personal", func(t *testing.T) {
		dob, _ := time.Parse("02-01-2006", "18-12-1977")
		req := PersonalProfileRequest{
			FirstName:   "Test",
			LastName:    "Person",
			DateOfBirth: TwDate{dob},
			PhoneNumber: "+37211223344",
		}

		p, err := api.CreatePersonalProfile(req)
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		if p.ID == 0 {
			t.Fatalf("expected profile to have an ID")
		}

		t.Logf("profile: %#v", p)
	})

	t.Run("business", func(t *testing.T) {
		req := BusinessProfileRequest{
			Name:               "Test company Ltd",
			RegistrationNumber: "01234567",
			CompanyType:        PrivateLimitedCompany,
			CompanyRole:        Director,
			Description:        "Software development",
			Webpage:            "www.example.com",
		}

		p, err := api.CreateBusinessProfile(req)
		if err != nil {
			t.Fatalf("expected to pass, but got %v", err)
		}

		t.Logf("profile: %#v", p)
	})
}
