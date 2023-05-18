package stannp

import (
	"copilotiq/stannp-client-golang/letter"
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

	// Test data
	postUnverified := false
	test := true

	// Initialize Stannp with test data
	api := New(
		WithAPIKey(envAPIKey),
		WithPostUnverified(postUnverified),
		WithTest(test),
	)

	// Assert that the Stannp client has been initialized with the correct values
	assert.Equal(t, envAPIKey, api.apiKey, "APIKey does not match expected")
	assert.Equal(t, BaseURL, api.baseUrl, "BaseURL does not match expected")
	assert.Equal(t, postUnverified, api.postUnverified, "PostUnverified does not match expected")
	assert.Equal(t, test, api.test, "Test does not match expected")
}

func TestSendLetter(t *testing.T) {
	// Call SendLetter with a new instance of Request
	request := letter.Request{
		Test:      true,
		Template:  305202,
		ClearZone: true,
		Duplex:    true,
		Recipient: letter.RecipientDetails{
			Title:     "Mr.",
			Firstname: "John",
			Lastname:  "Doe",
			Address1:  "123 Random St",
			Town:      "Townsville",
			Zipcode:   "12345",
			State:     "Stateville",
			Country:   "US",
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