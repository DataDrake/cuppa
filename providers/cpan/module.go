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

package cpan

import (
	"fmt"
	"github.com/DataDrake/cuppa/util"
)

// APIModule holds the name of a Perl Module
type APIModule struct {
	Module string `json:"main_module"`
}

func nameToModule(name string) (module string, err error) {
	// Query the Release API
	url := fmt.Sprintf(APIRelease, name)
	var r APIModule
	if err = util.FetchJSON(url, "CPAN module", &r); err == nil {
		module = r.Module
	}
	return
}
