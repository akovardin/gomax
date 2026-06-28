package types

type MaxApiError struct {
	Error            string `json:"error"`
	Message          string `json:"message"`
	Title            string `json:"title"`
	LocalizedMessage string `json:"localizedMessage"`
}
