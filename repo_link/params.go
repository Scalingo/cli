package repo_link

import (
	"gopkg.in/errgo.v1"

	"github.com/urfave/cli"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/go-scalingo"
)

func CheckAndFillParams(c *cli.Context, app string) (*scalingo.ScmRepoLinkParams, error) {
	sc, err := config.ScalingoClient()
	if err != nil {
		return nil, errgo.Notef(err, "fail to get Scalingo client")
	}

	// Get RepoLink of App
	repoLink, err := sc.ScmRepoLinkShow(app)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	// Get params
	paramsChecker := NewParamsChecker(repoLink, c)
	params := &scalingo.ScmRepoLinkParams{
		Branch:                   paramsChecker.lookupBranch(),
		AutoDeployEnabled:        paramsChecker.lookupAutoDeploy(),
		DeployReviewAppsEnabled:  paramsChecker.lookupDeployReviewApps(),
		DestroyOnCloseEnabled:    paramsChecker.lookupDeleteOnClose(),
		HoursBeforeDeleteOnClose: paramsChecker.lookupHoursBeforeDeleteOnClose(),
		DestroyStaleEnabled:      paramsChecker.lookupDeleteStale(),
		HoursBeforeDeleteStale:   paramsChecker.lookupHoursBeforeDeleteStale(),
	}

	return params, nil
}

type ParamsChecker struct {
	repoLink *scalingo.ScmRepoLink
	ctx      *cli.Context
}

func NewParamsChecker(repoLink *scalingo.ScmRepoLink, ctx *cli.Context) *ParamsChecker {
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

func (p *ParamsChecker) lookupDeleteOnClose() *bool {
	deleteOnClose := p.ctx.Bool("delete-on-close")

	if p.repoLink.DeployReviewAppsEnabled != deleteOnClose {
		return &deleteOnClose
	}
	return nil
}

func (p *ParamsChecker) lookupHoursBeforeDeleteOnClose() *uint {
	hoursBeforeDeleteOnClose := p.ctx.Uint("hours-before-delete-on-close")

	if p.repoLink.HoursBeforeDeleteOnClose != hoursBeforeDeleteOnClose {
		return &hoursBeforeDeleteOnClose
	}
	return nil
}

func (p *ParamsChecker) lookupDeleteStale() *bool {
	deleteStale := p.ctx.Bool("delete-stale")

	if p.repoLink.DeleteOnStaleEnabled != deleteStale {
		return &deleteStale
	}
	return nil
}

func (p *ParamsChecker) lookupHoursBeforeDeleteStale() *uint {
	hoursBeforeDeleteStale := p.ctx.Uint("hours-before-delete-stale")

	if p.repoLink.HoursBeforeDeleteStale != hoursBeforeDeleteStale {
		return &hoursBeforeDeleteStale
	}
	return nil
}
