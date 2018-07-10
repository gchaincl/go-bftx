// File: ./blockfreight/cmd/bftnode/bftnode.go
// Summary: Application code for Blockfreight™ | The blockchain of global freight.
// License: MIT License
// Company: Blockfreight, Inc.
// Author: Julian Nunez, Neil Tran, Julian Smith, Gian Felipe & contributors
// Site: https://blockfreight.com
// Support: <support@blockfreight.com>

// Copyright © 2017 Blockfreight, Inc. All Rights Reserved.

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
// OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// =================================================================================================================================================
// =================================================================================================================================================
//
// BBBBBBBBBBBb     lll                                kkk             ffff                         iii                  hhh            ttt
// BBBB``````BBBB   lll                                kkk            fff                           ```                  hhh            ttt
// BBBB      BBBB   lll      oooooo        ccccccc     kkk    kkkk  fffffff  rrr  rrr    eeeee      iii     gggggg ggg   hhh  hhhhh   tttttttt
// BBBBBBBBBBBB     lll    ooo    oooo    ccc    ccc   kkk   kkk    fffffff  rrrrrrrr eee    eeee   iii   gggg   ggggg   hhhh   hhhh  tttttttt
// BBBBBBBBBBBBBB   lll   ooo      ooo   ccc           kkkkkkk        fff    rrrr    eeeeeeeeeeeee  iii  gggg      ggg   hhh     hhh    ttt
// BBBB       BBB   lll   ooo      ooo   ccc           kkkk kkkk      fff    rrr     eeeeeeeeeeeee  iii   ggg      ggg   hhh     hhh    ttt
// BBBB      BBBB   lll   oooo    oooo   cccc    ccc   kkk   kkkk     fff    rrr      eee      eee  iii    ggg    gggg   hhh     hhh    tttt    ....
// BBBBBBBBBBBBB    lll     oooooooo       ccccccc     kkk     kkkk   fff    rrr       eeeeeeeee    iii     gggggg ggg   hhh     hhh     ttttt  ....
//                                                                                                        ggg      ggg
//   Blockfreight™ | The blockchain of global freight.                                                      ggggggggg
//
// =================================================================================================================================================
// =================================================================================================================================================

// Starts the Blockfreight™ Node to listen to all requests in the Blockfreight Network.
package main

import (
	// =======================
	// Golang Standard library
	// =======================
	// Implements command-line flag parsing.
	"fmt" // Implements formatted I/O with functions analogous to C's printf and scanf.
	"os"

	// ===============
	// Tendermint Core
	// ===============

	tmConfig "github.com/tendermint/tendermint/config"
	tmNode "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tmlibs/log"
	// ======================
	// Blockfreight™ packages
	// ======================
	"github.com/blockfreight/go-bftx/api/api"
	"github.com/blockfreight/go-bftx/lib/app/bft"
	// Implements the main functions to work with the Blockfreight™ Network.
)

var homeDir = os.Getenv("HOME")
var logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))

func BlockfreightAppClientCreator(addr, transport, dbDir string) proxy.ClientCreator {
	return proxy.NewLocalClientCreator(bft.NewBftApplication())
}

func main() {

	fmt.Println("Blockfreight™ Node")

	index := &tmConfig.TxIndexConfig{
		Indexer:      "kv",
		IndexTags:    "bftx.id",
		IndexAllTags: false,
	}

	config := tmConfig.DefaultConfig()

	config.P2P.Seeds = "0ce024c57fc1137bfbee70a1e520fba4c9163fbe@bftx0.blockfreight.net:8888,0537b4c4800b810858dc554e65f85b76217ff900@bftx1.blockfreight.net:8888,5a4833829cc5cec95a6194fb16e3ad75b605968b@bftx2.blockfreight.net:8888,5fe8f8847e4b87c6eea350bcd55269d3c492ffcb@bftx3.blockfreight.net:8888"
	config.Consensus.CreateEmptyBlocks = false

	config.TxIndex = index
	config.DBPath = homeDir + "/.blockfreight/config/bft-db"
	config.Genesis = homeDir + "/.blockfreight/config/genesis.json"
	config.PrivValidator = homeDir + "/.blockfreight/config/priv_validator.json"
	config.NodeKey = homeDir + "/.blockfreight/config/node_key.json"
	config.P2P.ListenAddress = "tcp://0.0.0.0:8888"

	logger.Info("Setting up config", "nodeInfo", config)

	node, err := tmNode.NewNode(config,
		privval.LoadOrGenFilePV(config.PrivValidatorFile()),
		BlockfreightAppClientCreator(config.ProxyApp, config.ABCI, config.DBDir()),
		tmNode.DefaultGenesisDocProviderFunc(config),
		tmNode.DefaultDBProvider,
		logger,
	)

	if err != nil {
		fmt.Errorf("Failed to create a node: %v", err)
	}

	if err = node.Start(); err != nil {
		fmt.Errorf("Failed to start node: %v", err)
	}

	logger.Info("Started node", "nodeInfo", node.Switch().NodeInfo())

	err = api.Start()
	if err != nil {
		logger.Error(err.Error())
	}

	// Trap signal, run forever.
	node.RunForever()

}

// =================================================
// Blockfreight™ | The blockchain of global freight.
// =================================================

// BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB
// BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB
// BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB
// BBBBBBB                    BBBBBBBBBBBBBBBBBBB
// BBBBBBB                       BBBBBBBBBBBBBBBB
// BBBBBBB                        BBBBBBBBBBBBBBB
// BBBBBBB       BBBBBBBBB        BBBBBBBBBBBBBBB
// BBBBBBB       BBBBBBBBB        BBBBBBBBBBBBBBB
// BBBBBBB       BBBBBBB         BBBBBBBBBBBBBBBB
// BBBBBBB                     BBBBBBBBBBBBBBBBBB
// BBBBBBB                        BBBBBBBBBBBBBBB
// BBBBBBB       BBBBBBBBBB        BBBBBBBBBBBBBB
// BBBBBBB       BBBBBBBBBBB       BBBBBBBBBBBBBB
// BBBBBBB       BBBBBBBBBB        BBBBBBBBBBBBBB
// BBBBBBB       BBBBBBBBB        BBB       BBBBB
// BBBBBBB                       BBBB       BBBBB
// BBBBBBB                    BBBBBBB       BBBBB
// BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB
// BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB
// BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB

// ==================================================
// Blockfreight™ | The blockchain for global freight.
// ==================================================
