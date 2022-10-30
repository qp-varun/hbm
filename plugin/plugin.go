package plugin

import (
	"encoding/json"
	"net/url"
	"regexp"

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
	log.Debug("AuthZReq:", req)
	uriinfo, err := uri.GetURIInfo(req)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	if req.RequestMethod == "OPTIONS" || stringInRegexpSlice(uriinfo.Path, p.skipEndpoints) {
		return authorization.Response{Allow: true}
	}

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
	log.Debug("AuthZRes:", req)
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		log.Debug("parse error:", err)
		return authorization.Response{Allow: true, Msg: err.Error()}
	}

	cname := u.Query().Get("name")
	if p.isCreateContainer(req, u) {
		log.Debug("setting owner for", cname)
		err = p.setContainerOwner(cname, req)
		log.Debug("setContainerOwner err:", err)
	}

	if p.isRemoveContainer(req, u) {
		log.Debug("calling p.removeContainerOwner")
		err = p.removeContainerOwner(cname, req)
		log.Debug("removeContainerOwner err:", err)
	}

	return authorization.Response{Allow: true}
}

func (p *plugin) isCreateContainer(req authorization.Request, u *url.URL) bool {
	log.Debug("entering isCreateContainer")
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
	log.Debug("entering isRemoveContainer")
	if req.ResponseStatusCode != 204 {
		return false
	}

	avm := regexp.MustCompile(`^/v\d+\.\d+/containers/[^/]+`)
	if avm.MatchString(u.Path) && req.RequestMethod == "DELETE" {
		log.Debug("it is removecontainer:", u.Path)
		log.Debug("isRemoveContainer req:", req)
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

	s.SetContainerOwner(username, cname, rjson.Id)

	log.Debug("did owner with:", username, cname, rjson.Id)

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

	var rjson struct {
		Id string
	}
	err = json.Unmarshal(req.ResponseBody, &rjson)
	if err != nil {
		return err
	}

	s.RemoveContainerOwner(username, cname, rjson.Id)

	log.Debug("removed owner with:", username, cname, rjson.Id)

	return nil
}
