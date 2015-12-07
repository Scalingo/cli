package apps

import (
	"fmt"
	"os"
	"time"

	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/Scalingo/go-scalingo"
	"github.com/Scalingo/cli/Godeps/_workspace/src/github.com/olekukonko/tablewriter"
	"github.com/Scalingo/cli/Godeps/_workspace/src/gopkg.in/errgo.v1"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func Stats(app string, stream bool) error {
	if stream {
		stats, err := scalingo.AppsStats(app)
		if err != nil {
			return errgo.Mask(err)
		}
		displayLiveStatsTable(stats.Stats)

		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				stats, err := scalingo.AppsStats(app)
				if err != nil {
					ticker.Stop()
					return errgo.Mask(err)
				}
				displayLiveStatsTable(stats.Stats)
			}
		}
	} else {
		stats, err := scalingo.AppsStats(app)
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
				toHuman(s.MemoryUsage),
				toHuman(s.MemoryLimit),
			),
			fmt.Sprintf(
				"%2d%% %v/%v",
				int(float64(s.SwapUsage)/float64(s.SwapLimit)*100),
				toHuman(s.SwapUsage),
				toHuman(s.SwapLimit),
			),
		})
		t.Append([]string{
			"", "",
			fmt.Sprintf("Highest: %v", toHuman(s.HighestMemoryUsage)),
			fmt.Sprintf("Highest: %v", toHuman(s.HighestSwapUsage)),
		})
		if i != len(stats)-1 {
			t.Append([]string{"", "", "", ""})
		}
	}

	t.Render()
	return nil
}

func toHuman(sizeInBytes int64) string {
	if sizeInBytes > GB {
		return fmt.Sprintf("%3dGB", sizeInBytes/GB)
	} else if sizeInBytes > MB {
		return fmt.Sprintf("%3dMB", sizeInBytes/MB)
	} else if sizeInBytes > KB {
		return fmt.Sprintf("%3dKB", sizeInBytes/KB)
	} else {
		return fmt.Sprintf("%3dB", sizeInBytes)
	}
}
