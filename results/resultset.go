//
// Copyright © 2017 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
	return rs.results[0]
}

/*
PrintAll pretty-prints an entire ResultSet
*/
func (rs *ResultSet) PrintAll() {
	fmt.Printf("%-25s: '%s'\n", "Results of Query", rs.query)
	fmt.Printf("%-25s: %d\n\n", "Total Number of Results", len(rs.results))
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
