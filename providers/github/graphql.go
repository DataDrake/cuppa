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

package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/config"
	"github.com/DataDrake/cuppa/results"
	"net/http"
	"os"
	"strings"
	"time"
)

// GraphQLAPI is the location of the GraphQL Endpoint
const GraphQLAPI = "https://api.github.com/graphql"

// RepoQueryFormat is the text format for the necessary GraphQL query
const RepoQueryFormat = `
query {
    repository(owner: "%s", name: "%s") {
        releases (last: %d) {
            nodes {
                name
                publishedAt
                isPrerelease
                tag {
                    name
                }
            }
        }
        refs (refPrefix: "refs/tags/", last: %d){
            nodes {
                name
            }
        }
    }
}
`

// RepoQuery is the JSON payload for this request
type RepoQuery struct {
	Query string `json:"query"`
}

// RepoQueryResult is the JSON payload of the response
type RepoQueryResult struct {
	Data struct {
		Repository struct {
			Releases struct {
				Nodes []struct {
					Name         string `json:"name"`
					PublishedAt  string `json:"publishedAt"`
					IsPrerelease bool   `json:"isPrerelease"`
					Tag          struct {
						Name string `json:"name"`
					} `json:"tag"`
				} `json:"nodes"`
			} `json:"releases"`
			Refs struct {
				Nodes []struct {
					Name string `json:"name"`
				} `json:"nodes"`
			} `json:"refs"`
		} `json:"repository"`
	} `json:"data"`
}

// Convert turns a RepoQueryResult into a Cuppa ResultSet
func (rqr RepoQueryResult) Convert(name string) (rs *results.ResultSet) {
	rs = results.NewResultSet(name)
	if nodes := rqr.Data.Repository.Releases.Nodes; len(nodes) > 0 {
		for _, node := range nodes {
			if node.IsPrerelease {
				continue
			}
			published, _ := time.Parse(time.RFC3339, node.PublishedAt)
			r := results.NewResult(node.Name, node.Tag.Name, fmt.Sprintf(SourceFormat, name, node.Tag.Name), published )
			rs.AddResult(r)
		}
	}
	if rs.Len() == 0 {
		nodes := rqr.Data.Repository.Refs.Nodes
		for _, node := range nodes {
			r := results.NewResult(name, node.Name, fmt.Sprintf(SourceFormat, name, node.Name), time.Time{})
			rs.AddResult(r)
		}
	}
	return
}

// GetReleases gets a number of releases for a given repo
func (c Provider) GetReleases(name string, max int) (rs *results.ResultSet, s results.Status) {
	names := strings.Split(name, "/")

	query := RepoQuery{
		Query: fmt.Sprintf(RepoQueryFormat, names[0], names[1], max, max),
	}
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	err := enc.Encode(&query)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	// Query the API
	req, _ := http.NewRequest("POST", GraphQLAPI, buff)
	if key := config.Global.Github.Key; len(key) > 0 {
		req.Header["Authorization"] = []string{"token " + key}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		s = results.OK
	case 404:
		s = results.NotFound
		return
	default:
		s = results.Unavailable
	}

	// Fail if not OK
	if s != results.OK {
		return
	}

	dec := json.NewDecoder(resp.Body)
	rqr := &RepoQueryResult{}
	err = dec.Decode(rqr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		s = results.Unavailable
		return
	}
	rs = rqr.Convert(name)
	return
}
