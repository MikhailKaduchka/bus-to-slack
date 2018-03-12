package departures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrl(t *testing.T) {
	url, err := getEnv(LOCATION_URL, URL_SAMPLE)
	assert.Equal(t, url, "")
	assert.NotEqual(t, err, nil)
}

func TestGetWebHook(t *testing.T) {
	url, err := getEnv(SLACK_WEBHOOK, WEBHOOK_SAMPLE)
	assert.Equal(t, url, "")
	assert.NotEqual(t, err, nil)
}

// func TestHandler(t *testing.T) {
// 	tests := []struct {
// 		request Request
// 		expect  Response
// 		err     error
// 	}{
// 		{
// 			// Test that the handler responds with the correct response
// 			// when a valid name is provided in the HTTP body
// 			request: Request{ID: "test-request-id", Value: "Test"},
// 			expect:  Response{Message: "Processed request ID test-request-id", Ok: true},
// 			err:     nil,
// 		},
// 	}

// 	os.Setenv(LOCATION_URL, URL_SAMPLE)
// 	os.Setenv(SLACK_WEBHOOK, WEBHOOK_SAMPLE)
// 	for _, test := range tests {
// 		response, err := Handler(test.request)
// 		assert.IsType(t, test.err, err)
// 		log.Println("Response: " + response.Message)
// 		assert.Equal(t, test.expect.Message, response.Message)
// 	}
// }
