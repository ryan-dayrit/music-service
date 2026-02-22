//go:build integration

package testutil

import (
	_ "embed"
)

//go:embed testdata/integration_init.sql
var InitSQL string
