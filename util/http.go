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

package util

import (
	"encoding/json"
	"github.com/DataDrake/cuppa/results"
	log "github.com/DataDrake/waterlog"
	"net/http"
)

// FetchJSON requests from a URL and converts the message body from JSON to a desired type
func FetchJSON(url, kind string, out interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Debugf("Failed to build request: %s\n", err)
		return results.Unavailable
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Debugf("Failed to get %s: %s\n", kind, err)
		return results.Unavailable
	}
	defer resp.Body.Close()
	// Translate Status Code
	switch resp.StatusCode {
	case 200:
		break
	case 404:
		return results.NotFound
	default:
		return results.Unavailable
	}
	// Decode response
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(out); err != nil {
		log.Debugf("Failed to decode response: %s\n", err)
		return results.Unavailable
	}
	return nil
}
