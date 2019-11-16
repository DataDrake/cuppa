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
    "fmt"
	"github.com/DataDrake/cuppa/results"
	"io"
	"os/exec"
	"strings"
	"time"
    "os"
)

// Provider provides a common interface for each of the backend providers
type Provider struct{}

// Name returns the name of this provider
func (p Provider) Name() string {
	return "Git"
}

// Latest finds the newest release for a Git package
func (p Provider) Latest(name string) (r *results.Result, s results.Status) {
	pieces := strings.Split(name, "/")
	repoName := strings.Split(pieces[len(pieces)-1], ".")[0]
    cwd, _ := os.Getwd()
    err := os.Chdir("/tmp")
    if err != nil {
        s = results.Unavailable
        return
    }
	cmd := exec.Command("git", "clone", "--depth=1", name)
	err = cmd.Run()
    if err == nil {
        err = os.Chdir(fmt.Sprintf("./%s", repoName))
        if err == nil {
        	cmd = exec.Command("git", "fetch", "--tags", "--depth=1")
	        err = cmd.Run()
        }
    }
	buff := new(bytes.Buffer)
	read := bufio.NewReader(buff)
    var line []byte
    var tag string
    var date time.Time
    if err != nil {
        s = results.Unavailable
        goto CLEANUP
    }
    cmd = exec.Command("git", "log", "--tags", "-n 10", "--format='%S %cI'")
	cmd.Stdout = buff
    cmd.Run()
	line, _, err = read.ReadLine()
	for err == nil {
		pieces = strings.Fields(string(line))
		tag = pieces[0]
        date, _ = time.Parse("2006-01-02T15:04:05-07:00", pieces[1])
		r = results.NewResult(repoName, tag, "git|"+name, date)
		line, _, err = read.ReadLine()
	}
	if err != io.EOF || r == nil {
		s = results.NotFound
	} else {
        s = results.OK
    }
CLEANUP:
    os.RemoveAll(fmt.Sprintf("/tmp/%s", repoName))
    os.Chdir(cwd)
	return
}

// Match checks to see if this provider can handle this kind of query
func (p Provider) Match(query string) string {
	if strings.HasPrefix(query, "git|") || strings.HasSuffix(query, ".git") {
		pieces := strings.Split(query, "|")
		if len(pieces) > 1 {
			return pieces[1]
		}
		return pieces[0]
	}
	return ""
}

// Releases finds all matching releases for a Git package
func (p Provider) Releases(name string) (*results.ResultSet, results.Status) {
	cmd := exec.Command("git", "ls-remote", "--tags", "--sort='-*authordate'", name)
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
			r = results.NewResult(repoName, tag, "git|"+name, time.Time{})
			rs.AddResult(r)
		}
		line, _, err = read.ReadLine()
	}
	if err != io.EOF || rs.Len() == 0 {
		return nil, results.NotFound
	}
	return rs, results.OK
}
