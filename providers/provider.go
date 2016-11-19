package providers

import "github.com/DataDrake/cuppa/results"

/*
Provider provides a common interface for each of the backend providers
*/
type Provider interface {
	Search(Name string) (*results.ResultSet, results.Status)
}
