package models

import (
	//nolint:gosec // We are not using this for security purposes
	"crypto/md5"
	"encoding/hex"
	"io"
	"regexp"
	"strings"

	"github.com/weslleycapelari/boss/pkg/env"

	"github.com/weslleycapelari/boss/pkg/msg"
)

type Dependency struct {
	Repository string
	version    string
	UseSSH     bool
}

func (p *Dependency) HashName() string {
	//nolint:gosec // We are not using this for security purposes
	hash := md5.New()
	if _, err := io.WriteString(hash, p.Repository); err != nil {
		msg.Warn("Failed on write dependency hash")
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (p *Dependency) GetVersion() string {
	return p.version
}

func (p *Dependency) sshURL() string {
	if strings.Contains(p.Repository, "@") {
		return p.Repository
	}
	re = regexp.MustCompile(`(?m)([\w\d.]*)(?:/)(.*)`)
	submatch := re.FindStringSubmatch(p.Repository)
	provider := submatch[1]
	repo := submatch[2]
	return "git@" + provider + ":" + repo
}

func (p *Dependency) GetURLPrefix() string {
	urlPrefixPattern := regexp.MustCompile(`^[^/^:]+`)
	return urlPrefixPattern.FindString(p.Repository)
}

func (p *Dependency) GetURL() string {
	prefix := p.GetURLPrefix()
	auth := env.GlobalConfiguration().Auth[prefix]
	if auth != nil {
		if auth.UseSSH {
			return p.sshURL()
		}
	}
	var hasHTTPS = regexp.MustCompile(`(?m)^https?:\/\/`)
	if hasHTTPS.MatchString(p.Repository) {
		return p.Repository
	}

	return "https://" + p.Repository
}

var re = regexp.MustCompile(`(?m)^(.|)(\d+)\.(\d+)$`)
var re2 = regexp.MustCompile(`(?m)^(.|)(\d+)$`)

func ParseDependency(repo string, info string) Dependency {
	parsed := strings.Split(info, ":")
	dependency := Dependency{}
	dependency.Repository = repo
	dependency.version = parsed[0]
	if re.MatchString(dependency.version) {
		msg.Warn("Current version for %s is not semantic (x.y.z), for comparison using %s -> %s",
			dependency.Repository, dependency.version, dependency.version+".0")
		dependency.version += ".0"
	}
	if re2.MatchString(dependency.version) {
		msg.Warn("Current version for %s is not semantic (x.y.z), for comparison using %s -> %s",
			dependency.Repository, dependency.version, dependency.version+".0.0")
		dependency.version += ".0.0"
	}
	if len(parsed) > 1 {
		dependency.UseSSH = parsed[1] == "ssh"
	}
	return dependency
}

func GetDependencies(deps map[string]string) []Dependency {
	dependencies := make([]Dependency, 0)
	for repo, info := range deps {
		dependencies = append(dependencies, ParseDependency(repo, info))
	}
	return dependencies
}

func GetDependenciesNames(deps []Dependency) []string {
	var dependencies []string
	for _, info := range deps {
		dependencies = append(dependencies, info.Name())
	}
	return dependencies
}

func (p *Dependency) Name() string {
	var re = regexp.MustCompile(`[^/]+(:?/$|$)`)
	return re.FindString(p.Repository)
}
