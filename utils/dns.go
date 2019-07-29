package utils

import (
	"net"

	errgo "gopkg.in/errgo.v1"
)

func IsCNAME(domain string) (bool, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return false, errgo.Notef(err, "fail to lookup CNAME")
	}
	return cname != domain, nil
}
