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

package cpan

import (
	"encoding/json"
	"fmt"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
)

// APIModule holds the name of a Perl Module
type APIModule struct {
	Module string `json:"main_module"`
}

func nameToModule(name string) (module string, err error) {
	// Query the Release API
	resp, err := http.Get(fmt.Sprintf(APIRelease, name))
	if err != nil {
		log.Debugf("Failed to get CPAN module: %s\n", err)
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
	// Decode response
	dec := json.NewDecoder(resp.Body)
	var r APIModule
	if err = dec.Decode(&r); err != nil {
		log.Debugf("Failed to decode response: %s\n", err)
		err = results.Unavailable
		return
	}
	module = r.Module
	return
}
