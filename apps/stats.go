package apps

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo/v4"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func Stats(ctx context.Context, app string, stream bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errgo.Notef(err, "fail to get Scalingo client")
	}

	if stream {
		stats, err := c.AppsStats(ctx, app)
		if err != nil {
			return errgo.Mask(err)
		}
		displayLiveStatsTable(stats.Stats)

		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				stats, err := c.AppsStats(ctx, app)
				if err != nil {
					ticker.Stop()
					return errgo.Mask(err)
				}
				displayLiveStatsTable(stats.Stats)
			}
		}
	} else {
		stats, err := c.AppsStats(ctx, app)
		if err != nil {
			return errgo.Mask(err)
		}
		return displayStatsTable(stats.Stats)
	}
}

func displayLiveStatsTable(stats []*scalingo.ContainerStat) {
	fmt.Print("\033[2J\033[;H")
	fmt.Printf("Refreshing every 10 seconds...\n\n")
	displayStatsTable(stats)
	fmt.Println("Last update at:", time.Now().Format(time.UnixDate))
}

func displayStatsTable(stats []*scalingo.ContainerStat) error {
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader([]string{"Name", "CPU", "Memory", "Swap"})

	for i, s := range stats {
		t.Append([]string{
			s.ID,
			fmt.Sprintf("%d%%", s.CpuUsage),
			fmt.Sprintf(
				"%2d%% %v/%v",
				int(float64(s.MemoryUsage)/float64(s.MemoryLimit)*100),
				humanize.IBytes(uint64(s.MemoryUsage)),
				humanize.IBytes(uint64(s.MemoryLimit)),
			),
			fmt.Sprintf(
				"%2d%% %v/%v",
				int(float64(s.SwapUsage)/float64(s.SwapLimit)*100),
				humanize.IBytes(uint64(s.SwapUsage)),
				humanize.IBytes(uint64(s.SwapLimit)),
			),
		})
		t.Append([]string{
			"", "",
			fmt.Sprintf("Highest: %v", humanize.IBytes(uint64(s.HighestMemoryUsage))),
			fmt.Sprintf("Highest: %v", humanize.IBytes(uint64(s.HighestSwapUsage))),
		})
		if i != len(stats)-1 {
			t.Append([]string{"", "", "", ""})
		}
	}

	t.Render()
	return nil
}
