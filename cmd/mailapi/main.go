// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/mailapi/internal/mailapi"
	"github.com/cfichtmueller/mailapi/internal/util"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	config := mailapi.Config{}

	if len(os.Args) > 1 {
		b, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}
	}

	config.Host = util.Nvl(util.Nvl(os.Getenv("HOST"), config.Host), "127.0.0.1:8000")
	config.ApiKey = util.Nvl(os.Getenv("API_KEY"), config.ApiKey)
	if config.ApiKey == "" {
		_, _ = fmt.Fprintln(os.Stderr, "API_KEY_MISSING")
		os.Exit(1)
	}

	engine := jug.New()

	smtpPort, err := strconv.Atoi(util.Nvl(os.Getenv("SMTP_PORT"), "25"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "invalid SMTP_PORT")
		os.Exit(1)
	}
	if config.Smtp.Port == 0 {
		config.Smtp.Port = smtpPort
	}
	config.Smtp.Host = util.Nvl(util.Nvl(os.Getenv("SMTP_HOST"), config.Smtp.Host), "localhost")
	if !config.Smtp.Tls {
		config.Smtp.Tls = util.Nvl(os.Getenv("SMTP_TLS"), "false") == "true"
	}
	config.Smtp.Username = util.Nvl(os.Getenv("SMTP_USERNAME"), config.Smtp.Username)
	config.Smtp.Password = util.Nvl(os.Getenv("SMTP_PASSWORD"), config.Smtp.Password)

	sender := mailapi.NewSender(config.Smtp)

	apiGroup := engine.Group("/api", func(c jug.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.RespondUnauthorized(nil)
			c.Abort()
			return
		}
		token := authHeader[7:]
		if token != config.ApiKey {
			c.RespondUnauthorized(nil)
			c.Abort()
			return
		}
	})

	apiGroup.POST("/send", func(c jug.Context) {
		var req mailapi.Email
		if !c.MustBindJSON(&req) {
			return
		}

		if err := sender.Send(req); err != nil {
			c.RespondInternalServerError(ErrorResponse{Error: err.Error()})
			c.Abort()
			return
		}

		c.RespondNoContent()
	})

	engine.ExpandMethods()

	log.Printf("Starting Server on %s", config.Host)
	if err := engine.Run(config.Host); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}
