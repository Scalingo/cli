package integrationlink

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v9"
)

func CheckAndFillParams(c *cli.Command) *scalingo.SCMRepoLinkUpdateParams {
	paramsChecker := newParamsChecker(c)
	params := &scalingo.SCMRepoLinkUpdateParams{
		Branch:                            paramsChecker.lookupBranch(),
		AutoDeployEnabled:                 paramsChecker.lookupAutoDeploy(),
		DeployReviewAppsEnabled:           paramsChecker.lookupDeployReviewApps(),
		DestroyOnCloseEnabled:             paramsChecker.lookupDestroyOnClose(),
		HoursBeforeDeleteOnClose:          paramsChecker.lookupHoursBeforeDestroyOnClose(),
		DestroyStaleEnabled:               paramsChecker.lookupDestroyOnStale(),
		HoursBeforeDeleteStale:            paramsChecker.lookupHoursBeforeDestroyOnStale(),
		AutomaticCreationFromForksAllowed: paramsChecker.lookupAllowReviewAppsFromForks(),
	}

	return params
}

type paramsChecker struct {
	c *cli.Command
}

func newParamsChecker(c *cli.Command) *paramsChecker {
	return &paramsChecker{c: c}
}

func (p *paramsChecker) lookupBranch() *string {
	if !p.c.IsSet("branch") {
		return nil
	}

	branch := p.c.String("branch")
	return &branch
}

func (p *paramsChecker) lookupAutoDeploy() *bool {
	if p.c.IsSet("auto-deploy") {
		return utils.BoolPtr(true)
	}
	if p.c.IsSet("no-auto-deploy") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupDeployReviewApps() *bool {
	if p.c.IsSet("deploy-review-apps") {
		return utils.BoolPtr(true)
	}
	if p.c.IsSet("no-deploy-review-apps") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupDestroyOnClose() *bool {
	if p.c.IsSet("destroy-on-close") {
		return utils.BoolPtr(true)
	}
	if p.c.IsSet("no-destroy-on-close") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupHoursBeforeDestroyOnClose() *uint {
	if !p.c.IsSet("hours-before-destroy-on-close") {
		return nil
	}

	hoursBeforeDestroyOnClose := p.c.Uint("hours-before-destroy-on-close")
	return &hoursBeforeDestroyOnClose
}

func (p *paramsChecker) lookupDestroyOnStale() *bool {
	if p.c.IsSet("destroy-on-stale") {
		return utils.BoolPtr(true)
	}
	if p.c.IsSet("no-destroy-on-stale") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupAllowReviewAppsFromForks() *bool {
	if p.c.IsSet("allow-review-apps-from-forks") {
		return utils.BoolPtr(true)
	}
	if p.c.IsSet("no-allow-review-apps-from-forks") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupHoursBeforeDestroyOnStale() *uint {
	if !p.c.IsSet("hours-before-destroy-on-stale") {
		return nil
	}

	hoursBeforeDestroyOnStale := p.c.Uint("hours-before-destroy-on-stale")
	return &hoursBeforeDestroyOnStale
}
