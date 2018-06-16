package main

func main() {
	rfc := new(RFC)
	pgStore := new(PostgresStore)
	rfc.SetStore(pgStore)

	pgStore.CreateStore()

	rfc.DownloadRFCList()
	rfcs := rfc.parseListConcurrent("list.txt")
	pgStore.StoreList(rfcs)

	rfc.WriteOutRandomRFC()
	rfc.WriteOutAllPreviousRFCHTML()
}
