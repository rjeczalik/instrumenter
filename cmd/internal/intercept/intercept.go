package intercept

import "time"

type Error struct {
	ID         string    `json:"id,omitempty"`
	Path       string    `json:"path"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"createdAt"`
	Stacktrace []byte    `json:"stacktrace"`
	Type       string    `json:"type"`
	Message    string    `json:"message"`
}

type NewRequest struct {
	Error *Error `json:"error"`
}

type NewResponse struct {
	ID string `json:"id"`
}
