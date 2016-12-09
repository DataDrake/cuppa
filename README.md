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
* Jetbrains

### Planned Providers
* FTP
* Github
* Hackage
* Launchpad
* RubyGems
* PyPi
* Sourceforge
 
## License
 
Copyright © 2016 Bryan T. Meyers <bmeyers@datadrake.com>
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 