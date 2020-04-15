// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coinbase/rosetta-sdk-go/fetcher"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	// serverURL is the URL of a Rosetta Server.
	serverURL = "http://localhost:8080"
)

func main() {
	ctx := context.Background()

	// Step 1: Create a new fetcher
	newFetcher := fetcher.New(
		ctx,
		serverURL,
	)

	// Step 2: Initialize the fetcher's asserter
	//
	// Behind the scenes this makes a call to get the
	// network status and uses the response to inform
	// the asserter what are valid responses.
	primaryNetwork, networkStatus, err := newFetcher.InitializeAsserter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Step 3: Print the primary network and network status
	prettyPrimaryNetwork, err := json.MarshalIndent(
		primaryNetwork,
		"",
		" ",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Primary Network: %s\n", string(prettyPrimaryNetwork))

	prettyNetworkStatus, err := json.MarshalIndent(
		networkStatus,
		"",
		" ",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Network Status: %s\n", string(prettyNetworkStatus))

	// Step 4: Fetch the current block with retries (automatically
	// asserted for correctness)
	//
	// It is important to note that this assertion only ensures
	// required fields are populated and that operations
	// in the block only use types and statuses that were
	// provided in the networkStatusResponse. To run more
	// intensive validation, use the Rosetta Validator. It
	// can be found at: https://github.com/coinbase/rosetta-validator
	//
	// On another note, notice that fetcher.BlockRetry
	// automatically fetches all transactions that are
	// returned in BlockResponse.OtherTransactions. If you use
	// the client directly, you will need to implement a mechanism
	// to fully populate the block by fetching all these
	// transactions.
	block, err := newFetcher.BlockRetry(
		ctx,
		primaryNetwork,
		types.ConstructPartialBlockIdentifier(
			networkStatus.CurrentBlockIdentifier,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Step 5: Print the block
	prettyBlock, err := json.MarshalIndent(
		block,
		"",
		" ",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current Block: %s\n", string(prettyBlock))

	// Step 6: Get a range of blocks
	blockMap, err := newFetcher.BlockRange(
		ctx,
		primaryNetwork,
		networkStatus.GenesisBlockIdentifier.Index,
		networkStatus.GenesisBlockIdentifier.Index+10,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Step 7: Print the block range
	prettyBlockRange, err := json.MarshalIndent(
		blockMap,
		"",
		" ",
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block Range: %s\n", string(prettyBlockRange))
}
