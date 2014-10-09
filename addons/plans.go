package addons

import (
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
)

type Plan struct {
	LogoURL     string `json:"logo_url"`
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PlansParams struct {
	Plans []*Plan `json:"plans"`
}

func Plans(addon string) error {
	res, err := api.AddonPlansList(addon)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var params PlansParams
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return err
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name", "Description"})
	for _, plan := range params.Plans {
		t.Append([]string{plan.Name, plan.DisplayName, plan.Description})
	}
	t.Render()
	return nil
}
