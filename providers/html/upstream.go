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
	"github.com/DataDrake/cuppa/results"
    "io"
	"regexp"
)

type Upstream struct {
    Name        string
    HostPattern *regexp.Regexp
    Conf        Config
}

func (u Upstream) Match(path string) string {
    sm := u.HostPattern.FindStringSubmatch(path)
    if len(sm) == 0 {
        return ""
    }
    return sm[0]
}

func (u Upstream) Parse(name string, in io.Reader) (*results.ResultSet, error){
    sm := u.HostPattern.FindStringSubmatch(name)
    return u.Conf.Parse(sm[2],sm[1], in)
}
