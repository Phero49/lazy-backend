package io

import (
	"os"
	"path"
)


var home, _ = os.UserHomeDir()
// AppDataDir should be defined elsewhere in your package;
// declared here as a placeholder for compilation completeness.
// path.Join(home, ".lazy-backend")
var APP_FOLDER = path.Join(home, ".lazy-backend")
  