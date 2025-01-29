/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package akamai

import (
	"fmt"
	"os"

	internalssh "github.com/konstructio/kubefirst-api/pkg/ssh"
	"github.com/konstructio/kubefirst/internal/catalog"
	"github.com/konstructio/kubefirst/internal/progress"
	"github.com/konstructio/kubefirst/internal/provision"
	"github.com/konstructio/kubefirst/internal/utilities"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func createAkamai(cmd *cobra.Command, _ []string) error {
	cliFlags, err := utilities.GetFlags(cmd, "akamai")
	if err != nil {
		progress.Error(err.Error())
		return fmt.Errorf("failed to get flags: %w", err)
	}

	progress.DisplayLogHints(25)

	isValid, catalogApps, err := catalog.ValidateCatalogApps(cliFlags.InstallCatalogApps)
	if !isValid {
		return fmt.Errorf("catalog validation failed: %w", err)
	}

	err = ValidateProvidedFlags(cliFlags.GitProvider, cliFlags.DNSProvider)
	if err != nil {
		progress.Error(err.Error())
		return fmt.Errorf("failed to validate flags: %w", err)
	}

	if err := provision.ManagementCluster(cliFlags, catalogApps); err != nil {
		return fmt.Errorf("failed to provision management cluster: %w", err)
	}

	return nil
}

func ValidateProvidedFlags(gitProvider, dnsProvider string) error {
	progress.AddStep("Validate provided flags")

	if os.Getenv("LINODE_TOKEN") == "" {
		return fmt.Errorf("your LINODE_TOKEN is not set - please set and re-run your last command")
	}

	if dnsProvider == "cloudflare" {
		if os.Getenv("CF_API_TOKEN") == "" {
			return fmt.Errorf("your CF_API_TOKEN environment variable is not set. Please set and try again")
		}
	}

	switch gitProvider {
	case "github":
		key, err := internalssh.GetHostKey("github.com")
		if err != nil {
			return fmt.Errorf("failed to fetch github host key: %w", err)
		}
		log.Info().Msgf("%q %s", "github.com", key.Type())
	case "gitlab":
		key, err := internalssh.GetHostKey("gitlab.com")
		if err != nil {
			return fmt.Errorf("failed to fetch gitlab host key: %w", err)
		}
		log.Info().Msgf("%q %s", "gitlab.com", key.Type())
	}

	progress.CompleteStep("Validate provided flags")

	return nil
}
