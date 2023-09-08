package demand

// `interface.go` module contains declaration of interfaces,
// provided by `demand` package
type (
	Service interface {
		BidOrtbReq(demandParams DemandExtParams) (int, []byte)
	}
)
