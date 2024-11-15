package model

type WebResponse struct {
	Message string  `json:"message"`
	Error   *string `json:"error"`
	Data    any     `json:"data"`
}

type PaginationMetadata struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

func NewWebResponse(message string, err error, data any) WebResponse {
	if err != nil {
		errMsg := err.Error()

		return WebResponse{Message: message, Error: &errMsg, Data: data}
	}

	return WebResponse{Message: message, Data: data}
}
