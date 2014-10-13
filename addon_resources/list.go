package addon_resources

import (
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
)

type AddonResource struct {
	ResourceID string `json:"resource_id"`
	Plan       string `json:"plan"`
	PlanID     string `json:"plan_id"`
	Addon      string `json:"addon"`
	AddonID    string `json:"addon_id"`
}

type ListAddonResourcesParams struct {
	AddonResource []*AddonResource `json:"addon_resources"`
}

func List(app string) error {
	res, err := api.AddonResourcesList(app)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var params ListAddonResourcesParams
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return err
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Plan", "Addon"})

	for _, resource := range params.AddonResource {
		t.Append([]string{resource.ResourceID, resource.Plan, resource.Addon})
	}
	t.Render()

	return nil

}
