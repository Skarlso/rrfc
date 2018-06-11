package main

// Store describes a storage mechanism.
type Store interface {
	Store(string, string)
	StoreConcurrently([]string)
	LoadRandom() RFC
	LoadAllPrevious() []RFC
}
