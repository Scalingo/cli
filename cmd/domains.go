package cmd

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"gopkg.in/errgo.v1"

	"github.com/Scalingo/cli/cmd/autocomplete"
	"github.com/Scalingo/cli/detect"
	"github.com/Scalingo/cli/domains"
	"github.com/Scalingo/cli/utils"
	"github.com/Scalingo/go-scalingo/v8"
)

var (
	DomainsListCommand = cli.Command{
		Name:     "domains",
		Category: "Custom Domains",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "List the domains of an application",
		Description: CommandDescription{
			Description: "List all the custom domains of an application",
			Examples:    []string{"scalingo --app my-app domains"},
			SeeAlso:     []string{"domains-add", "domains-remove"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() > 0 {
				err := cli.ShowCommandHelp(c, "domains")
				if err != nil {
					errorQuit(ctx, err)
				}
				return nil
			}

			currentApp := detect.CurrentApp(c)
			err := domains.List(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "domains")
		},
	}

	DomainsAddCommand = cli.Command{
		Name:      "domains-add",
		Category:  "Custom Domains",
		Usage:     "Add a custom domain to an application",
		ArgsUsage: "domain",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate"},
			&cli.StringFlag{Name: "key", Usage: "SSL Keypair"},
			&cli.BoolFlag{Name: "no-letsencrypt", Usage: "Disable automatic Let's Encrypt certificate generation", Value: false},
		},
		Description: CommandDescription{
			Description: "Add a custom domain to an application",
			Examples: []string{
				"scalingo --app my-app domains-add example.com",
				"scalingo --app my-app domains-add --cert ./cert.pem --key ./key.pem example.com",
				"scalingo --app my-app domains-add --no-letsencrypt example.com",
			},
			SeeAlso: []string{"domains", "domains-remove"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			if c.Args().Len() == 1 {
				cert := c.String("cert")
				key := c.String("key")
				certContent, keyContent, err := validateTLSParams(cert, key)
				if err != nil {
					errorQuit(ctx, err)
				}

				params := scalingo.DomainsAddParams{
					Name: c.Args().First(),
				}
				if certContent != "" {
					params.TLSCert = &certContent
				}
				if keyContent != "" {
					params.TLSKey = &keyContent
				}
				if c.Bool("no-letsencrypt") {
					params.LetsEncryptEnabled = utils.BoolPtr(false)
				}

				err = domains.Add(ctx, currentApp, params)
				if err != nil {
					errorQuit(ctx, err)
				}
			} else {
				_ = cli.ShowCommandHelp(ctx, c, "domains-add")
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "domains-add")
		},
	}

	DomainsRemoveCommand = cli.Command{
		Name:      "domains-remove",
		Category:  "Custom Domains",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Remove a custom domain from an application",
		ArgsUsage: "domain",
		Description: CommandDescription{
			Description: "Remove a custom domain from an application",
			Examples:    []string{"scalingo --app my-app domains-remove example.com"},
			SeeAlso:     []string{"domains", "domains-add"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
			var err error
			if c.Args().Len() == 1 {
				err = domains.Remove(ctx, currentApp, c.Args().First())
			} else {
				_ = cli.ShowCommandHelp(ctx, c, "domains-remove")
			}

			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(ctx context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "domains-remove")
			_ = autocomplete.DomainsRemoveAutoComplete(ctx, c)
		},
	}

	// TODO: Split the two operations (enable/disable) into two subcommands
	DomainsSSLCommand = cli.Command{
		Name:      "domains-ssl",
		Category:  "Custom Domains",
		Usage:     "Enable or disable SSL for your custom domains",
		ArgsUsage: "domain [disable]",
		Flags: []cli.Flag{&appFlag,
			&cli.StringFlag{Name: "cert", Usage: "SSL Signed Certificate", Value: "domain.crt"},
			&cli.StringFlag{Name: "key", Usage: "SSL Keypair", Value: "domain.key"},
		},
		Description: CommandDescription{
			Description: "Enable or disable SSL for your custom domains",
			Examples: []string{
				"scalingo --app my-app domains-ssl --cert certificate.crt --key keyfile.key example.com",
				"scalingo --app my-app domains-ssl example.com disable",
			},
			SeeAlso: []string{"domains", "domains-add"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)
			if c.Args().Len() == 2 && c.Args().Slice()[1] == "disable" {
				err := domains.DisableSSL(ctx, currentApp, c.Args().First())
				if err != nil {
					errorQuit(ctx, err)
				}
			} else if c.Args().Len() == 1 {
				certContent, keyContent, err := validateTLSParams(c.String("cert"), c.String("key"))
				if err != nil {
					errorQuit(ctx, err)
				}
				err = domains.EnableSSL(ctx, currentApp, c.Args().First(), certContent, keyContent)
				if err != nil {
					errorQuit(ctx, err)
				}
			} else {
				_ = cli.ShowCommandHelp(ctx, c, "domains-ssl")
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "domains-ssl")
		},
	}

	setCanonicalDomainCommand = cli.Command{
		Name:      "set-canonical-domain",
		Category:  "App Management",
		Flags:     []cli.Flag{&appFlag},
		Usage:     "Set a canonical domain.",
		ArgsUsage: "domain",
		Description: CommandDescription{
			Description: `After defining multiple domains on an application, one can need to redirect all requests towards its application to a specific domain.
This domain is called the canonical domain. This command sets the canonical domain`,
			Examples: []string{"scalingo -a myapp set-canonical-domain example.com"},
			SeeAlso:  []string{"domains", "domains-add", "unset-canonical-domain"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 1 {
				_ = cli.ShowCommandHelp(ctx, c, "set-canonical-domain")
				return nil
			}
			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err := domains.SetCanonical(ctx, currentApp, c.Args().First())
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "set-canonical-domain")
		},
	}

	unsetCanonicalDomainCommand = cli.Command{
		Name:     "unset-canonical-domain",
		Category: "App Management",
		Flags:    []cli.Flag{&appFlag},
		Usage:    "Unset a canonical domain.",
		Description: CommandDescription{
			Description: "Unset the canonical domain of this app",
			Examples:    []string{"scalingo --app my-app unset-canonical-domain"},
			SeeAlso:     []string{"domains", "domains-add", "set-canonical-domain"},
		}.Render(),

		Action: func(ctx context.Context, c *cli.Command) error {
			currentApp := detect.CurrentApp(c)
			if c.Args().Len() != 0 {
				_ = cli.ShowCommandHelp(ctx, c, "unset-canonical-domain")
				return nil
			}

			utils.CheckForConsent(ctx, currentApp, utils.ConsentTypeContainers)

			err := domains.UnsetCanonical(ctx, currentApp)
			if err != nil {
				errorQuit(ctx, err)
			}
			return nil
		},
		ShellComplete: func(_ context.Context, c *cli.Command) {
			_ = autocomplete.CmdFlagsAutoComplete(c, "unset-canonical-domain")
		},
	}
)

func validateTLSParams(cert, key string) (string, string, error) {
	if cert == "" && key == "" {
		return "", "", nil
	}

	if cert == "" && key != "" {
		return "", "", errgo.New("--cert <certificate path> should be defined")
	}

	if key == "" && cert != "" {
		return "", "", errgo.New("--key <key path> should be defined")
	}

	certContent, err := os.ReadFile(cert)
	if err != nil {
		return "", "", errgo.Notef(err, "fail to read the TLS certificate")
	}
	keyContent, err := os.ReadFile(key)
	if err != nil {
		return "", "", errgo.Notef(err, "fail to read the private key")
	}
	return string(certContent), string(keyContent), nil
}
