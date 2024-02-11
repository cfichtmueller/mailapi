// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import "github.com/cfichtmueller/jug"

type Attachment struct {
	ContentType string `json:"contentType"`
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	Data        string `json:"data"`
}

func (a Attachment) Validate(v *jug.Validator) {
	v.RequireNotEmpty(a.ContentType, "content_type_missing").
		RequireNotEmpty(a.Name, "name_missing").
		RequireNotEmpty(a.Filename, "filename_missing").
		RequireNotEmpty(a.Data, "data_missing")
}
