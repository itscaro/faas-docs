// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"fmt"
	"log"

	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-cli/stack"
	"github.com/spf13/cobra"
)

func init() {
	// Setup flags that are used by multiple commands (variables defined in faas.go)
	removeCmd.Flags().StringVar(&gateway, "gateway", defaultGateway, "Gateway URI")

	faasCmd.AddCommand(removeCmd)
}

// removeCmd deletes/removes OpenFaaS function containers
var removeCmd = &cobra.Command{
	Use: `remove FUNCTION_NAME [--gateway GATEWAY_URL]
  faas-cli remove -f YAML_FILE [--regex "REGEX"] [--filter "WILDCARD"]`,
	Aliases: []string{"rm"},
	Short:   "Remove deployed OpenFaaS functions",
	Long: `Removes/deletes deployed OpenFaaS functions either via the supplied YAML config
using the "--yaml" flag (which may contain multiple function definitions), or by
explicitly specifying a function name.`,
	Example: `  faas-cli remove -f https://domain/path/myfunctions.yml
  faas-cli remove -f ./samples.yml
  faas-cli remove -f ./samples.yml --filter "*gif*"
  faas-cli remove -f ./samples.yml --regex "fn[0-9]_.*"
  faas-cli remove url-ping
  faas-cli remove img2ansi --gateway==http://remote-site.com:8080`,
	Run: runDelete,
}

func runDelete(cmd *cobra.Command, args []string) {
	var services stack.Services
	if len(yamlFile) > 0 {
		parsedServices, err := stack.ParseYAMLFile(yamlFile, regex, filter)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		if parsedServices != nil {
			services = *parsedServices
		}
	}

	if len(services.Functions) > 0 {
		if len(services.Provider.Network) == 0 {
			services.Provider.Network = defaultNetwork
		}

		for k, function := range services.Functions {
			function.Name = k
			fmt.Printf("Deleting: %s.\n", function.Name)

			proxy.DeleteFunction(services.Provider.GatewayURL, function.Name)
		}
	} else {
		if len(args) < 1 {
			fmt.Println("Please provide the name of a function to delete")
			return
		}

		functionName = args[0]
		fmt.Printf("Deleting: %s.\n", functionName)
		proxy.DeleteFunction(gateway, functionName)
	}
}
