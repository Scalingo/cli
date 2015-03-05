package domains

import (
	"fmt"
	"os"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
	"github.com/Scalingo/cli/api"
)

func List(app string) error {
	domains, err := api.DomainsList(app)
	if err != nil {
		return errgo.Mask(err)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Domain", "SSL"})

	for _, domain := range domains {
		if !domain.SSL {
			t.Append([]string{domain.Name, "-"})
		} else {
			t.Append([]string{domain.Name, fmt.Sprintf("Valid until %v", domain.Validity)})
		}
	}
	t.Render()
	return nil
}
