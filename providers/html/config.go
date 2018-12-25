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

package html

import (
    "encoding/xml"
    "github.com/DataDrake/cuppa/results"
    "io"
	"regexp"
    "strings"
    "time"
)

var ArchiveRegex = regexp.MustCompile("^(.+)-(.*)\\.(?:tar\\.[^.]+|zip)$")

type DownloadList struct {
    Name xml.Name
    Entry []struct {
        Columns []struct{
            XML string `xml:",innerxml"`
            Raw string `xml:",chardata"`
        } `xml:"td"`
    } `xml:"body>table>tr"`
}

type LocationConfig struct {
    Index int
    XML   bool
    Pattern *regexp.Regexp
}

type TimeConfig struct {
    Index int
    Layout string
}

type Config struct {
    Location LocationConfig
    Modified TimeConfig
}

func (c Config) Parse(name, path string, in io.Reader) (rs *results.ResultSet, err error) {
    dec := xml.NewDecoder(in)
    dec.Strict = false
    dec.AutoClose = xml.HTMLAutoClose
    dec.Entity = xml.HTMLEntity
    var list DownloadList
    err = dec.Decode(&list)
    if err != nil {
        return
    }
    rs = results.NewResultSet(name)
    for _, entry := range list.Entry {
        if len(entry.Columns) >= c.Location.Index && len(entry.Columns) >= c.Modified.Index {
            var loc string
            if c.Location.XML {
                loc = entry.Columns[c.Location.Index].XML
                sm := c.Location.Pattern.FindStringSubmatch(loc)
                if len(sm) != 2 {
                    continue
                }
                loc = sm[1]
            } else {
                loc = entry.Columns[c.Location.Index].Raw
            }
            sm := ArchiveRegex.FindStringSubmatch(loc)
            if len(sm) != 3 {
                continue
            }
            n := sm[1]
            version := sm[2]
            if n != name {
                continue
            }
            mod, e := time.Parse(c.Modified.Layout, strings.TrimSpace(entry.Columns[c.Modified.Index].Raw))
            if e != nil {
                continue
            }
            r := &results.Result {
                Name: n,
                Version: version,
                Published: mod,
                Location: path + loc,
            }
            rs.AddResult(r)
        }
    }
    err = nil
    return
}
