//
// Copyright 2016-2020 Bryan T. Meyers <root@datadrake.com>
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
	"testing"
)

func newVersionTest(t *testing.T, raw string, actual Version) {
	v := NewVersion(raw)
	if len(v) != len(actual) {
		t.Logf("%#v", v)
		t.Errorf("Expected length '%d', found '%d'", len(actual), len(v))
	}
	for i, piece := range actual {
		if v[i] != piece {
			t.Error("Versions do not match")
		}
	}
}

const rawVersion1 = "1.2.3.4"

var version1 = Version{"1", "2", "3", "4"}

func TestNewVersion1(t *testing.T) {
	newVersionTest(t, rawVersion1, version1)
}

const rawVersion2 = "v1.2.3.4."

var version2 = Version{"1", "2", "3", "4"}

func TestNewVersion2(t *testing.T) {
	newVersionTest(t, rawVersion2, version2)
}

const rawVersion3 = "release-1.2.3.5"

var version3 = Version{"1", "2", "3", "5"}

func TestNewVersion3(t *testing.T) {
	newVersionTest(t, rawVersion3, version3)
}

const rawVersion4 = "v1.2-bob-4"

var version4 = Version{"1", "2", "bob", "4"}

func TestNewVersion4(t *testing.T) {
	newVersionTest(t, rawVersion4, version4)
}

const rawVersion5 = "v1.2.3rc1"

var version5 = Version{"1", "2", "3", "rc", "1"}

func TestNewVersion5(t *testing.T) {
	newVersionTest(t, rawVersion5, version5)
}

func TestVersionCompareEqual1(t *testing.T) {
	if c := version1.Compare(version1); c != 0 {
		t.Errorf("Should be equal, found '%d'", c)
	}
}

func TestVersionCompareEqual2(t *testing.T) {
	if c := version5.Compare(version5); c != 0 {
		t.Errorf("Should be equal, found '%d'", c)
	}
}

func TestVersionCompare1(t *testing.T) {
	if version1.Compare(version3) <= 0 {
		t.Error("Should have been less")
	}
	if version3.Compare(version1) >= 0 {
		t.Error("Should have been greater")
	}
}

func TestVersionCompare2(t *testing.T) {
	if version5.Compare(version1) >= 0 {
		t.Error("Should have been less")
	}
	if version1.Compare(version5) <= 0 {
		t.Error("Should have been greater")
	}
}

func TestVersionLess1(t *testing.T) {
	if version5.Less(version5) {
		t.Error("Should not be less: equal")
	}
}

func TestVersionLess2(t *testing.T) {
	if version1.Less(version5) {
		t.Error("Should not be less: greater")
	}
}

func TestVersionLess3(t *testing.T) {
	if !version5.Less(version1) {
		t.Error("Should be less: less")
	}
}
