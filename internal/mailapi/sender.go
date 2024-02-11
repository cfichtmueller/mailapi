// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
)

type SenderConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Tls      bool   `yaml:"tls"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c SenderConfig) Address() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func (c SenderConfig) Auth() smtp.Auth {
	if c.Username == "" || c.Password == "" {
		return nil
	}
	return smtp.PlainAuth("", c.Username, c.Password, c.Host)
}

type Sender struct {
	config SenderConfig
}

func NewSender(config SenderConfig) *Sender {
	return &Sender{config: config}
}

func (s *Sender) Send(m Email) error {
	buf := new(bytes.Buffer)

	headers := make(map[string]string)
	headers["From"] = m.From.toMailAddress().String()
	headers["To"] = m.To.toMailAddress().String()
	headers["Subject"] = m.Subject

	writeHeaders(buf, headers)

	if m.HasAttachments() {
		writeAttachmentMail(buf, m)
	} else if m.IsTextMail() {
		writeTextMail(buf, m.TextContent)
	} else if m.IsHtmlMail() {
		writeHtmlMail(buf, m.HtmlContent)
	} else if m.IsAlternativeMail() {
		writeAlternativeMail(buf, m)
	}

	conn, err := s.createConnection()
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("unable to create SMTP client: %v", err)
	}

	if auth := s.config.Auth(); auth != nil {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("unable to authenticate at SMTP server: %v", err)
		}
	}

	if err = c.Mail(m.From.Address); err != nil {
		return err
	}

	if err = c.Rcpt(m.To.Address); err != nil {
		return err
	}

	d, err := c.Data()
	if err != nil {
		return err
	}

	_, err = d.Write(buf.Bytes())
	if err != nil {
		return err
	}

	if err := d.Close(); err != nil {
		return err
	}

	if err := c.Close(); err != nil {
		return err
	}

	return nil
}

func (s *Sender) createConnection() (net.Conn, error) {
	if !s.config.Tls {
		conn, err := net.Dial("tcp", s.config.Address())
		if err != nil {
			return nil, fmt.Errorf("unable to dial server: %v", err)
		}
		return conn, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.config.Host,
	}

	conn, err := tls.Dial("tcp", s.config.Address(), tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to establish tls connection: %v", err)
	}
	return conn, nil
}
