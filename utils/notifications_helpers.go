package utils

import (
	"errors"
	"strings"
	// "sync"
)

// This method is beign called inside ValidateBulkNotificationPayload, So make sure changes in this won't negatively affect the flow
func ValidateNotificationPayload(payload map[string]any) (map[string]any, error) {
	notificationType, exists := payload["channel"].(string)
	if !exists {
		return nil, errors.New("channel is required in notification payload")
	}

	if !IsValidChannelType(notificationType) {
		return nil, errors.New("invalid channel type provided")
	}

	if recipient, exists := payload["recipient"].(string); !exists || recipient == "" {
		return nil, errors.New("'recipient' field is required in notification payload")
	}

	// only validated email format for now
	// TODO: add validations for other channels recipient formats
	if notificationType == "email" && !ValidateEmail(payload["recipient"].(string)) {
		return nil, errors.New("recipient is invalid email format")
	}

	// only one of body or html_body is required (if body is provided then html_body is ignored)
	body, exists := payload["body"]
	if !exists || body == nil || strings.TrimSpace(body.(string)) == "" {
		payload["body"] = ""
	}

	htmlBody, exists := payload["html_body"]
	if !exists || htmlBody == nil || strings.TrimSpace(htmlBody.(string)) == "" {
		payload["html_body"] = ""
	}

	if htmlBody == "" && body == "" {
		return nil, errors.New("either 'body' or 'html_body' must be provided in notification payload")
	}

	return payload, nil
}

// func ValidateBulkNotificationPayload(payload []map[string]any) (valid []map[string]any, invalid []map[string]any) {
// 	payloadCollection := make(chan map[string]any)   // contain all payload
// 	validPayloads := make(chan map[string]any)   // valid chunks of payloads
// 	invalidPayloads := make(chan map[string]any) // chunks of invalid payloads

// 	var workers int
// 	wg := sync.WaitGroup{}

// 	// set the number of workers based on the payload size
// 	switch payloadSize := len(payload); {
// 	case payloadSize < 5:
// 		workers = 1
// 	case payloadSize > 5 && payloadSize < 60:
// 		workers = 3
// 	default:
// 		workers = 5
// 	}

// 	// worker pool pattern
// 	// Initializing go routines(workers)
// 	for range workers {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for p := range payloadCollection {
// 				vp, err := ValidateNotificationPayload(p)
// 				if err != nil {
// 					invalidPayloads <- p
// 				} else {
// 					validPayloads <- vp
// 				}
// 			}
// 		}()
// 	}

// 	// send all payloads to workers via channel
// 	go func() {
// 		defer close(payloadCollection)
// 		for _, p := range payload {
// 			payloadCollection <- p
// 		}
// 	}()

// 	// close channels
// 	wg.Wait()
// 	close(validPayloads)
// 	close(invalidPayloads)

// 	validPayloadSlice := make([]map[string]any, 0)
// 	invalidePayloadSlice := make([]map[string]any, 0)
// 	// adding data to slices from channels
// 	for vp := range validPayloads {
// 		validPayloadSlice = append(validPayloadSlice, vp)
// 	}

// 	for ivp := range invalidPayloads {
// 		invalidePayloadSlice = append(invalidePayloadSlice, ivp)
// 	}

// 	return validPayloadSlice, invalidePayloadSlice
// }
