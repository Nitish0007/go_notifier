package rabbitmq_utils

import (
	"encoding/json"
	"github.com/google/uuid"
)

type JobMessage struct {
	JobID string `json:"job_id"`
	Payload map[string]any `json:"payload"`
}

func NewJobMessage(payload map[string]any) *JobMessage {
	jobID, exists := payload["job_id"].(string)
	if !exists || jobID == "" {
		jobID = uuid.New().String()
	}

	body, exists := payload["payload"].(map[string]any)
	if !exists {
		body = map[string]any{}
	}
	return &JobMessage{
		JobID: jobID,
		Payload: body,
	}
}

func (m *JobMessage) GetJobID() string {
	return m.JobID
}

func (m *JobMessage) GetPayload() map[string]any {
	return m.Payload
}

func (m *JobMessage) SetPayload(payload map[string]any) {
	m.Payload = payload
}

func (m *JobMessage) SetJobID(jobID string) {
	m.JobID = jobID
}

func (m *JobMessage) ToJSON() ([]byte, error) {
	json, err := json.Marshal(m)
	return json, err
}

func (m *JobMessage) FromJSON(jsonData []byte) error {
	return json.Unmarshal(jsonData, m)
}

func (m *JobMessage) ToMap() map[string]any {
	return m.Payload
}

// func (m *JobMessage) FromMap(data map[string]any) error {
// 	m.Payload = data
// 	return nil
// }