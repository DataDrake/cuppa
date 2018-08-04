//
// Copyright 2016-2018 Bryan T. Meyers <bmeyers@datadrake.com>
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

import (
	"fmt"
	"sort"
	"strings"
)

// ResultSet is a collection of the Results of a Provider query
type ResultSet struct {
	results []*Result
	query   string
}

// NewResultSet creates as empty ResultSet for the provided query
func NewResultSet(query string) *ResultSet {
	return &ResultSet{make([]*Result, 0), query}
}

// AddResult appends a new Result
func (rs *ResultSet) AddResult(r *Result) {
	rs.results = append(rs.results, r)
}

// Empty checks if there were no results
func (rs *ResultSet) Empty() bool {
	return len(rs.results) == 0
}

// First retrieves the first result from a query
func (rs *ResultSet) First() *Result {
	return rs.results[0]
}

// Last retrieves the first result from a query
func (rs *ResultSet) Last() *Result {
    switch len(rs.results) {
    case 0:
        return nil
    case 1:
        return rs.results[0]
    default:
	    return rs.results[len(rs.results)-1]
    }
}

// PrintAll pretty-prints an entire ResultSet
func (rs *ResultSet) PrintAll() {
	fmt.Printf("%-25s: '%s'\n", "Results of Query", rs.query)
	fmt.Printf("%-25s: %d\n\n", "Total Number of Results", len(rs.results))
	sort.Sort(rs)
	for _, r := range rs.results {
		r.Print()
	}
}

// PrintFirst pretty-prints the first result
func (rs *ResultSet) PrintFirst() {
	fmt.Printf("First Result of Query: '%s'\n", rs.query)
	rs.results[0].Print()
}

// Len is the number of elements in the ResultSet (sort.Interface)
func (rs *ResultSet) Len() int {
	return len(rs.results)
}

// Less reports whether the element with
// index i should sort before the element with index j. (sort.Interface)
func (rs *ResultSet) Less(i, j int) bool {
	return strings.Compare(rs.results[i].Version, rs.results[j].Version) == -1
}

// Swap swaps the elements with indexes i and j.
func (rs *ResultSet) Swap(i, j int) {
	rs.results[i], rs.results[j] = rs.results[j], rs.results[i]
}
