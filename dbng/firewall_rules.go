package dbng

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func FirewallRulesList(ctx context.Context, databaseIdentifier string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, addon, err := databaseAndAddonByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return err
	}

	rules, err := c.Preview().FirewallRulesList(ctx, database.Name, addon.ID)
	if err != nil {
		return errors.Wrap(ctx, err, "list firewall rules")
	}

	if len(rules) == 0 {
		io.Status("No firewall rules configured for this database.")
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "Type", "CIDR or Range ID", "Label"})

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

func FirewallRulesAdd(ctx context.Context, databaseIdentifier string, params scalingo.FirewallRuleCreateParams) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, addon, err := databaseAndAddonByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return err
	}

	rule, err := c.Preview().FirewallRulesCreate(ctx, database.Name, addon.ID, params)
	if err != nil {
		return errors.Wrap(ctx, err, "add firewall rule")
	}

	io.Statusf("Firewall rule '%s' has been added to database '%s'.\n", rule.ID, database.ID)
	io.Warning("Expect some delay for the firewall rules to be applied.")

	return nil
}

func FirewallRulesRemove(ctx context.Context, databaseIdentifier, ruleID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, addon, err := databaseAndAddonByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return err
	}

	err = c.Preview().FirewallRulesDestroy(ctx, database.Name, addon.ID, ruleID)
	if err != nil {
		return errors.Wrap(ctx, err, "remove firewall rule")
	}

	io.Statusf("Firewall rule '%s' has been removed from database '%s'.\n", ruleID, database.ID)
	io.Warning("Expect some delay for the firewall rules to be applied.")

	return nil
}

func FirewallManagedRangesList(ctx context.Context, databaseIdentifier string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get Scalingo client")
	}

	database, addon, err := databaseAndAddonByNameOrID(ctx, c, databaseIdentifier)
	if err != nil {
		return err
	}

	ranges, err := c.Preview().FirewallRulesGetManagedRanges(ctx, database.Name, addon.ID)
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
