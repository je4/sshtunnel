// Copyright 2021 Jürgen Enge, info-age GmbH, Basel
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/BurntSushi/toml"
	"github.com/goph/emperror"
	"log"
	"net"
	"strconv"
)

type Endpoint struct {
	Host string
	Port int
}

type Forward struct {
	Local  *Endpoint
	Remote *Endpoint
}

type SSHTunnel struct {
	User       string             `toml:"user"`
	PrivateKey string             `toml:"privatekey"`
	Endpoint   *Endpoint          `toml:"endpoint"`
	Forward    map[string]Forward `toml:"forward"`
}

type Config struct {
	Logfile  string               `toml:"logfile"`
	Loglevel string               `toml:"loglevel"`
	Tunnel   map[string]SSHTunnel `toml:"tunnel"`
}

func (e *Endpoint) UnmarshalText(text []byte) error {
	var err error
	var port string
	e.Host, port, err = net.SplitHostPort(string(text))
	if err == nil {
		var longPort int64
		longPort, err = strconv.ParseInt(port, 10, 64)
		if err != nil {
			return emperror.Wrapf(err, "cannot parse port %s of %s", port, string(text))
		}
		e.Port = int(longPort)
	}
	return err
}

func LoadConfig(filepath string) Config {
	var conf Config
	_, err := toml.DecodeFile(filepath, &conf)
	if err != nil {
		log.Fatalln("Error on loading config: ", err)
	}
	return conf
}
