// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func createBoundary() string {
	b := make([]byte, 15)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	enc := base64.StdEncoding.EncodeToString(b)
	return "-------" + strings.ToUpper(enc)
}
