package dbng

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v8"
	"github.com/Scalingo/go-utils/errors/v2"
)

func FirewallRulesList(ctx context.Context, databaseID, addonID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	rules, err := c.Preview().FirewallRulesList(ctx, databaseID, addonID)
	if err != nil {
		return errors.Wrap(ctx, err, "list firewall rules")
	}

	if len(rules) == 0 {
		io.Status("No firewall rules configured for this database.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Type", "CIDR/Range", "Label"})

	for _, rule := range rules {
		cidrOrRange := rule.CIDR
		if rule.Type == scalingo.FirewallRuleTypeManagedRange {
			cidrOrRange = rule.RangeID
		}

		_ = t.Append([]string{
			rule.ID,
			string(rule.Type),
			cidrOrRange,
			rule.Label,
		})
	}
	_ = t.Render()

	return nil
}

func FirewallRulesAdd(ctx context.Context, databaseID, addonID string, params scalingo.FirewallRuleCreateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	rule, err := c.Preview().FirewallRulesCreate(ctx, databaseID, addonID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add firewall rule")
	}

	io.Statusf("Firewall rule '%s' has been added to database '%s'.\n", rule.ID, databaseID)
	io.Warning("Firewall rules take time to be applied due to infrastructure provisioning.")

	return nil
}

func FirewallRulesRemove(ctx context.Context, databaseID, addonID, ruleID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	err = c.Preview().FirewallRulesDestroy(ctx, databaseID, addonID, ruleID)
	if err != nil {
		return errors.Wrap(ctx, err, "remove firewall rule")
	}

	io.Statusf("Firewall rule '%s' has been removed from database '%s'.\n", ruleID, databaseID)
	io.Warning("Firewall rules take time to be applied due to infrastructure provisioning.")

	return nil
}

func FirewallManagedRangesList(ctx context.Context, databaseID, addonID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	ranges, err := c.Preview().FirewallRulesGetManagedRanges(ctx, databaseID, addonID)
	if err != nil {
		return errors.Wrap(ctx, err, "list managed ranges")
	}

	if len(ranges) == 0 {
		io.Status("No managed ranges available.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Name"})

	for _, r := range ranges {
		_ = t.Append([]string{
			r.ID,
			r.Name,
		})
	}
	_ = t.Render()

	return nil
}
