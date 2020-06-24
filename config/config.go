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

package config

import (
	"github.com/BurntSushi/toml"
	log "github.com/DataDrake/waterlog"
	"os/user"
	"path/filepath"
)

// Config is the configuration for cuppa
type Config struct {
	Github struct {
		Key string `toml:"key"`
	} `toml:"github"`
}

// Global is the config for all of cuppa at runtime
var Global Config

// init loads the config if it exists
func init() {
	user, err := user.Current()
	if err != nil {
		return
	}
	configPath := filepath.Join(user.HomeDir, ".config", "cuppa")
	if _, err = toml.DecodeFile(configPath, &Global); err != nil {
		log.Fatalf("Failed to read config, reason: '%s'\n", err)
	}
}
