// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import (
	"strings"
	"testing"
)

func Test_createBoundary(t *testing.T) {
	boundary := createBoundary()
	if strings.Contains(boundary, "=") {
		t.Errorf("boundary contains \"=\": %s", boundary)
	}
}
