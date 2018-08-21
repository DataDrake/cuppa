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
    "bufio"
    "bytes"
    "github.com/DataDrake/cuppa/results"
    "io"
    "os/exec"
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
    cmd := exec.Command("git", "ls-remote", "--tags", name)
    buff := new(bytes.Buffer)
    cmd.Stdout = buff
    read := bufio.NewReader(buff)
    err := cmd.Run()
    line, _, err := read.ReadLine()
    var r *results.Result
    pieces := strings.Split(name, "/")
    repoName := strings.Split(pieces[len(pieces)-1], ".")[0]
    for err == nil {
        pieces := strings.Split(string(line), "/")
        tag := pieces[0]
        if len(pieces) > 1 {
            tag = pieces[len(pieces)-1]
        }
        if !strings.HasSuffix(tag, "{}") {
            r = &results.Result{
                Name: repoName,
                Version: tag,
                Location: "git|" + name,
            }
        }
        line, _, err = read.ReadLine()
    }
    if err != io.EOF || r == nil{
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
    cmd := exec.Command("git", "ls-remote", "--tags", name)
    buff := new(bytes.Buffer)
    cmd.Stdout = buff
    read := bufio.NewReader(buff)
    err := cmd.Run()
    line, _, err := read.ReadLine()
    var r *results.Result
    pieces := strings.Split(name, "/")
    repoName := strings.Split(pieces[len(pieces)-1], ".")[0]
    rs := results.NewResultSet(repoName)
    for err == nil {
        pieces := strings.Split(string(line), "/")
        tag := pieces[0]
        if len(pieces) > 1 {
            tag = pieces[len(pieces)-1]
        }
        if !strings.HasSuffix(tag, "{}") {
            r = &results.Result{
                Name: repoName,
                Version: tag,
                Location: "git|" + name,
            }
            rs.AddResult(r)
        }
        line, _, err = read.ReadLine()
    }
    if err != io.EOF || rs.Len() == 0 {
        return nil, results.NotFound
    }
    return rs, results.OK
}
