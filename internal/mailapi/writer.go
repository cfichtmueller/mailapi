// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import (
	"bytes"
	"io"
	"mime/multipart"
)

const newline = "\r\n"

var (
	textHeaders = map[string][]string{
		"Content-Type":              {"text/plain; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
		"Content-Disposition":       {"inline"},
	}
	htmlHeaders = map[string][]string{
		"Content-Type":              {"text/html; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
		"Content-Disposition":       {"inline"},
	}
)

func writeHeaders(buf *bytes.Buffer, headers map[string]string) {
	for k, v := range headers {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v)
		buf.WriteString(newline)
	}
}

func writeAttachmentMail(buf *bytes.Buffer, e Email) {
	outer := multipart.NewWriter(buf)

	buf.WriteString("Content-Type: multipart/mixed; boundary=\"")
	buf.WriteString(outer.Boundary())
	buf.WriteString("\"")
	buf.WriteString(newline)

	var headers map[string][]string
	mpp := multipart.NewWriter(buf)

	if e.IsTextMail() {
		headers = textHeaders
	} else if e.IsHtmlMail() {
		headers = htmlHeaders
	} else {
		headers = map[string][]string{
			"Content-Type": {"multipart/alternative; boundary=\"" + mpp.Boundary() + "\""},
		}
	}

	p, _ := outer.CreatePart(headers)

	if e.IsTextMail() {
		p.Write([]byte(e.TextContent + newline))
	} else if e.IsHtmlMail() {
		p.Write([]byte(e.HtmlContent))
	} else if e.IsAlternativeMail() {
		tp, _ := mpp.CreatePart(textHeaders)
		tp.Write([]byte(e.TextContent))
		hp, _ := mpp.CreatePart(htmlHeaders)
		hp.Write([]byte(e.HtmlContent))
		mpp.Close()
	}

	for _, a := range e.Attachments {
		apart, _ := outer.CreatePart(map[string][]string{
			"Content-Type":              {a.ContentType + ";name=\"" + a.Name + "\""},
			"Content-Disposition":       {"attachment; filename=\"" + a.Filename + "\""},
			"Content-Transfer-Encoding": {"base64"},
		})
		apart.Write([]byte(a.Data))
	}

	outer.Close()
}

func writeTextMail(w io.Writer, text string) {
	w.Write([]byte("Content-Type: text/plain; charset=\"utf-8\""))
	w.Write([]byte(newline))
	w.Write([]byte("Content-Transfer-Encoding: quoted-printable"))
	w.Write([]byte(newline))
	w.Write([]byte("Content-Disposition: inline"))
	w.Write([]byte(newline))
	w.Write([]byte(newline))
	w.Write([]byte(text))
	w.Write([]byte(newline))
}

func writeHtmlMail(w io.Writer, text string) {
	w.Write([]byte("Content-Type: text/html; charset=\"utf-8\""))
	w.Write([]byte(newline))
	w.Write([]byte("Content-Transfer-Encoding: quoted-printable"))
	w.Write([]byte(newline))
	w.Write([]byte("Content-Disposition: inline"))
	w.Write([]byte(newline))
	w.Write([]byte(newline))
	w.Write([]byte(text))
	w.Write([]byte(newline))
}

func writeAlternativeMail(w io.Writer, e Email) {
	mp := multipart.NewWriter(w)

	w.Write([]byte("Content-Type: multipart/alternative; boundary=\""))
	w.Write([]byte(mp.Boundary()))
	w.Write([]byte("\""))
	w.Write([]byte(newline))
	w.Write([]byte(newline))

	textPart, _ := mp.CreatePart(map[string][]string{
		"Content-Type":              {"text/plain; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
		"Content-Disposition":       {"inline"},
	})
	textPart.Write([]byte(e.TextContent))

	htmlPart, _ := mp.CreatePart(map[string][]string{
		"Content-Type":              {"text/html; charset=\"utf-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
		"Content-Disposition":       {"inline"},
	})
	htmlPart.Write([]byte(e.HtmlContent))

	mp.Close()
}
