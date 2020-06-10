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

func getcontainerid(p string) (bool, string) {
	ts := strings.Trim(p, "/")
	up := strings.Split(ts, "/") // api version / type / id
	if len(up) < 3 {
		return false, ""
	}
	if up[1] != "containers" {
		return false, ""
	}

	return true, up[2]
}

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

	corr, up := getcontainerid(u.Path)
	if corr == false {
		return ar
	}
	
	ar.Allow = p.ValidateOwner(config.Username, "containers", up)
	if !ar.Allow {
		ar.Error = "you are not the owner of the container"
	}
	
	return ar
}

/*func RemoveIfOwner(req authorization.Request, config *types.Config) *types.AllowResult {
	ar := &ContainerOwner(req, config)
	if ar.Allow == false {
		return ar
	}

	corr, up := getcontainerid(u.Path)
	if corr == false {
		return ar
	}
	
	return ar
}*/
