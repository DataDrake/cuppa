package results

import "fmt"

/*
ResultSet is a collection of the Results of a Provider query
*/
type ResultSet struct {
	results []*Result
	query   string
}

/*
NewResultSet creates as empty ResultSet for the provided query
*/
func NewResultSet(query string) *ResultSet {
	return &ResultSet{make([]*Result, 0), query}
}

/*
AddResult appends a new Result
*/
func (rs *ResultSet) AddResult(r *Result) {
	rs.results = append(rs.results, r)
}

/*
First retrieves the first result from a query
*/
func (rs *ResultSet) First() *Result {
	return rs[0]
}

/*
PrintAll pretty-prints an entire ResultSet
*/
func (rs *ResultSet) PrintAll() {
	fmt.Printf("Results of Query: '%s'\n", rs.query)
	fmt.Printf("Total Number of Results: %d\n\n", len(rs.results))
	for _, r := range rs.results {
		r.Print()
	}
}

/*
PrintFirst pretty-prints the first result
*/
func (rs *ResultSet) PrintFirst() {
	fmt.Printf("First Result of Query: '%s'\n", rs.query)
	rs.results[0].Print()
}
