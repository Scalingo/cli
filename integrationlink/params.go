package integrationlink

import (
	"github.com/urfave/cli/v2"

	"github.com/Scalingo/go-scalingo/v7"
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
		t := true
		return &t
	}
	if p.ctx.IsSet("no-auto-deploy") {
		f := false
		return &f
	}
	return nil
}

func (p *paramsChecker) lookupDeployReviewApps() *bool {
	if p.ctx.IsSet("deploy-review-apps") {
		t := true
		return &t
	}
	if p.ctx.IsSet("no-deploy-review-apps") {
		f := false
		return &f
	}
	return nil
}

func (p *paramsChecker) lookupDestroyOnClose() *bool {
	if p.ctx.IsSet("destroy-on-close") {
		t := true
		return &t
	}
	if p.ctx.IsSet("no-destroy-on-close") {
		f := false
		return &f
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
		t := true
		return &t
	}
	if p.ctx.IsSet("no-destroy-on-stale") {
		f := false
		return &f
	}
	return nil
}

func (p *paramsChecker) lookupAllowReviewAppsFromForks() *bool {
	if p.ctx.IsSet("allow-review-apps-from-forks") {
		t := true
		return &t
	}
	if p.ctx.IsSet("no-allow-review-apps-from-forks") {
		f := false
		return &f
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
