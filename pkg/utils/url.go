package utils

import (
	"fmt"
	"net/url"
	"strings"
)

type ParsedURL struct {
	Protocol string
	Slashes  bool
	Auth     *string
	Host     string
	Port     *string
	Hostname string
	Hash     *string
	Search   *string
	Query    *string
	Pathname string
	Path     string
	Href     string
	Origin   *string
}

func SafelyParseURL(raw string) (*ParsedURL, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("empty URL")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	// NOTE: require at least a host or path to be meaningful
	if u.Host == "" && u.Path == "" {
		return nil, fmt.Errorf("invalid URL: no host or path found")
	}

	parsed := &ParsedURL{}

	parsed.Protocol = u.Scheme
	if parsed.Protocol != "" {
		parsed.Protocol += ":"
	}

	parsed.Slashes = u.Host != "" || strings.HasPrefix(u.Scheme, "http") ||
		strings.HasPrefix(u.Scheme, "ftp")

	// Auth (userinfo)
	if u.User != nil {
		auth := u.User.String()
		parsed.Auth = &auth
	}

	parsed.Host = u.Host

	parsed.Hostname = u.Hostname()
	port := u.Port()
	if port != "" {
		parsed.Port = &port
	}

	// Hash (fragment)
	if u.Fragment != "" {
		hash := "#" + u.Fragment
		parsed.Hash = &hash
	}

	if u.RawQuery != "" {
		search := "?" + u.RawQuery
		parsed.Search = &search
		parsed.Query = &u.RawQuery
	}

	parsed.Pathname = u.Path
	if parsed.Pathname == "" {
		parsed.Pathname = "/"
	}

	// Path = pathname + search
	parsed.Path = parsed.Pathname
	if parsed.Search != nil {
		parsed.Path += *parsed.Search
	}

	parsed.Href = u.String()

	// Origin — scheme + "://" + host (only for absolute URLs with a host)
	if u.Scheme != "" && u.Host != "" {
		origin := u.Scheme + "://" + u.Host
		parsed.Origin = &origin
	}

	return parsed, nil
}
