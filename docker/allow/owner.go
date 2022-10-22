package allow

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/kassisol/hbm/docker/allow/types"
	policyobj "github.com/kassisol/hbm/object/policy"
	"github.com/kassisol/hbm/version"
	log "github.com/sirupsen/logrus"
)

func ContainerOwner(req authorization.Request, config *types.Config) *types.AllowResult {
	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		log.WithFields(log.Fields{
			"version": version.Version,
		}).Fatal(err)
	}
	defer p.End()

	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return &types.AllowResult{Allow: false}
	}

	ts := strings.Trim(u.Path, "/")
	up := strings.Split(ts, "/") // api version / type / id
	if len(up) < 3 {
		return &types.AllowResult{Allow: false}
	}
	if up[1] != "containers" {
		return &types.AllowResult{Allow: false}
	}

	if !p.ValidateOwner(config.Username, "containers", up[2]) {
		return &types.AllowResult{
			Allow: false,
			Msg: map[string]string{
				"text": fmt.Sprintf("Container %s is not owned by %s", up[2], config.Username),
			},
		}
	}

	return &types.AllowResult{Allow: true}
}
