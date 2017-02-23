package transport

type ConnectionOpts struct {
	DisableTransportKeyPinning bool
	AllowTransportSelfSigned   bool
	ID                         string
	TransportHost              string
	TransportPort              int
	PasscodeHash               string
	PasscodeSalt               string
	Tags                       []string
}
