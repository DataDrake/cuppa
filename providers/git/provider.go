//
// Copyright 2016-2020 Bryan T. Meyers <root@datadrake.com>
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
	"os"
	"os/exec"
	"strings"
	"time"
)

// Provider provides a common interface for each of the backend providers
type Provider struct{}

// Name returns the name of this provider
func (p Provider) Name() string {
	return "Git"
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

// Latest finds the newest release for a Git package
func (p Provider) Latest(name string) (r *results.Result, err error) {
	rs, err := p.Releases(name)
	if err == nil {
		r = rs.Last()
	}
	return
}

// Releases finds all matching releases for a Git package
func (p Provider) Releases(name string) (rs *results.ResultSet, err error) {
	pieces := strings.Split(name, "/")
	repoName := strings.Split(pieces[len(pieces)-1], ".")[0]
	tmp := fmt.Sprintf("/tmp/%s", repoName)
	defer os.RemoveAll(tmp)
	// Shallow clone repo to temp directory
	cmd := exec.Command("git", "clone", "--depth=1", name)
	cmd.Dir = "/tmp"
	if err := cmd.Run(); err == nil {
		// Fetch tags from remote
		cmd = exec.Command("git", "fetch", "--tags", "--depth=1")
		cmd.Dir = tmp
		err = cmd.Run()
	}
	if err != nil {
		err = results.Unavailable
		return
	}
	// Read git tags
	var buff bytes.Buffer
	read := bufio.NewReader(&buff)
	var tag string
	var date time.Time
	cmd = exec.Command("git", "log", "--tags", "-n 10", "--format='%S %cI'")
	cmd.Dir = tmp
	cmd.Stdout = &buff
	cmd.Run()
	// Convert tags to releases
	rs = results.NewResultSet(name)
	line, _, err := read.ReadLine()
	for err == nil {
		pieces = strings.Fields(string(line))
		tag = pieces[0]
		date, _ = time.Parse("2006-01-02T15:04:05-07:00", pieces[1])
		r := results.NewResult(repoName, tag, "git|"+name, date)
		rs.AddResult(r)
		line, _, err = read.ReadLine()
	}
	if err != io.EOF || rs.Len() == 0 {
		err = results.NotFound
		return
	}
	err = nil
	return
}
