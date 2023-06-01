package regionmigrations

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/utils"
	scalingo "github.com/Scalingo/go-scalingo/v6"
	"github.com/Scalingo/go-utils/retry"
)

type RefreshOpts struct {
	ShowHints        bool
	ExpectedStatuses []scalingo.RegionMigrationStatus
	HiddenSteps      []string
	CurrentStep      scalingo.RegionMigrationStep
}

type Refresher struct {
	appID       string
	migrationID string
	client      *scalingo.Client
	opts        RefreshOpts

	lock                 *sync.Mutex
	migration            *scalingo.RegionMigration
	errCount             int
	maxErrors            int
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
		errCount:             0,
		maxErrors:            0,
		opts:                 opts,
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
		errCount := r.errCount
		maxErrors := r.maxErrors
		r.lock.Unlock()
		r.currentLoadersStep = (r.currentLoadersStep + 1) % len(spinner.CharSets[11])

		r.writeMigration(writer, migration, errCount, maxErrors)
		if stop {
			return
		}

		time.Sleep(r.screenRefreshTime)
	}
}

func (r *Refresher) onRefreshError(ctx context.Context, err error, currentAttempt, maxAttempt int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.errCount = currentAttempt + 1
	r.maxErrors = maxAttempt
}

func (r *Refresher) migrationRefresher() error {
	defer r.wg.Done()
	r.lock.Lock()
	client := r.client
	r.lock.Unlock()

	errCount := 0

	retrier := retry.New(
		retry.WithMaxAttempts(10),
		retry.WithWaitDuration(10*time.Second),
		retry.WithErrorCallback(r.onRefreshError),
	)

	for {
		err := retrier.Do(context.Background(), func(ctx context.Context) error {
			migration, err := client.ShowRegionMigration(ctx, r.appID, r.migrationID)
			if err != nil {
				return err
			}
			r.lock.Lock()
			r.errCount = 0
			r.migration = &migration
			r.errCount = errCount
			r.lock.Unlock()
			return nil
		})

		if err != nil {
			r.Stop()
			return errgo.Notef(err, "fail to get migration")
		}

		r.lock.Lock()
		migration := r.migration
		stop := r.stop
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

func (r *Refresher) writeMigration(w *uilive.Writer, migration *scalingo.RegionMigration, errCount, maxErrors int) {
	defer w.Flush()

	if errCount != 0 {
		fmt.Fprint(w.Newline(), color.RedString("Connection lost. Retrying (%v/%v)\n", errCount, maxErrors))
	}

	if migration == nil {
		fmt.Fprint(w.Newline(), color.BlueString("%s Loading migration information\n", r.loader()))
		return
	}

	fmt.Fprintf(w.Newline(), "Migration ID: %s\n", migration.ID)
	fmt.Fprintf(w.Newline(), "Migrating app: %s\n", migration.SrcAppName)
	fmt.Fprintf(w.Newline(), "Destination: %s\n", migration.Destination)
	if migration.NewAppID == "" {
		fmt.Fprintf(w.Newline(), "New app ID: %s\n", color.BlueString("N/A"))
	} else {
		fmt.Fprintf(w.Newline(), "New app ID: %s\n", migration.NewAppID)
	}
	fmt.Fprintf(w.Newline(), "Status: %s\n", formatMigrationStatus(migration.Status))
	if r.opts.ShowHints {
		fmt.Fprintf(w.Newline(), "%s\n", r.hintFor(migration))
	}
	if migration.Status == scalingo.RegionMigrationStatusCreated {
		fmt.Fprintf(w.Newline(), "%s Waiting for the migration to start\n", r.loader())
	}

	for _, step := range migration.Steps {
		if r.shouldShowStep(step) {
			r.writeStep(w, step)
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

func (r *Refresher) shouldStop(m *scalingo.RegionMigration) bool {
	if m == nil {
		return false
	}
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

func (r *Refresher) hintFor(m *scalingo.RegionMigration) string {
	switch m.Status {
	case scalingo.RegionMigrationStatusAborted:
		return "The migration has been aborted. No update will be posted here."
	case scalingo.RegionMigrationStatusCreated:
		return "The migration has been created. The preflight checks will begin shortly."
	case scalingo.RegionMigrationStatusPreflightError:
		return "There was an error during the preflight checks. No update will be posted here."
	case scalingo.RegionMigrationStatusPreflightSuccess:
		return "The preflight checks were successful. Waiting on the user to start the 'prepare' step."
	case scalingo.RegionMigrationStatusRunning:
		return "The migration is currently running."
	case scalingo.RegionMigrationStatusPrepared:
		return "The migration has been prepared. Waiting on the user to start the 'data' or 'finalize' step."
	case scalingo.RegionMigrationStatusDataMigrated:
		return "The addon has been migrated. Waiting on the user to start the 'finalize' step."
	case scalingo.RegionMigrationStatusError:
		return "There was an error while running the migration. Waiting on the user to 'abort' it."
	case scalingo.RegionMigrationStatusDone:
		return "The migration is done. No update will be posted here."
	case scalingo.RegionMigrationStatusAborting:
		return "The migration will be aborted shortly. The abort process will begin shortly."
	}
	return ""
}
