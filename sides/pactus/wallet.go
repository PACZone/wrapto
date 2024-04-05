package pactus

import (
	"os"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	pWallet "github.com/pactus-project/pactus/wallet"
)

type Wallet struct {
	address  string
	password string
	wallet   pWallet.Wallet
}

func openWallet(path, addr, rpcURL, pass string) (*Wallet, error) {
	if !doesWalletExist(path) {
		return nil, WalletNotExistError{
			path: path,
		}
	}

	wt, err := pWallet.Open(path, true)
	if err != nil {
		return nil, err
	}

	err = wt.Connect(rpcURL)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		wallet:   *wt,
		address:  addr,
		password: pass,
	}, nil
}

func (w *Wallet) TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error) {
	fee, err := w.wallet.CalculateFee(amt, payload.TypeTransfer)
	if err != nil {
		return "", err
	}

	opts := []pWallet.TxOption{
		pWallet.OptionFee(fee),
		pWallet.OptionMemo(memo),
	}

	tx, err := w.wallet.MakeTransferTx(w.address, toAddress, amt, opts...)
	if err != nil {
		return "", err
	}

	// sign transaction
	err = w.wallet.SignTransaction(w.password, tx)
	if err != nil {
		return "", err
	}

	// broadcast transaction
	res, err := w.wallet.BroadcastTransaction(tx)
	if err != nil {
		return "", err
	}

	err = w.wallet.Save()
	if err != nil {
		return "", SaveWalletError{}
	}

	return res, nil
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Balance() amount.Amount {
	blnc, err := w.wallet.Balance(w.address)
	if err != nil {
		return 0
	}

	return blnc
}

func doesWalletExist(fileName string) bool {
	_, err := os.Stat(fileName)

	return !os.IsNotExist(err)
}
