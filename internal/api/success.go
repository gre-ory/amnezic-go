package api

// //////////////////////////////////////////////////
// encode

func toJsonSuccess() *JsonSuccessResponse {
	return &JsonSuccessResponse{
		Success: true,
	}
}

type JsonSuccessResponse struct {
	Success bool `json:"success,omitempty"`
}
