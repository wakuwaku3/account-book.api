package responses

type (
	// Response はレスポンスです
	Response struct {
		Errors []string    `json:"errors"`
		Result interface{} `json:"result"`
	}
)
