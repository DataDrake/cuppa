//
// Copyright 2018 Bryan T. Meyers <bmeyers@datadrake.com>
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

package version

import (
	"strconv"
	"strings"
	"unicode"
)

// Version is a record of a new version for a single source
type Version []string

func splitDigit(raw string) []string {
	pieces := make([]string, 0)
	if len(raw) == 0 {
		return pieces
	}
	i := 0
	for _, char := range raw {
		if !unicode.IsDigit(char) {
			break
		}
		i++
	}
	if i > 0 {
		pieces = append(pieces, raw[:i])
	}
	if i < len(raw) {
		pieces = append(pieces, splitChar(raw[i:])...)
	}
	return pieces
}

func splitChar(raw string) []string {
	pieces := make([]string, 0)
	if len(raw) == 0 {
		return pieces
	}
	i := 0
	for _, char := range raw {
		if unicode.IsDigit(char) {
			break
		}
		i++
	}
	if i > 0 {
		pieces = append(pieces, raw[:i])
	}
	if i < len(raw) {
		pieces = append(pieces, splitDigit(raw[i:])...)
	}
	return pieces
}

// NewVersion creates a new version by parsing from a string
func NewVersion(raw string) Version {
	dots := strings.Split(raw, ".")
	dashes := make([]string, 0)
	for _, dot := range dots {
		dashes = append(dashes, strings.Split(dot, "-")...)
	}
	unclean := make([]string, 0)
	for _, dash := range dashes {
		unclean = append(unclean, strings.Split(dash, "_")...)
	}
	pieces := make([]string, 0)
	for _, u := range unclean {
		if u != "" {
			pieces = append(pieces, u)
		}
	}
	if len(pieces) == 0 {
		return []string{"N/A"}
	}
	v := make(Version, 0)
	started := false
	for _, piece := range pieces {
		parts := splitChar(piece)
		for _, part := range parts {
			if !started && !unicode.IsDigit(rune(part[0])) {
				continue
			}
			started = true
			v = append(v, part)
		}
	}
	if len(v) == 0 {
		v = append(v, "N/A")
	} else {
        // Strip trailing words
        i := len(v)-1
        for i >= 0 && !unicode.IsDigit(rune(v[i][0])) {
            i--
        }
        v = v[:i+1]
    }
	return v
}

// Compare allows to version nubmers to be compared to see which is newer (higher)
func (v Version) Compare(old Version) int {
	result := 0
	for i, piece := range v {
		if len(old) == i {
			return -1
		}
		if old[i] == piece {
			continue
		}
		curr, e1 := strconv.Atoi(piece)
		prev, e2 := strconv.Atoi(old[i])
		if e1 != nil && e2 != nil {
			goto HARD
		}
		if e1 != nil {
			return -1
		}
		if e2 != nil {
			return 1
		}
		result = prev - curr
		goto CHECK
	HARD:
		result = strings.Compare(piece, old[i])
	CHECK:
		if result != 0 {
			return result
		}
	}
	return result
}

// Less checks if this version is less than another
func (v Version) Less(other Version) bool {
	return v.Compare(other) < 0
}
