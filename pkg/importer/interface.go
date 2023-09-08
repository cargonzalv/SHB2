package importer

// `interface.go` module contains declaration of interfaces,
// provided by `importer` package
type (
	Service interface {
		LoadTifas()
		ProcessFile(string, int, int) bool
	}
)
