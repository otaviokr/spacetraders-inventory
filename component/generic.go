package component

// Error is the encapsulation of an error message sent by the API.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
