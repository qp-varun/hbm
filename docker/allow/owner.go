package allow

import (
	"strings"
	"net/url"
	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/kassisol/hbm/docker/allow/types"
	"github.com/juliengk/go-log"
	logdriver "github.com/juliengk/go-log/driver"
	policyobj "github.com/kassisol/hbm/object/policy"
	"github.com/kassisol/hbm/version"
)

func ContainerOwner(req authorization.Request, config *types.Config) *types.AllowResult {
	ar := &types.AllowResult{Allow: false}

	l, _ := log.NewDriver("standard", nil)

        p, err := policyobj.New("sqlite", config.AppPath)
        if err != nil {
                l.WithFields(logdriver.Fields{
                        "storagedriver": "sqlite",
                        "logdriver":     "standard",
                        "version":       version.Version,
                }).Fatal(err)
        }
        defer p.End()

	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return ar
	}

	ts := strings.Trim(u.Path, "/")
	up := strings.Split(ts, "/") // api version / type / id
	if len(up) < 3 {
		return ar
	}
	if up[1] != "containers" {
		return ar
	}

	ar.Allow = p.ValidateOwner(config.Username, "containers", up[2])

	return ar
}
