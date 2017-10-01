// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Build constraints copied from go's src/os/dir_unix.go
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package settings

import (
	"os/user"
	"path/filepath"
)

var defaultConfigPaths = []string{
	// This will be prepended by $HOME/.config/webborer.conf
	"/etc/webborer.conf",
}

func init() {
	if usr, err := user.Current(); err == nil {
		path := filepath.Join(usr.HomeDir, ".config", "webborer.conf")
		defaultConfigPaths = append([]string{path}, defaultConfigPaths...)
	}
}
