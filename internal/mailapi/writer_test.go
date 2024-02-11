// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"testing"
)

func TestWriter_writeTextMail(t *testing.T) {
	buf := new(bytes.Buffer)
	writeTextMail(buf, "Hello World")

	fmt.Println(buf.String())
}

func TestWriter_writeHtmlMail(t *testing.T) {
	buf := new(bytes.Buffer)
	writeHtmlMail(buf, "<p>Hello World</p>")

	fmt.Println(buf.String())
}

func TestWriter_writeAlternativeMail(t *testing.T) {
	buf := new(bytes.Buffer)
	writeAlternativeMail(buf, Email{
		TextContent: "Hello Bob",
		HtmlContent: "<p>Hello,</p><p>Bob</p>",
	})

	fmt.Println(buf.String())
}

func TestWriter_writeAttachmentMail(t *testing.T) {
	buf := new(bytes.Buffer)
	writeAttachmentMail(buf, Email{
		TextContent: "Hello Bob",
		HtmlContent: "<p>Hello,</p><p>Bob</p>",
		Attachments: []Attachment{
			{
				Name:        "invite.ics",
				ContentType: "application/ics",
				Filename:    "invite.ics",
				Data:        "ABCDEFHI",
			},

			{Name: "logo", ContentType: "image/png", Filename: "logo.png", Data: "fdc82b"},
		},
	})

	fmt.Println(buf.String())
}

func TestMailAttachment(t *testing.T) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	buf.WriteString("Content-Type: multipart/mixed; boundary=\"")
	buf.WriteString(w.Boundary())
	buf.WriteString("\"\r\n")

	buf2 := new(bytes.Buffer)
	w2 := multipart.NewWriter(buf2)

	c, _ := w.CreatePart(map[string][]string{
		"Content-Type": {"multipart/alternative; boundary=\"" + w2.Boundary() + "\""},
	})

	p1, _ := w2.CreatePart(map[string][]string{
		"Content-Type": {"text/plain"},
	})

	p1.Write([]byte("Hello World"))

	p2, _ := w2.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/html; charset=\"utf-8\""},
	})

	p2.Write([]byte("<p>Hello World</p>"))

	w2.Close()

	c.Write(buf2.Bytes())

	c2, _ := w.CreatePart(map[string][]string{
		"Content-Type": {"application/ics"},
	})

	c2.Write([]byte("somebase64content"))

	w.Close()

	fmt.Println(buf.String())
}
