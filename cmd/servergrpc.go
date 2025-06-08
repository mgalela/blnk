/*
Copyright 2024 Blnk Finance Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"log"

	"github.com/jerry-enebeli/blnk/config"
	"github.com/spf13/cobra"
)

/*
serverCommands returns the Cobra command responsible for starting the Blnk GRPC server.
It sets up the API routes, traces, and TypeSense client before launching the server.
*/
func serverGrpcCommands(b *blnkInstance) *cobra.Command {
	// Define the `startgrpc` command for starting the server
	cmd := &cobra.Command{
		Use:   "startgrpc",
		Short: "start grpc blnk server", // Short description of the command
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			// Load configuration
			cfg, err := config.Fetch()
			if err != nil {
				log.Println(err)
			}
	
			// Start server
			if err := startGrpcServer(ctx, b.blnk, cfg.Server); err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
