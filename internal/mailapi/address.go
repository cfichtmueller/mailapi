// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import "net/mail"

type Address struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (a Address) toMailAddress() *mail.Address {
	return &mail.Address{Name: a.Name, Address: a.Address}
}
