package types

// CallbackHandler is a struct that contains the details needed to use the 9 Spokes callback handler (cb.9spokes.io)
type CallbackHandler struct {
	URL         string
	ReturnURL   string
	RedirectURI string
	IV          string
	Secret      string
}
