package main

// Store describes a storage mechanism.
type Store interface {
	StorePreviousRFC(string, string) error
	StoreList([]rfcEntity)
	LoadRandom() (RFC, error)
	LoadAllPrevious() ([]RFC, error)
	Wipe() error
	CreateStore() error
	Connect()
}
