package apps

import (
	"context"
	"net/netip"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/Scalingo/cli/config"
	"github.com/Scalingo/cli/io"
	"github.com/Scalingo/go-scalingo/v11"
	"github.com/Scalingo/go-utils/errors/v3"
)

func ForceHTTPS(ctx context.Context, appName string, enable bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.AppsForceHTTPS(ctx, appName, enable)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to configure force-https feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Force HTTPS has been %sd on %s\n", action, appName)
	return nil
}

func StickySession(ctx context.Context, appName string, enable bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}
	_, err = c.AppsStickySession(ctx, appName, enable)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to configure sticky-session feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Sticky session has been %sd on %s\n", action, appName)
	return nil
}

func RouterLogs(ctx context.Context, appName string, enable bool) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	_, err = c.AppsRouterLogs(ctx, appName, enable)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to configure router-logs feature")
	}

	var action string
	if enable {
		action = "enable"
	} else {
		action = "disable"
	}

	io.Statusf("Router logs have been %sd on %s\n", action, appName)
	return nil
}

func AppFirewallRulesList(ctx context.Context, appName string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	rules, err := c.AppsFirewallRulesList(ctx, appName)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to list app firewall rules")
	}

	if len(rules) == 0 {
		io.Statusf("No app firewall rules configured on %s\n", appName)
		return nil
	}

	t := tablewriter.NewWriter(os.Stdout)
	t.Header([]string{"ID", "CIDR", "Label"})
	for _, rule := range rules {
		_ = t.Append([]string{rule.ID, rule.CIDR, rule.Label})
	}
	_ = t.Render()
	return nil
}

func AppFirewallRuleAdd(ctx context.Context, appName string, cidr string, label string) error {
	err := ValidateIPv4CIDR(cidr)
	if err != nil {
		return err
	}

	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	rule, err := c.AppsFirewallRuleCreate(ctx, appName, scalingo.AppFirewallRuleParams{
		CIDR:  cidr,
		Label: label,
	})
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to add app firewall rule")
	}

	io.Statusf("app firewall rule %s has been added to %s\n", rule.ID, appName)
	return nil
}

func AppFirewallRuleRemove(ctx context.Context, appName string, ruleID string) error {
	c, err := config.ScalingoClient(ctx)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to get Scalingo client")
	}

	err = c.AppsFirewallRuleDelete(ctx, appName, ruleID)
	if err != nil {
		return errors.Wrapf(ctx, err, "fail to remove app firewall rule")
	}

	io.Statusf("app firewall rule %s has been removed from %s\n", ruleID, appName)
	return nil
}

func ValidateIPv4CIDR(cidr string) error {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return errors.Newf(context.Background(), "invalid IPv4 CIDR %q", cidr)
	}
	if !prefix.Addr().Is4() {
		return errors.Newf(context.Background(), "invalid IPv4 CIDR %q: IPv6 prefixes are not supported", cidr)
	}
	return nil
}
