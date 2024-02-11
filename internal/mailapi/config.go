// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mailapi

type Config struct {
	Host   string       `yaml:"host"`
	ApiKey string       `yaml:"apiKey"`
	Smtp   SenderConfig `yaml:"smtp"`
}
