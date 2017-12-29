package resource

import (
	"github.com/aelsabbahy/goss/system"
	"github.com/aelsabbahy/goss/util"
)

type HTTP struct {
	Title             string   `json:"title,omitempty" yaml:"title,omitempty"`
	Meta              meta     `json:"meta,omitempty" yaml:"meta,omitempty"`
	HTTP              string   `json:"-" yaml:"-"`
	Status            matcher  `json:"status" yaml:"status"`
	Header            []string `json:"header" yaml:"header"`
	Body              []string `json:"body" yaml:"body"`
	AllowInsecure     bool     `json:"allow-insecure" yaml:"allow-insecure"`
	NoFollowRedirects bool     `json:"no-follow-redirects" yaml:"no-follow-redirects"`
	XForwardedSSL     bool     `json:"x-forwarded-ssl" yaml:"x-forwarded-ssl"`
	Host              string   `json:"host" yaml:"host"`
	Timeout           int      `json:"timeout" yaml:"timeout"`
}

func (u *HTTP) ID() string      { return u.HTTP }
func (u *HTTP) SetID(id string) { u.HTTP = id }

// FIXME: Can this be refactored?
func (r *HTTP) GetTitle() string { return r.Title }
func (r *HTTP) GetMeta() meta    { return r.Meta }

func (u *HTTP) Validate(sys *system.System) []TestResult {
	skip := false
	if u.Timeout == 0 {
		u.Timeout = 5000
	}
	sysHTTP := sys.NewHTTP(u.HTTP, sys, util.Config{AllowInsecure: u.AllowInsecure, NoFollowRedirects: u.NoFollowRedirects, XForwardedSSL: u.XForwardedSSL, Host: u.Host, Timeout: u.Timeout})
	sysHTTP.SetAllowInsecure(u.AllowInsecure)
	sysHTTP.SetNoFollowRedirects(u.NoFollowRedirects)
	sysHTTP.SetXForwardedSSL(u.XForwardedSSL)
	sysHTTP.SetHost(u.Host)

	var results []TestResult
	results = append(results, ValidateValue(u, "status", u.Status, sysHTTP.Status, skip))
	if shouldSkip(results) {
		skip = true
	}
	if len(u.Header) > 0 {
        results = append(results, ValidateContains(u, "Header", u.Header, sysHTTP.Header, skip))
    }
	if len(u.Body) > 0 {
		results = append(results, ValidateContains(u, "Body", u.Body, sysHTTP.Body, skip))
	}

	return results
}

func NewHTTP(sysHTTP system.HTTP, config util.Config) (*HTTP, error) {
	http := sysHTTP.HTTP()
	status, err := sysHTTP.Status()
	u := &HTTP{
		HTTP:              http,
		Status:            status,
		Header:            []string{},
		Body:              []string{},
		AllowInsecure:     config.AllowInsecure,
		NoFollowRedirects: config.NoFollowRedirects,
		XForwardedSSL:     config.XForwardedSSL,
		Host:              config.Host,
		Timeout:           config.Timeout,
	}
	return u, err
}
