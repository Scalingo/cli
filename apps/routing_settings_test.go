package apps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateIPv4CIDR(t *testing.T) {
	runs := map[string]struct {
		cidr      string
		expectErr bool
	}{
		"IPv4 CIDR":      {cidr: "192.0.2.0/24"},
		"IPv4 host CIDR": {cidr: "198.51.100.4/32"},
		"single IP":      {cidr: "192.0.2.1", expectErr: true},
		"malformed CIDR": {cidr: "192.0.2.0/33", expectErr: true},
		"IPv6 CIDR":      {cidr: "2001:db8::/32", expectErr: true},
	}

	for name, run := range runs {
		t.Run(name, func(t *testing.T) {
			err := ValidateIPv4CIDR(run.cidr)
			if run.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
