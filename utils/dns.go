package utils

import (
	"context"
	"net"

	"github.com/Scalingo/go-utils/errors/v2"
)

func IsCNAME(domain string) (bool, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return false, errors.Wrapf(context.Background(), err, "fail to lookup CNAME")
	}
	return cname != domain, nil
}
