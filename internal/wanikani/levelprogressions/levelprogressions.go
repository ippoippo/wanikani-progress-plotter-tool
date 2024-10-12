package levelprogressions

import (
	"encoding/json"
	"time"
)

type ResponseBody struct {
	TotalCount    int                 `json:"total_count"`
	DataUpdatedAt time.Time           `json:"data_updated_at"`
	Data          []*LevelProgression `json:"data"`
}

type LevelProgression struct {
	Id   int64                 `json:"id"`
	Data *LevelProgressionData `json:"data"`
}

type LevelProgressionData struct {
	Level      int       `json:"level"`
	UnlockedAt time.Time `json:"unlocked_at"`
	PassedAt   time.Time `json:"passed_at"`
}

func NewResponseFrom(b []byte) (*ResponseBody, error) {
	response := &ResponseBody{}
	err := json.Unmarshal(b, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
