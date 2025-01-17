package stannp

import (
	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/jgroeneveld/trial/assert"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

const ApiKeyEnvKey = "STANNP_API_KEY"

var TestClient *Stannp

func TestMain(m *testing.M) {
	setup()
	m.Run()
}

func setup() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Unable to load .env file: %s", err)
	}

	apiKey := os.Getenv(ApiKeyEnvKey)
	if apiKey == "" {
		log.Fatalf("Cannot proceed when apiKey is the empty string [%s]", apiKey)
	}

	// Initialize Stannp with test data
	TestClient = New(
		WithAPIKey(apiKey),
		WithClearZone(false),
		WithDuplex(false),
		WithIdempotencyFunc(DefaultIdemFunc),
		WithPostUnverified(false),
		WithTest(true),
	)

	if !TestClient.IsTest() {
		log.Fatalf("Cannot proceed when API key is live [%s]", apiKey)
	}
}

//goland:noinspection GoBoolExpressions
func TestNew(t *testing.T) {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	// Get API key from environment variable
	envAPIKey, exists := os.LookupEnv("STANNP_API_KEY")
	if !exists {
		t.Fatal("STANNP_API_KEY not set in .env file")
	}

	// Initialize Stannp with test data
	api := New(
		WithAPIKey(envAPIKey),
		WithClearZone(false),
		WithDuplex(false),
		WithPostUnverified(true),
		WithTest(true),
		WithIdempotencyFunc(DefaultIdemFunc),
	)

	// Assert that the Stannp client has been initialized with the correct values
	assert.Equal(t, envAPIKey, api.apiKey)
	assert.Equal(t, BaseURL, api.baseUrl)
	assert.Equal(t, false, api.clearZone)
	assert.Equal(t, false, api.duplex)
	assert.Equal(t, true, api.postUnverified)
	assert.Equal(t, true, api.test)
	assert.NotNil(t, api.idemFunc)
	assert.Equal(t, 36, len(api.idemFunc()))
}

func TestSendLetter(t *testing.T) {
	// Call SendLetter with a new instance of SendReq
	request := &letter.SendReq{
		Template: "305202",
		Recipient: letter.RecipientDetails{
			Address1:  "9355 Burton Way",
			Address2:  "Courthouse",
			Country:   "United States",
			Firstname: "Judge",
			Lastname:  "Judy",
			State:     "CA",
			Title:     "Mrs.",
			Town:      "Beverly Hills",
			Zipcode:   "90210",
		},
	}

	// Note: This call is not actually sending a request.
	response, apiErr := TestClient.SendLetter(request)
	assert.True(t, reflect.ValueOf(apiErr).IsNil())

	assert.True(t, response.Success)
	assert.Equal(t, "0.81", response.Data.Cost)
	assert.Equal(t, "US-LETTER", response.Data.Format)
	assert.Equal(t, "test", response.Data.Status)
	assert.True(t, strings.HasPrefix(response.Data.Pdf, "https://us.stannp.com/api/v1/storage/get/"))
}

func TestValidateAddress(t *testing.T) {
	t.Run("verify is_valid is false for fake data", func(t *testing.T) {
		request := &address.ValidateReq{
			Address1: "9354444445 Burton Way",
			City:     "Beverly Hills",
			Company:  "Beverly Hills Courthouse",
			Country:  "US",
			State:    "CA",
			Zipcode:  "90210",
		}

		validateRes, apiErr := TestClient.ValidateAddress(request)
		assert.True(t, reflect.ValueOf(apiErr).IsNil())
		assert.False(t, validateRes.Data.IsValid)
		assert.True(t, validateRes.Success)
	})
	t.Run("verify is_valid is true for real data", func(t *testing.T) {
		request := &address.ValidateReq{
			Address1: "9355 Burton Way",
			City:     "Beverly Hills",
			Company:  "Beverly Hills Courthouse",
			Country:  "US",
			State:    "CA",
			Zipcode:  "90210",
		}

		validateRes, apiErr := TestClient.ValidateAddress(request)
		assert.True(t, reflect.ValueOf(apiErr).IsNil())
		assert.True(t, validateRes.Data.IsValid)
		assert.True(t, validateRes.Success)
	})
}
