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

type dummyStore struct {
	Error error
	RFC   RFC
	RFCS  []RFC
}

func (ds *dummyStore) StorePreviousRFC(string, string) error {
	return ds.Error
}
func (ds *dummyStore) StoreList([]rfcEntity) {

}
func (ds *dummyStore) LoadRandom() (RFC, error) {
	return ds.RFC, ds.Error
}
func (ds *dummyStore) LoadAllPrevious() ([]RFC, error) {
	return ds.RFCS, ds.Error
}
func (ds *dummyStore) Wipe() error {
	return ds.Error
}
func (ds *dummyStore) CreateStore() error {
	return ds.Error
}
func (ds *dummyStore) Connect() {

}
