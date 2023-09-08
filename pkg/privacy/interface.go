package privacy

// `interface.go` module contains declaration of interfaces,
// provided by `privacy` package
type (
	Service interface {
		Optout([]byte) bool
		SetOptout([]byte) []byte
	}
)
