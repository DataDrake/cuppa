package results

const (
	// OK - Query completed successfully, with results
	OK = uint8(0)
	// NotFound - Query completed successfully, without results
	NotFound = uint8(1)
	// Unavailable - Provider could not be reached
	Unavailable = uint8(2)
)

/*
Status indicates the state of a query upon completion
*/
type Status uint8
