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

package providers

import (
	"github.com/DataDrake/cuppa/providers/cpan"
	"github.com/DataDrake/cuppa/providers/git"
	"github.com/DataDrake/cuppa/providers/github"
	"github.com/DataDrake/cuppa/providers/gitlab"
	"github.com/DataDrake/cuppa/providers/gnome"
	"github.com/DataDrake/cuppa/providers/gnu"
	"github.com/DataDrake/cuppa/providers/hackage"
	"github.com/DataDrake/cuppa/providers/html"
	"github.com/DataDrake/cuppa/providers/jetbrains"
	"github.com/DataDrake/cuppa/providers/kde"
	"github.com/DataDrake/cuppa/providers/launchpad"
	"github.com/DataDrake/cuppa/providers/pypi"
	"github.com/DataDrake/cuppa/providers/rubygems"
	"github.com/DataDrake/cuppa/providers/sourceforge"
	"github.com/DataDrake/cuppa/results"
)

// Provider provides a common interface for each of the backend providers
type Provider interface {
	Name() string
	Match(query string) string
	Latest(name string) (*results.Result, error)
	Releases(name string) (*results.ResultSet, error)
}

// All returns a list of all available providers
func All() []Provider {
	return []Provider{
		cpan.Provider{},
		github.Provider{},
		gnome.Provider{},
		gnu.Provider{},
		hackage.Provider{},
		jetbrains.Provider{},
		kde.Provider{},
		launchpad.Provider{},
		pypi.Provider{},
		rubygems.Provider{},
		sourceforge.Provider{},
		git.Provider{},
		gitlab.Provider{},
		html.Provider{},
	}
}
