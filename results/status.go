//
// Copyright © 2017 Bryan T. Meyers <bmeyers@datadrake.com>
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

package results

/*
Status indicates the state of a query upon completion
*/
type Status uint8

const (
	// OK - Query completed successfully, with results
	OK = Status(0)
	// NotFound - Query completed successfully, without results
	NotFound = Status(1)
	// Unavailable - Provider could not be reached
	Unavailable = Status(2)
)
