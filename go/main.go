package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/glifio/go-pools/constants"
	"github.com/glifio/go-pools/deploy"
	"github.com/glifio/go-pools/sdk"
)

var PRIVATE_KEY = ""

func main() {
	example := flag.Uint("example", 0, "Example to run (1 or 2)")
	flag.Parse()

	switch *example {
	case 1:
		deployCounter()
	case 2:
		getCount()
	default:
		fmt.Println("No example selected")
	}
}

func deployCounter() {
	sdk, err := sdk.New(context.Background(), big.NewInt(constants.MainnetChainID), deploy.Extern)
	if err != nil {
		log.Fatalf("Failed to initialize pools sdk %s", err)
	}

	ethClient, err := sdk.Extern().ConnectEthClient()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyECDSA, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println("Deploying Counter contract from address: ", fromAddress.Hex())

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(constants.MainnetChainID))
	if err != nil {
		log.Fatal(err)
	}

	address, tx, _, err := DeployCounter(auth, ethClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deployed Counter to %s, tx: %s\n", address.Hex(), tx.Hash().Hex())
	log.Printf("Confirming deployment...")

	for {
		time.Sleep(time.Millisecond * 5000)

		tx, err := ethClient.TransactionReceipt(context.TODO(), tx.Hash())
		if err == nil && tx != nil {
			log.Printf("Counter deploy confirmed %s\n", address.Hex())
			return
		}
	}
}

func getCount() {
	sdk, err := sdk.New(context.Background(), big.NewInt(constants.MainnetChainID), deploy.Extern)
	if err != nil {
		log.Fatalf("Failed to initialize pools sdk %s", err)
	}

	ethClient, err := sdk.Extern().ConnectEthClient()
	if err != nil {
		log.Fatal(err)
	}

	counterAddress := "0x47bE44CdB532634e4d454d7F2Ff829DfF1560D13"
	counter, err := NewCounterCaller(common.HexToAddress(counterAddress), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	count, err := counter.Number(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Counter count: %d\n", count)
}
