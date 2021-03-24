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
	"flag"
	"github.com/je4/sshtunnel/pkg/sshtunnel"
	"github.com/op/go-logging"
	"os"
	"time"
)

var _logFormat = logging.MustStringFormatter(
	`%{time:2006-01-02T15:04:05.000} %{module}::%{shortfunc} [%{shortfile}] > %{level:.5s} - %{message}`,
)

func CreateLogger(module string, logfile string, loglevel string) (log *logging.Logger, lf *os.File) {
	log = logging.MustGetLogger(module)
	var err error
	if logfile != "" {
		lf, err = os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Errorf("Cannot open logfile %v: %v", logfile, err)
		}
		//defer lf.CloseInternal()

	} else {
		lf = os.Stderr
	}
	backend := logging.NewLogBackend(lf, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.GetLevel(loglevel), "")

	logging.SetFormatter(_logFormat)
	logging.SetBackend(backendLeveled)

	return
}

func main() {
	cfgfile := flag.String("config", "./sshtunnel.toml", "location of config file")
	flag.Parse()
	config := LoadConfig(*cfgfile)

	// create logger instance
	log, lf := CreateLogger("memostream", config.Logfile, config.Loglevel)
	defer lf.Close()

	log.Info("initializing ssh tunnels...")

	for name, tunnel := range config.Tunnel {
		log.Infof("starting tunnel %s", name)

		forwards := make(map[string]*sshtunnel.SourceDestination)
		for fwname, fw := range tunnel.Forward {
			forwards[fwname] = &sshtunnel.SourceDestination{
				Local: &sshtunnel.Endpoint{
					Host: fw.Local.Host,
					Port: fw.Local.Port,
				},
				Remote: &sshtunnel.Endpoint{
					Host: fw.Remote.Host,
					Port: fw.Remote.Port,
				},
			}
		}

		t, err := sshtunnel.NewSSHTunnel(
			tunnel.User,
			tunnel.PrivateKey,
			&sshtunnel.Endpoint{
				Host: tunnel.Endpoint.Host,
				Port: tunnel.Endpoint.Port,
			},
			forwards,
			log,
		)
		if err != nil {
			log.Errorf("cannot create tunnel %v@%v:%v - %v", tunnel.User, tunnel.Endpoint.Host, tunnel.Endpoint.Port, err)
			return
		}
		if err := t.Start(); err != nil {
			log.Errorf("cannot create sshtunnel %v - %v", t.String(), err)
			return
		}
		defer t.Close()
	}

	// do something interesting or just wait...
	time.Sleep(2 * time.Second)
}
