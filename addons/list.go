package addons

import (
	"encoding/json"
	"os"

	"github.com/Scalingo/cli/api"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

type Addon struct {
	LogoURL   string `json:"logo_url"`
	Name      string `json:"name"`
	NameParam string `json:"name_param"`
}

type ListParams struct {
	Addons []*Addon `json:"addons"`
}

func List() error {
	res, err := api.AddonsList()
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}
	defer res.Body.Close()

	var params ListParams
	err = json.NewDecoder(res.Body).Decode(&params)
	if err != nil {
		return errgo.Mask(err, errgo.Any)
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"ID", "Name"})

	for _, addon := range params.Addons {
		t.Append([]string{addon.NameParam, addon.Name})
	}

	t.Render()
	return nil
}
