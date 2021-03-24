# je4/sshtunnel

This Go package wraps the [crypto/ssh
package](https://godoc.org/golang.org/x/crypto/ssh) with a higher-level API for
building SSH tunnels.

```go
import (
    "flag"
    "github.com/je4/sshtunnel/v2/pkg/sshtunnel"
    "github.com/op/go-logging"
    "os"
    "time"
)


func main() {
	log := logging.MustGetLogger("sshtunnel")
	t, err := sshtunnel.NewSSHTunnel(
		"root",
		"ed25519.priv.openssh",
		&sshtunnel.Endpoint{
			Host: "somehwere.earth",
			Port: 22,
		},
		map[string]*sshtunnel.SourceDestination{
			"postgresql": &sshtunnel.SourceDestination{
				Local: &sshtunnel.Endpoint{
					Host: "localhost",
					Port: 3306,
				},
				Remote: &sshtunnel.Endpoint{
					Host: "db.server.earth",
					Port: 3306,
				},
			},
		},
		log,
	)
	if err != nil {
		log.Errorf("cannot create tunnel - %v", err)
		return
	}
	if err := t.Start(); err != nil {
		log.Errorf("cannot start sshtunnel %v - %v", t.String(), err)
		return
	}
	defer t.Close()
	
	time.Sleep(2 * time.Second)
	
	// [...]
}
```

