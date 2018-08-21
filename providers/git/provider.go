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

package git

import (
    "github.com/DataDrake/cuppa/results"
    git2 "gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/config"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "gopkg.in/src-d/go-git.v4/storage/memory"
    "regexp"
    "strings"
)


// SourceRegex is the regex for Git sources
var SourceRegex = regexp.MustCompile("^(?:git|)?(.+)(?:\\.git)?$")

// Provider provides a common interface for each of the backend providers
type Provider struct {}

func (p Provider) Name() string {
    return "Git"
}

// Latest finds the newest release for a Git package
func (p Provider) Latest(name string) (*results.Result, results.Status){
    repo, err := git2.Init(memory.NewStorage(), nil)
    if err != nil {
        return nil, results.NotFound
    }
    _, err = repo.CreateRemote(&config.RemoteConfig{
        Name: "origin",
        URLs: []string{name},
    })
    if err != nil {
        return nil, results.NotFound
    }
    err = repo.Fetch(&git2.FetchOptions{
        RemoteName: "origin",
        Tags: git2.AllTags,
        Depth: 0,
    })
    if err != nil {
        return nil, results.NotFound
    }
    tags, err := repo.TagObjects()
    if err != nil {
        return nil, results.NotFound
    }
    pieces := strings.Split(name, "/")
    repoName := strings.Split(pieces[len(pieces)-1], ".")[0]
    var r *results.Result
    err = tags.ForEach(func(tag *object.Tag) error {
        r2 := &results.Result{
            Name: repoName,
            Published: tag.Tagger.When,
            Version: tag.Name,
        }
        if r == nil || r.Published.Before(r2.Published) {
            r = r2
        }
        return nil
    })
    repo.DeleteRemote("origin")
    if r == nil {
        return nil, results.NotFound
    }
    return r, results.OK
}

// Match checks to see if this provider can handle this kind of query
func (p Provider) Match(query string) string {
    sm := SourceRegex.FindStringSubmatch(query)
    if len(sm) != 2 {
        return ""
    }
    pieces := strings.Split(sm[1], "|")
    if len(pieces) > 1 {
        return pieces[1]
    }
    return pieces[0]
}

// Releases finds all matching releases for a Git package
func (p Provider) Releases(name string) (*results.ResultSet, results.Status) {
    rs := results.NewResultSet(name)
    repo, err := git2.Init(memory.NewStorage(), nil)
    if err != nil {
        panic(err.Error())
    }
    _, err = repo.CreateRemote(&config.RemoteConfig{
        Name: "origin",
        URLs: []string{name},
    })
    if err != nil {
        panic(err.Error())
    }
    err = repo.Fetch(&git2.FetchOptions{
        RemoteName: "origin",
        Tags: git2.AllTags,
        Depth: 100,
    })
    if err != nil {
        panic(err.Error())
    }
    tags, err := repo.TagObjects()
    if err != nil {
        panic(err.Error())
    }
    pieces := strings.Split(name, "/")
    repoName := strings.Split(pieces[len(pieces)-1], ".")[0]
    err = tags.ForEach(func(tag *object.Tag) error {
        r := &results.Result{
            Name: repoName,
            Published: tag.Tagger.When,
            Version: tag.Name,
        }
        rs.AddResult(r)
        return nil
    })
    repo.DeleteRemote("origin")
    if rs.Len() == 0 {
        return nil, results.NotFound
    }
    return rs, results.OK
}
