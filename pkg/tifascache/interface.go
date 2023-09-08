package tifascache

import "time"

// `interface.go` module contains declaration of interfaces,
// provided by `tifascache` package
type (
	Service interface {
		GetTifa(string) bool
		SetTifa(string, string, time.Duration)
		IsLastLoadTsExpired(string) bool
	}
)
