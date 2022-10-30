package plugin

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/kassisol/hbm/pkg/uri"
	"github.com/kassisol/hbm/storage"
	log "github.com/sirupsen/logrus"
)

type plugin struct {
	appPath       string
	skipEndpoints []*regexp.Regexp
}

func stringInRegexpSlice(s string, regexps []*regexp.Regexp) bool {
	for _, re := range regexps {
		if re.MatchString(s) {
			return true
		}
	}

	return false
}

func NewPlugin(appPath string) (*plugin, error) {
	p := plugin{
		appPath: appPath,
		skipEndpoints: []*regexp.Regexp{
			regexp.MustCompile(`^/_ping`),
			regexp.MustCompile(`^/distribution/(.+)/json`),
		},
	}

	return &p, nil
}

func (p *plugin) AuthZReq(req authorization.Request) authorization.Response {
	uriinfo, err := uri.GetURIInfo(req)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	if req.RequestMethod == "OPTIONS" || stringInRegexpSlice(uriinfo.Path, p.skipEndpoints) {
		return authorization.Response{Allow: true}
	}

	log.WithFields(log.Fields{
		"user":               req.User,
		"RequestMethod":      req.RequestMethod,
		"RequestURI":         req.RequestURI,
		"RequestBody":        string(req.RequestBody),
		"ResponseStatusCode": req.ResponseStatusCode,
		"ResponseBody":       string(req.ResponseBody),
	}).Debug("AuthZReq")

	a, err := NewApi(&uriinfo, p.appPath)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	r := a.Allow(req)
	if r.Error != "" {
		return authorization.Response{Err: r.Error}
	}
	if !r.Allow {
		return authorization.Response{Msg: r.Msg["text"]}
	}

	return authorization.Response{Allow: true}
}

func (p *plugin) AuthZRes(req authorization.Request) authorization.Response {
	uriinfo, err := uri.GetURIInfo(req)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	if stringInRegexpSlice(uriinfo.Path, p.skipEndpoints) {
		return authorization.Response{Allow: true}
	}

	log.WithFields(log.Fields{
		"user":               req.User,
		"RequestMethod":      req.RequestMethod,
		"RequestURI":         req.RequestURI,
		"RequestBody":        string(req.RequestBody),
		"ResponseStatusCode": req.ResponseStatusCode,
		"ResponseBody":       string(req.ResponseBody),
	}).Debug("AuthZRes")

	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return authorization.Response{Allow: true, Msg: err.Error()}
	}

	cname := u.Query().Get("name")
	if p.isCreateContainer(req, u) {
		err = p.setContainerOwner(cname, req)
		if err != nil {
			return authorization.Response{Allow: false, Msg: err.Error()}
		}
	}

	if p.isRemoveContainer(req, u) {
		err = p.removeContainerOwner(cname, req)
		if err != nil {
			return authorization.Response{Allow: false, Msg: err.Error()}
		}
	}

	return authorization.Response{Allow: true}
}

func (p *plugin) isCreateContainer(req authorization.Request, u *url.URL) bool {
	if req.ResponseStatusCode != 201 {
		return false
	}

	avm := regexp.MustCompile(`^/v\d+\.\d+/containers/create`)
	if avm.MatchString(u.Path) || u.Path == "/containers/create" {
		return true
	}

	return false
}

func (p *plugin) isRemoveContainer(req authorization.Request, u *url.URL) bool {
	if req.ResponseStatusCode != 204 {
		return false
	}

	avm := regexp.MustCompile(`^/v\d+\.\d+/containers/[^/]+`)
	if avm.MatchString(u.Path) && req.RequestMethod == "DELETE" {
		return true
	}

	return false
}

func (p *plugin) setContainerOwner(cname string, req authorization.Request) error {
	username := req.User
	if username == "" {
		username = "root"
	}

	s, err := storage.NewDriver("sqlite", p.appPath)
	if err != nil {
		return err
	}
	defer s.End()

	var rjson struct {
		Id string
	}
	err = json.Unmarshal(req.ResponseBody, &rjson)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"username": username,
		"cname":    cname,
		"rJson.Id": rjson.Id,
	}).Debug("Calling SetContainerOwner() with")

	err = s.SetContainerOwner(username, cname, rjson.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *plugin) removeContainerOwner(cname string, req authorization.Request) error {
	username := req.User
	if username == "" {
		username = "root"
	}

	s, err := storage.NewDriver("sqlite", p.appPath)
	if err != nil {
		return err
	}
	defer s.End()

	log.Debug(req)

	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return err
	}

	ts := strings.Trim(u.Path, "/")
	up := strings.Split(ts, "/") // api version / type / id
	if len(up) < 3 {
		return nil
	}
	if up[1] != "containers" {
		return nil
	}

	log.WithFields(log.Fields{
		"username": username,
		"cname":    cname,
		"up[2]":    up[2],
	}).Debug("Calling RemoveContainerOwner() with")

	err = s.RemoveContainerOwner(username, cname, up[2])
	if err != nil {
		return err
	}

	return nil
}
