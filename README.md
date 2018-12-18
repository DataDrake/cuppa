# cuppa
Comprehensive Upstream Provider Polling Assistant

[![Go Report Card](https://goreportcard.com/badge/github.com/DataDrake/cuppa)](https://goreportcard.com/report/github.com/DataDrake/cuppa) [![license](https://img.shields.io/github/license/DataDrake/cuppa.svg)]() 

## Motivation

As a package maintainer, it's a challenging task to keep track of every new release of a piece of software. Usually this involves subscribing to mailing lists, signing up for notifications from FOSS portals like Github, or even subscribing to news sites. For a distro, this might also mean a repeated effort amongst its package maintainers. The inefficiency and time requirements of such an approach is inevitable. This has led several distros to create their own upstream tracking platforms to automate the process of tracking upstream releases. However, these platforms are often distro specific, leading to further duplication of effort between distros.

## Goals

 * Support as many upstream providers as possible
 * Be completely distro agnostic
 * Extensibility
 * A+ Rating on [Report Card](https://goreportcard.com/report/github.com/DataDrake/cuppa)
 
## Progress

### Supported Providers
* CPAN
* Github (with API Key support)
* GNOME
* Hackage
* Jetbrains
* KDE
* Launchpad
* PyPi
* RubyGems

### Planned Providers
* BitBucket
* FTP
* Git
* GitLab

### Stretch Goal Providers
* RSS
* Sourceforge

Both of these will require some level of scraping to get useful info.

### Unsupported Providers
* NPM
  Completely pointless as this will just pivot to another provider
* Stackage
  Not really in scope for this project and they seem to be missing a web API

## Installation

1. Clone repo and enter its
2. `make`
3. `sudo make install`

## Configuration

Your configuration file must be located at `$HOME/.config/cuppa`

### Github Personal Access Keys

Github limits the number of requests per day for unauthenticated clients. If you would like to get 
around this limitation, you can configure Cuppa to use a Personal Access Key (PAK) by following the
instructions [here](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/#creating-a-token). You do **not** need to enable any OAuth Scopes.

Example:
``` toml
[github]
key = "<personal access key>"
```

## Usage

All `cuppa` commands follow the format:

`cuppa CMD URL`

where CMD is the action to perform and URL is the link to an existing source.

### Commands (CMD)

| CMD      | Alias | Description                                        |
| -------- | ----- | -------------------------------------------------- |
| help     |   ?   | Get help for the other commands.                   |
| latest   |   l   | Get the details for the latest (non-beta) release. |
| quick    |   q   | Get just the new version number and URL if found.  |
| releases |   r   | Get all known previous (non-beta) releases.        |

### Example Sources

| Provider  | URL |
| --------- | --- |
| CPAN      | https://cpan.metacpan.org/authors/id/T/TO/TODDR/IO-1.39.tar.gz |
| Github    | https://github.com/DataDrake/cuppa/archive/v1.0.4.tar.gz |
| GNOME     | https://download.gnome.org/sources/gnome-music/3.28/gnome-music-3.28.2.tar.xz |
| Hackage   | https://hackage.haskell.org/package/mtl-2.2.2/mtl-2.2.2.tar.gz |
| JetBrains | https://download.jetbrains.com/ruby/RubyMine-2017.3.3.tar.gz |
| KDE       | https://download.kde.org/stable/applications/18.12.0/src/akonadi-18.12.0.tar.xz |
| Launchpad | https://launchpad.net/catfish-search/1.4/1.4.4/+download/catfish-1.4.4.tar.gz |
| PyPi      | https://pypi.python.org/packages/2c/a9/69f67f6d5d2fd80ef3d60dc5bef4971d837dc741be0d53295d3aabb5ec7f/pyparted-3.10.7.tar.gz |
| Rubygems  | https://rubygems.org/downloads/sass-3.4.25.gem |

## License
 
Copyright 2016-2018 Bryan T. Meyers <bmeyers@datadrake.com>
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 
