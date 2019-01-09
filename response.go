package ce

import "encoding/json"

type formDefenderResponse struct {
	Status string `json:"status"`
	// can be either ValidResponse or ErrorResponse
	Data json.RawMessage `json:"data"`
}

type validResponse struct {
	Valid    bool    `json:"valid"`
	Status   Status  `json:"status"`
	Duration float64 `json:"duration"`
}

type errorResponse struct {
	ID   string `json:"id"`
	Code int    `json:"code"`
}
