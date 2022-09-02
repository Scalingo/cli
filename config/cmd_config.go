package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"
)

func SetRegion(ctx context.Context, regionName string) error {
	region, err := GetRegion(ctx, C, regionName, GetRegionOpts{})
	if err != nil {
		return errgo.Notef(err, "fail to select region")
	}

	C.ConfigFile.Region = region.Name
	fd, err := os.OpenFile(C.ConfigFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return errgo.Notef(err, "fail to open config file")
	}
	defer fd.Close()

	err = json.NewEncoder(fd).Encode(C.ConfigFile)
	if err != nil {
		return errgo.Notef(err, "fail to persist config file %v", C.ConfigFilePath)
	}

	return nil
}

func Display() {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetColWidth(60)
	t.SetHeader([]string{"Configuration key", "Value"})
	t.Append([]string{"region", C.ConfigFile.Region})
	t.Render()
}
