package variable

// Meta defines meta format format for api format
type Meta struct {
	Version   string `json:"version,omitempty"`
	APIEnv    string `json:"api_env,omitempty"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	NodeName  string `json:"node_name,omitempty"`
}

// ResponseCodeMessages ..
type ResponseCodeMessages struct {
	EN string `mapstructure:"en"`
	ID string `mapstructure:"id"`
}

// ResponseCode ..
type ResponseCode struct {
	ErrorCode  string                 `mapstructure:"error_code"`
	Messages   []ResponseCodeMessages `mapstructure:"message"`
	StatusCode int                    `mapstructure:"status_code"`
}

// ResponseDesc defines details data response
type ResponseDesc struct {
	ID string `json:"id,omitempty"`
	EN string `json:"en,omitempty"`
}

// Response ...
type Response struct {
	Meta         interface{}   `json:"meta"`
	ResponseDesc *ResponseDesc `json:"response_desc,omitempty"`
	ResponseCode string        `json:"response_code,omitempty"`
}

// SuccessResponse ...
type SuccessResponse struct {
	Response
	Data interface{} `json:"data,omitempty"`
	Next string      `json:"next,omitempty"`
}

// ErrorResponse ...
type ErrorResponse struct {
	Data interface{} `json:"data,omitempty"`
	Response
	HTTPStatus int `json:"-"`
}
