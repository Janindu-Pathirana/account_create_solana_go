package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func create_new_account() string {
	// create new account
	newAccount := types.NewAccount()
	fmt.Println(newAccount.PublicKey.ToBase58())
	fmt.Println(newAccount.PrivateKey)

	// recover account by its private key
	recoverAccount, err := types.AccountFromBytes(
		newAccount.PrivateKey,
	)
	if err != nil {
		log.Fatalf("failed to retrieve account from bytes, err: %v", err)
	}
	fmt.Println("account : ", base58.Encode(recoverAccount.PrivateKey))

	return string(base58.Encode(recoverAccount.PrivateKey))
}

func fund_account(_account string, c *client.Client) {

	account, err := types.AccountFromBase58(_account)
	if err != nil {
		log.Fatalln("error: ", err)
	}

	txhash, err := c.RequestAirdrop(
		context.Background(),
		account.PublicKey.ToBase58(),
		1e9, // 1 SOL = 10^9 lamports
	)
	if err != nil {
		log.Fatalf("failed to request airdrop, err: %v", err)
	}

	fmt.Println("txhash:", txhash)

}

func check_balance(_account string, c *client.Client) {

	account, err := types.AccountFromBase58(_account)
	if err != nil {
		log.Fatalln("error: ", err)
	}

	balance, err := c.GetBalance(
		context.Background(),
		account.PublicKey.ToBase58(),
	)
	if err != nil {
		log.Fatalln("get balance error", err)
	}
	fmt.Println("balance : ", balance)
}

func transfer(_from string, _to string, c *client.Client) {

	from, from_error := types.AccountFromBase58(_from)
	to, to_error := types.AccountFromBase58(_to)

	if from_error != nil || to_error != nil {
		log.Fatalln("from Error : ", from_error, "\n", "to_error : ", to_error)
	}

	// to fetch recent blockhash
	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	// create a message
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        from.PublicKey,
		RecentBlockhash: res.Blockhash, // recent blockhash
		Instructions: []types.Instruction{
			sysprog.Transfer(sysprog.TransferParam{
				From:   from.PublicKey, // from
				To:     to.PublicKey,   // to
				Amount: 1000,           //  SOL val
			}),
		},
	})

	// create tx by message + signer
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{from},
	})
	if err != nil {
		log.Fatalf("failed to new transaction, err: %v", err)
	}

	// send tx
	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	log.Println("txhash:", txhash)

}

func main() {

	c := client.NewClient(rpc.DevnetRPCEndpoint)

	// account1 := create_new_account()
	// account2 := create_new_account()

	account1 := "3ztwHpJiLjE7gngFG2yzbbeEnrFGvAMi6f8SaR7XDpRGnCxSaz6UnBoe5ppP4h4Q7bK6zpKFnfwGZPMRXA81yQGJ"
	account2 := "276up8ht5MCgzkF8rxXANMtacSKvyoK8aPTjPideX9LjnwWUdWXZGWFmP9B1euvyVGyjozLXhPFedQBeYmBCi6MQ"

	// fund_account(account1, c)
	// fund_account(account2, c)

	fmt.Println("account 1")
	check_balance(account1, c)
	fmt.Println("account 2")
	check_balance(account2, c)

	transfer(account1, account2, c)

	fmt.Println("account 1")
	check_balance(account1, c)
	fmt.Println("account 2")
	check_balance(account2, c)

}
