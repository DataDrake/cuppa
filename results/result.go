package results

import (
	"fmt"
	"net/url"
	"time"
)

/*
Result contains the information for a single query result
*/
type Result struct {
	Name      string
	Version   string
	Location  url.URL
	Published time.Time
}

/*
NewResult creates a result with the specified values
*/
func NewResult(name, version string, location url.URL, published time.Time) *Result {
	return &Result{name, version, location, published}
}

/*
Print pretty-prints a single Result
*/
func (r *Result) Print() {
	fmt.Printf("Name: %s\n", r.Name)
	fmt.Printf("Version: %s\n", r.Version)
	fmt.Printf("Location: %s\n", r.Location.String())
	fmt.Printf("Published: %s\n", r.Published.Format(time.RFC3339))
}
