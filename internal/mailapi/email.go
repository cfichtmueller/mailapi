// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import "github.com/cfichtmueller/jug"

type Email struct {
	From        Address      `json:"from"`
	To          Address      `json:"to"`
	Cc          []Address    `json:"cc"`
	Bcc         []Address    `json:"bcc"`
	Subject     string       `json:"subject"`
	ContentType string       `json:"contentType"`
	Content     string       `json:"content"`
	TextContent string       `json:"textContent"`
	HtmlContent string       `json:"htmlContent"`
	Attachments []Attachment `json:"attachments"`
}

func (e Email) Validate() error {
	v := jug.NewValidator().
		RequireNotEmpty(e.From.Address, "from_missing").
		RequireNotEmpty(e.To.Address, "to_missing").
		RequireNotEmpty(e.Subject, "subject_missing").
		V(func(v *jug.Validator) {
			if e.Content != "" {
				v.RequireNotEmpty(e.ContentType, "content_type_missing")
				return
			}
			v.Require(e.TextContent != "" || e.HtmlContent != "", "no_content_given")
		})

	for _, a := range e.Attachments {
		v.V(a.Validate)
	}

	return v.Validate()
}

func (e Email) HasAttachments() bool {
	return len(e.Attachments) > 0
}

func (e Email) IsTextMail() bool {
	return e.Content == "" && e.TextContent != "" && e.HtmlContent == ""
}

func (e Email) IsHtmlMail() bool {
	return e.Content == "" && e.TextContent == "" && e.HtmlContent != ""
}

func (e Email) IsAlternativeMail() bool {
	return e.Content == "" && e.TextContent != "" && e.HtmlContent != ""
}
