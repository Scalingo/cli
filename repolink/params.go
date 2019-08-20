package repolink

import (
	"gopkg.in/errgo.v1"

	"github.com/urfave/cli"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
)

func CheckAndFillParams(c *cli.Context, app string) (*scalingo.SCMRepoLinkParams, error) {
	sc, err := config.ScalingoClient()
	if err != nil {
		return nil, errgo.Notef(err, "fail to get Scalingo client")
	}

	repoLink, err := sc.SCMRepoLinkShow(app)
	if err != nil {
		return nil, errgo.Notef(err, "fail to get repo link")
	}

	paramsChecker := newParamsChecker(repoLink, c)
	params := &scalingo.SCMRepoLinkParams{
		Branch:                   paramsChecker.lookupBranch(),
		AutoDeployEnabled:        paramsChecker.lookupAutoDeploy(),
		DeployReviewAppsEnabled:  paramsChecker.lookupDeployReviewApps(),
		DestroyOnCloseEnabled:    paramsChecker.lookupDestroyOnClose(),
		HoursBeforeDeleteOnClose: paramsChecker.lookupHoursBeforeDestroyOnClose(),
		DestroyStaleEnabled:      paramsChecker.lookupDestroyOnStale(),
		HoursBeforeDeleteStale:   paramsChecker.lookupHoursBeforeDestroyOnStale(),
	}

	return params, nil
}

type ParamsChecker struct {
	repoLink *scalingo.SCMRepoLink
	ctx      *cli.Context
}

func newParamsChecker(repoLink *scalingo.SCMRepoLink, ctx *cli.Context) *ParamsChecker {
	return &ParamsChecker{repoLink: repoLink, ctx: ctx}
}

func (p *ParamsChecker) lookupBranch() *string {
	branch := p.ctx.String("branch")

	if branch != "" && p.repoLink.Branch != branch {
		return &branch
	}
	return &p.repoLink.Branch
}

func (p *ParamsChecker) lookupAutoDeploy() *bool {
	autoDeploy := p.ctx.Bool("auto-deploy")

	if p.repoLink.AutoDeployEnabled != autoDeploy {
		return &autoDeploy
	}
	return nil
}

func (p *ParamsChecker) lookupDeployReviewApps() *bool {
	deployReviewApps := p.ctx.Bool("deploy-review-apps")

	if p.repoLink.DeployReviewAppsEnabled != deployReviewApps {
		return &deployReviewApps
	}
	return nil
}

func (p *ParamsChecker) lookupDestroyOnClose() *bool {
	destroyOnClose := p.ctx.Bool("destroy-on-close")

	if p.repoLink.DeployReviewAppsEnabled != destroyOnClose {
		return &destroyOnClose
	}
	return nil
}

func (p *ParamsChecker) lookupHoursBeforeDestroyOnClose() *uint {
	hoursBeforeDestroyOnClose := p.ctx.Uint("hours-before-destroy-on-close")

	if p.repoLink.HoursBeforeDeleteOnClose != hoursBeforeDestroyOnClose {
		return &hoursBeforeDestroyOnClose
	}
	return nil
}

func (p *ParamsChecker) lookupDestroyOnStale() *bool {
	destroyOnStale := p.ctx.Bool("destroy-on-stale")

	if p.repoLink.DeleteStaleEnabled != destroyOnStale {
		return &destroyOnStale
	}
	return nil
}

func (p *ParamsChecker) lookupHoursBeforeDestroyOnStale() *uint {
	hoursBeforeDestroyOnStale := p.ctx.Uint("hours-before-destroy-on-stale")

	if p.repoLink.HoursBeforeDeleteStale != hoursBeforeDestroyOnStale {
		return &hoursBeforeDestroyOnStale
	}
	return nil
}
