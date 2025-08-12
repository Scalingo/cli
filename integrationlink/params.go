package integrationlink

import (
	"github.com/urfave/cli/v3"

	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

func CheckAndFillParams(c *cli.Context) *scalingo.SCMRepoLinkUpdateParams {
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
	ctx *cli.Context
}

func newParamsChecker(ctx *cli.Context) *paramsChecker {
	return &paramsChecker{ctx: ctx}
}

func (p *paramsChecker) lookupBranch() *string {
	if !p.ctx.IsSet("branch") {
		return nil
	}

	branch := p.ctx.String("branch")
	return &branch
}

func (p *paramsChecker) lookupAutoDeploy() *bool {
	if p.ctx.IsSet("auto-deploy") {
		return utils.BoolPtr(true)
	}
	if p.ctx.IsSet("no-auto-deploy") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupDeployReviewApps() *bool {
	if p.ctx.IsSet("deploy-review-apps") {
		return utils.BoolPtr(true)
	}
	if p.ctx.IsSet("no-deploy-review-apps") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupDestroyOnClose() *bool {
	if p.ctx.IsSet("destroy-on-close") {
		return utils.BoolPtr(true)
	}
	if p.ctx.IsSet("no-destroy-on-close") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupHoursBeforeDestroyOnClose() *uint {
	if !p.ctx.IsSet("hours-before-destroy-on-close") {
		return nil
	}

	hoursBeforeDestroyOnClose := p.ctx.Uint("hours-before-destroy-on-close")
	return &hoursBeforeDestroyOnClose
}

func (p *paramsChecker) lookupDestroyOnStale() *bool {
	if p.ctx.IsSet("destroy-on-stale") {
		return utils.BoolPtr(true)
	}
	if p.ctx.IsSet("no-destroy-on-stale") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupAllowReviewAppsFromForks() *bool {
	if p.ctx.IsSet("allow-review-apps-from-forks") {
		return utils.BoolPtr(true)
	}
	if p.ctx.IsSet("no-allow-review-apps-from-forks") {
		return utils.BoolPtr(false)
	}
	return nil
}

func (p *paramsChecker) lookupHoursBeforeDestroyOnStale() *uint {
	if !p.ctx.IsSet("hours-before-destroy-on-stale") {
		return nil
	}

	hoursBeforeDestroyOnStale := p.ctx.Uint("hours-before-destroy-on-stale")
	return &hoursBeforeDestroyOnStale
}
