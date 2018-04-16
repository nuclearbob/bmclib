// Copyright © 2018 Joel Rebello <joel.rebello@booking.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/ncode/bmcbutler/asset"
	"github.com/ncode/bmcbutler/butler"
	"github.com/ncode/bmcbutler/inventory"
	"github.com/ncode/bmcbutler/resource"
	"time"
	//"os"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply config to bmcs.",
	Run: func(cmd *cobra.Command, args []string) {
		apply()
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func apply() {

	// Read in declared resources
	resourceInstance := resource.Resource{Log: log}
	config := resourceInstance.ReadResources()

	// A channel to recieve inventory assets
	inventoryChan := make(chan []asset.Asset)
	inventoryInstance := inventory.Inventory{Log: log, BatchSize: 10, Channel: inventoryChan}

	//// Spawn a goroutine that returns a slice of assets over inventoryChan
	//// the number of assets in the slice is determined by the batch size.
	go inventoryInstance.AssetIter()

	// Spawn butlers to work
	butlerChan := make(chan butler.ButlerMsg)
	butlerInstance := butler.Butler{Log: log, SpawnCount: 5, Channel: butlerChan}
	go butlerInstance.Spawn()

	////iterate over assets, pass them to the spawned butlers
	for asset := range inventoryChan {
		butlerMsg := butler.ButlerMsg{Assets: asset, Config: config}
		butlerChan <- butlerMsg
		time.Sleep(5000 * time.Millisecond)
	}
}