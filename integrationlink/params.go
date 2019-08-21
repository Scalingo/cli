package integrationlink

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

type paramsChecker struct {
	repoLink *scalingo.SCMRepoLink
	ctx      *cli.Context
}

func newParamsChecker(repoLink *scalingo.SCMRepoLink, ctx *cli.Context) *paramsChecker {
	return &paramsChecker{repoLink: repoLink, ctx: ctx}
}

func (p *paramsChecker) lookupBranch() *string {
	branch := p.ctx.String("branch")

	if branch != "" && p.repoLink.Branch != branch {
		return &branch
	}
	return nil
}

func (p *paramsChecker) lookupAutoDeploy() *bool {
	autoDeploy := p.ctx.Bool("auto-deploy")
	noAutoDeploy := p.ctx.Bool("no-auto-deploy")

	if autoDeploy {
		return &autoDeploy
	}
	if noAutoDeploy {
		f := false
		return &f
	}
	return nil
}

func (p *paramsChecker) lookupDeployReviewApps() *bool {
	deployReviewApps := p.ctx.Bool("deploy-review-apps")
	noDeployReviewApps := p.ctx.Bool("no-deploy-review-apps")

	if deployReviewApps {
		return &deployReviewApps
	}
	if noDeployReviewApps {
		f := false
		return &f
	}
	return nil
}

func (p *paramsChecker) lookupDestroyOnClose() *bool {
	destroyOnClose := p.ctx.Bool("destroy-on-close")
	noDestroyOnClose := p.ctx.Bool("no-destroy-on-close")

	if destroyOnClose {
		return &destroyOnClose
	}
	if noDestroyOnClose {
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

	if p.repoLink.HoursBeforeDeleteOnClose != hoursBeforeDestroyOnClose {
		return &hoursBeforeDestroyOnClose
	}
	return nil
}

func (p *paramsChecker) lookupDestroyOnStale() *bool {
	destroyOnStale := p.ctx.Bool("destroy-on-stale")
	noDestroyOnStale := p.ctx.Bool("no-destroy-on-stale")

	if destroyOnStale {
		return &destroyOnStale
	}
	if noDestroyOnStale {
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

	if p.repoLink.HoursBeforeDeleteStale != hoursBeforeDestroyOnStale {
		return &hoursBeforeDestroyOnStale
	}
	return nil
}
