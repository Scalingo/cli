package gitremote

import (
	"fmt"
	"net/url"
	"strings"
)

type Remote struct {
	Name  string
	URL   string
	Fetch string
}

type Remotes []*Remote

func (remotes Remotes) Hosts() []string {
	hosts := make([]string, 0, len(remotes))
	for _, r := range remotes {
		hosts = append(hosts, r.Host())
	}
	return hosts
}

func (remotes Remotes) FilterByName(name string) Remotes {
	var res Remotes
	for _, r := range remotes {
		if r.Name == name {
			res = append(res, r)
		}
	}
	return res
}

func (remotes Remotes) FilterByNamePrefix(prefix string) Remotes {
	var res Remotes
	for _, r := range remotes {
		if strings.HasPrefix(r.Name, prefix) {
			res = append(res, r)
		}
	}
	return res
}

func (remotes Remotes) FilterByHost(host string) Remotes {
	var res Remotes
	for _, r := range remotes {
		if r.Host() == host {
			res = append(res, r)
		}
	}
	return res
}

func (r *Remote) Host() string {
	url, err := url.Parse(r.URL)
	if err == nil && url.Host != "" {
		return url.Host
	}

	matches := sshURLHostRe.FindStringSubmatch(r.URL)
	if len(matches) == 2 {
		return matches[1]
	}

	return "invalid url"
}

func (r *Remote) Repository() string {
	url, err := url.Parse(r.URL)
	if err == nil && url.Host != "" {
		// Remove initial '/'
		return url.Path[1:]
	}

	matches := sshURLRepoRe.FindStringSubmatch(r.URL)
	if len(matches) == 2 {
		return matches[1]
	}

	return "invalid url"
}

func (r *Remote) ToConfig() (string, error) {
	if r.Fetch == "" {
		r.Fetch = defaultFetch(r.Name)
	}

	remoteStr := fmt.Sprintf(`[remote "%s"]
	url = %s
	fetch = %s`,
		r.Name, r.URL, r.Fetch,
	)

	return remoteStr, nil
}

func defaultFetch(name string) string {
	return fmt.Sprintf("+refs/heads/*:refs/remotes/%s/*", name)
}
