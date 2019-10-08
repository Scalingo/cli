package region_migrations

import (
	"fmt"
	"sync"
	"time"

	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	errgo "gopkg.in/errgo.v1"
)

type RefreshOpts struct {
	ExpectedStatuses []scalingo.RegionMigrationStatus
	HiddenSteps      []string
}

type Refresher struct {
	appID       string
	migrationID string
	client      *scalingo.Client
	opts        RefreshOpts

	lock                 *sync.Mutex
	migration            *scalingo.RegionMigration
	stop                 bool
	screenRefreshTime    time.Duration
	migrationRefreshTime time.Duration
	wg                   *sync.WaitGroup

	currentLoadersStep int
}

func NewRefresher(client *scalingo.Client, appID, migrationID string, opts RefreshOpts) *Refresher {
	return &Refresher{
		appID:                appID,
		migrationID:          migrationID,
		client:               client,
		lock:                 &sync.Mutex{},
		migration:            nil,
		stop:                 false,
		screenRefreshTime:    100 * time.Millisecond,
		migrationRefreshTime: 1 * time.Second,
		wg:                   &sync.WaitGroup{},
		currentLoadersStep:   0,
	}
}

func (r *Refresher) Start() (*scalingo.RegionMigration, error) {
	r.wg.Add(2)
	go r.screenRefresher()
	err := r.migrationRefresher()
	r.wg.Wait()
	if err != nil {
		return r.migration, errgo.Notef(err, "fail to refresh migration")
	}
	return r.migration, nil
}

func (r *Refresher) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.stop = true
}

func (r *Refresher) screenRefresher() {
	defer r.wg.Done()
	writer := uilive.New()
	for {
		r.lock.Lock()
		stop := r.stop
		migration := r.migration
		r.lock.Unlock()
		r.currentLoadersStep = (r.currentLoadersStep + 1) % len(spinner.CharSets[11])

		r.writeMigration(writer, migration)
		if stop {
			return
		}

		time.Sleep(r.screenRefreshTime)
	}
}

func (r *Refresher) migrationRefresher() error {
	defer r.wg.Done()
	r.lock.Lock()
	client := r.client
	r.lock.Unlock()

	for {
		migration, err := client.ShowRegionMigration(r.appID, r.migrationID)
		if err != nil {
			r.Stop()
			return errgo.Notef(err, "fail to get migration")
		}

		r.lock.Lock()
		stop := r.stop
		r.migration = &migration
		r.lock.Unlock()
		if stop {
			return nil
		}

		if r.shouldStop(migration) {
			r.Stop()
			return nil
		}

		time.Sleep(r.migrationRefreshTime)
	}
}

func (r *Refresher) writeMigration(w *uilive.Writer, migration *scalingo.RegionMigration) {
	defer w.Flush()

	if migration == nil {
		fmt.Fprint(w, color.BlueString("%s Loading migration information\n", r.loader()))
		return
	}

	fmt.Fprintf(w, "Migrating app: %s\n", migration.AppName)
	fmt.Fprintf(w.Newline(), "Destination: %s\n", migration.Destination)
	if migration.NewAppID == "" {
		fmt.Fprintf(w.Newline(), "New app ID: %s\n", color.BlueString("N/A"))
	} else {
		fmt.Fprintf(w.Newline(), "New app ID: %s\n", migration.NewAppID)
	}
	fmt.Fprintf(w.Newline(), "Status: %s\n", formatMigrationStatus(migration.Status))
	if migration.Status == scalingo.RegionMigrationStatusScheduled {
		fmt.Fprintf(w.Newline(), "%s Waiting for the migration to start\n", r.loader())
	}

	for i, _ := range migration.Steps {
		step := migration.Steps[len(migration.Steps)-1-i]
		if r.shouldShowStep(step) {
			r.writeStep(w, migration.Steps[len(migration.Steps)-1-i])
		}
	}

}

func (r *Refresher) writeStep(w *uilive.Writer, step scalingo.Step) {
	result := ""
	switch step.Status {
	case scalingo.StepStatusRunning:
		result = color.BlueString(fmt.Sprintf("%s %s...", r.loader(), step.Name))
	case scalingo.StepStatusDone:
		result = color.GreenString(fmt.Sprintf("%s %s Done!", utils.Success, step.Name))
	case scalingo.StepStatusError:
		result = color.RedString(fmt.Sprintf("%s %s FAILED!", utils.Error, step.Name))
	}
	fmt.Fprintf(w.Newline(), "%s\n", result)
}

func (r *Refresher) loader() string {
	return spinner.CharSets[11][r.currentLoadersStep]
}

func (r *Refresher) shouldStop(m scalingo.RegionMigration) bool {
	switch m.Status {
	case scalingo.RegionMigrationStatusError:
		return true
	case scalingo.RegionMigrationStatusDone:
		return true
	}

	if r.opts.ExpectedStatuses == nil {
		return false
	}

	for _, status := range r.opts.ExpectedStatuses {
		if m.Status == status {
			return true
		}
	}

	return false
}

func (r *Refresher) shouldShowStep(step scalingo.Step) bool {
	if r.opts.HiddenSteps == nil {
		return true
	}

	for _, id := range r.opts.HiddenSteps {
		if id == step.ID {
			return false
		}
	}
	return true
}
