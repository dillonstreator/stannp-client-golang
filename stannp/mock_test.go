package stannp

import (
	"github.com/CopilotIQ/stannp-client-golang/address"
	"github.com/CopilotIQ/stannp-client-golang/letter"
	"github.com/CopilotIQ/stannp-client-golang/util"
	"github.com/jgroeneveld/trial/assert"
	"reflect"
	"testing"
)

func TestNewMockClient(t *testing.T) {
	tests := []struct {
		name   string
		opts   []MockOption
		expect MockClient
	}{
		{
			name:   "no options",
			expect: MockClient{},
		},
		{
			name: "with failNext",
			opts: []MockOption{
				WithFailNext(true),
			},
			expect: MockClient{failNext: true},
		},
		{
			name: "with invalidNext",
			opts: []MockOption{
				WithInvalidNext(true),
			},
			expect: MockClient{invalidNext: true},
		},
		{
			name: "with codeNext",
			opts: []MockOption{
				WithCodeNext(400),
			},
			expect: MockClient{codeNext: 400},
		},
		{
			name: "with errNext",
			opts: []MockOption{
				WithErrNext("error"),
			},
			expect: MockClient{errNext: "error"},
		},
		{
			name: "with all options",
			opts: []MockOption{
				WithFailNext(true),
				WithInvalidNext(true),
				WithCodeNext(400),
				WithErrNext("simulated error"),
			},
			expect: MockClient{failNext: true, invalidNext: true, codeNext: 400, errNext: "simulated error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockClient(tt.opts...)
			assert.Equal(t, tt.expect, *client)
		})
	}
}

func TestMockSendLetter(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		expectedSuccess   bool
		expectedError     *util.APIError
	}{
		{
			name:              "success expected err not expected",
			mockClientOptions: []MockOption{},
			expectedSuccess:   true,
			expectedError:     nil,
		},
		{
			name:              "success not expected err expected",
			mockClientOptions: []MockOption{WithFailNext(true)},
			expectedSuccess:   false,
			expectedError: &util.APIError{
				Code:    500,
				Error:   "failNext is true",
				Success: false,
			},
		},
		{
			name: "err expected code expected custom err expected",
			mockClientOptions: []MockOption{
				WithFailNext(true),
				WithCodeNext(404),
				WithErrNext("custom message"),
			},
			expectedSuccess: false,
			expectedError: &util.APIError{
				Code:    404,
				Error:   "custom message",
				Success: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(tt.mockClientOptions...)
			sendLetterRes, apiErr := mockClient.SendLetter(&letter.SendReq{})

			if tt.expectedError != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, *tt.expectedError, *apiErr)
				assert.True(t, reflect.ValueOf(sendLetterRes).IsNil())
			} else {
				assert.True(t, reflect.ValueOf(apiErr).IsNil())
				assert.NotNil(t, sendLetterRes)
				assert.Equal(t, tt.expectedSuccess, sendLetterRes.Success)

				assert.Equal(t, sendLetterRes.Data.Status, "received")
				assert.True(t, sendLetterRes.Data.Cost != "")
				assert.True(t, sendLetterRes.Data.Created != "")
				assert.True(t, sendLetterRes.Data.Format != "")
				assert.True(t, sendLetterRes.Data.Id != "")
				assert.True(t, sendLetterRes.Data.Pdf != "")
				assert.True(t, sendLetterRes.Data.Status != "")
			}
		})
	}
}

func TestMockValidateAddress(t *testing.T) {
	tests := []struct {
		name              string
		mockClientOptions []MockOption
		isValidExpected   bool
		errExpected       *util.APIError
	}{
		{
			name:              "valid expected err not expected",
			mockClientOptions: []MockOption{},
			isValidExpected:   true,
			errExpected:       nil,
		},
		{
			name:              "valid not expected err not expected",
			mockClientOptions: []MockOption{WithInvalidNext(true)},
			isValidExpected:   false,
			errExpected:       nil,
		},
		{
			name:              "err expected",
			mockClientOptions: []MockOption{WithFailNext(true)},
			isValidExpected:   false,
			errExpected: &util.APIError{
				Code:    500,
				Error:   "failNext is true",
				Success: false,
			},
		},
		{
			name: "fail next code next err next",
			mockClientOptions: []MockOption{
				WithCodeNext(400),
				WithErrNext("custom message"),
				WithFailNext(true),
			},
			isValidExpected: false,
			errExpected: &util.APIError{
				Code:    400,
				Error:   "custom message",
				Success: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := NewMockClient(tt.mockClientOptions...)
			validateAddressRes, apiErr := mockClient.ValidateAddress(&address.ValidateReq{})

			if tt.errExpected != nil {
				assert.NotNil(t, apiErr)
				assert.Equal(t, *tt.errExpected, *apiErr)
			} else {
				assert.True(t, reflect.ValueOf(apiErr).IsNil())

				if tt.isValidExpected {
					assert.True(t, validateAddressRes.Data.IsValid)
				} else {
					assert.False(t, validateAddressRes.Data.IsValid)
				}
			}
		})
	}
}