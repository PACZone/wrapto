package wallet

import (
	"log"
	"os"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	pWallet "github.com/pactus-project/pactus/wallet"
)

type Wallet struct {
	address  string
	password string
	wallet   *pWallet.Wallet
}

func Open(path, addr, rpcURL, pass string) *Wallet {
	if doesWalletExist(path) {
		wt, err := pWallet.Open(path, true)
		if err != nil {
			log.Fatal("error opening existing wallet", "err", err)
		}

		err = wt.Connect(rpcURL)
		if err != nil {
			log.Fatal("error establishing connection", "err", err)
		}

		return &Wallet{
			wallet:   wt,
			address:  addr,
			password: pass,
		}
	}
	log.Panic("wallet is required")

	return nil
}

func (w *Wallet) TransferTransaction(toAddress, memo string, amount int64) (string, error) {
	fee, err := w.wallet.CalculateFee(amount, payload.TypeTransfer)
	if err != nil {
		return "", err
	}

	opts := []pWallet.TxOption{
		pWallet.OptionFee(fee),
		pWallet.OptionMemo(memo),
	}

	tx, err := w.wallet.MakeTransferTx(w.address, toAddress, amount, opts...)
	if err != nil {
		log.Print("error creating transfer transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))

		return "", err
	}

	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Print("error signing transfer transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))

		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Print("error broadcasting transfer transaction", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))

		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		log.Print("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", util.ChangeToCoin(amount))
	}

	return res, nil
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Balance() int64 {
	balance, _ := w.wallet.Balance(w.address)

	return balance
}

func IsValidData(address, pubKey string) bool {
	addr, err := crypto.AddressFromString(address)
	if err != nil {
		return false
	}
	pub, err := bls.PublicKeyFromString(pubKey)
	if err != nil {
		return false
	}
	err = pub.VerifyAddress(addr)

	return err == nil
}

// function to check if file exists.
func doesWalletExist(fileName string) bool {
	_, err := os.Stat(fileName)

	return !os.IsNotExist(err)
}
