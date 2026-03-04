package utils

import (
	"context"
	"net"

	"github.com/Scalingo/go-utils/errors/v3"
)

func IsCNAME(ctx context.Context, domain string) (bool, error) {
	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return false, errors.Wrapf(ctx, err, "fail to lookup CNAME")
	}
	return cname != domain, nil
}
