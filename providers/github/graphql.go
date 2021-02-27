//
// Copyright 2016-2021 Bryan T. Meyers <root@datadrake.com>
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
	log "github.com/DataDrake/waterlog"
	"net/http"
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
                target {
                    ... on Commit {
                        committedDate
                    }
                    ... on Tag {
                        tagger {
                            date
                        }
                    }
                }
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
					Name   string `json:"name"`
					Target struct {
						Date   string `json:"committedDate"`
						Tagger struct {
							Date string `json:"date"`
						} `json:"tagger"`
					} `json:"target"`
				} `json:"nodes"`
			} `json:"refs"`
		} `json:"repository"`
	} `json:"data"`
}

// Convert turns a RepoQueryResult into a Cuppa ResultSet
func (rqr RepoQueryResult) Convert(name string) (rs *results.ResultSet) {
	rs = results.NewResultSet(name)
	var err error
	for _, tag := range rqr.Data.Repository.Refs.Nodes {
		pre := false
		found := false
		var published time.Time
		for _, node := range rqr.Data.Repository.Releases.Nodes {
			if node.Tag.Name == tag.Name {
				found = true
				if node.IsPrerelease {
					pre = true
				}
				published, _ = time.Parse(time.RFC3339, node.PublishedAt)
			}
		}
		if pre {
			continue
		}
		if !found {
			if len(tag.Target.Date) > 0 {
				published, _ = time.Parse(time.RFC3339, tag.Target.Date)
			} else if len(tag.Target.Tagger.Date) > 0 {
				published, err = time.Parse(time.RFC3339, tag.Target.Tagger.Date)
				if err != nil {
					published, _ = time.Parse("2006-01-02T15:04:05-07:00", tag.Target.Tagger.Date)
				}
			}
		}
		r := results.NewResult(name, tag.Name, fmt.Sprintf(SourceFormat, name, tag.Name), published)
		rs.AddResult(r)
	}
	return
}

// GetReleases gets a number of releases for a given repo
func (c Provider) GetReleases(name string, max int) (rs *results.ResultSet, err error) {
	names := strings.Split(name, "/")
	query := RepoQuery{
		Query: fmt.Sprintf(RepoQueryFormat, names[0], names[1], max, max),
	}
	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	if err = enc.Encode(&query); err != nil {
		log.Debugf("Failed to encode request: %s\n", err)
		err = results.Unavailable
		return
	}
	// Query the API
	req, _ := http.NewRequest("POST", GraphQLAPI, &buff)
	if key := config.Global.Github.Key; len(key) > 0 {
		req.Header["Authorization"] = []string{"token " + key}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Debugf("Failed to get releases: %s\n", err)
		err = results.Unavailable
		return
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		break
	case 404:
		err = results.NotFound
		return
	default:
		err = results.Unavailable
		return
	}
	// Decode resposne
	dec := json.NewDecoder(resp.Body)
	var rqr RepoQueryResult
	if err = dec.Decode(&rqr); err != nil {
		log.Debugf("Failed to decode response: %s\n", err)
		err = results.Unavailable
		return
	}
	rs = rqr.Convert(name)
	return
}
