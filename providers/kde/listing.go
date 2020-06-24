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

package kde

import (
	"compress/bzip2"
	log "github.com/DataDrake/waterlog"
	"io/ioutil"
	"net/http"
)

const ListingURL = "https://download.kde.org/ls-lR.bz2"

var listing []byte

func getListing() {
	// Query the API
	resp, err := http.Get(ListingURL)
	if err != nil {
		log.Debugf("Failed to get listing: %s\n", err)
		return
	}
	defer resp.Body.Close()
	// Translate Status Code
	if resp.StatusCode != 200 {
		return
	}
	body := bzip2.NewReader(resp.Body)
	listing, _ = ioutil.ReadAll(body)
}
