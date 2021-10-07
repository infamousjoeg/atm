package atm

import (
	"fmt"
	"os"
	"time"
)

type IATM interface {
	accountDB() IAccountDB
	transactionDB() ITransactionDB
	balance() float64
	Authorize(int, string) bool
	Withdraw(Account, int) (bool, error)
	Deposit(Account, float64) error
	Balance(Account)
	History(Account) []Transaction
	Logout(bool) error
	End()
}

type ATM struct {
	AccountDB     IAccountDB
	TransactionDB ITransactionDB
	ATMBalance    float64
	Session       *Session
}

func (atm *ATM) accountDB() IAccountDB {
	return atm.AccountDB
}

func (atm *ATM) transactionDB() ITransactionDB {
	return atm.TransactionDB
}

func (atm *ATM) balance() float64 {
	return atm.ATMBalance
}

func (atm *ATM) Authorize(accountID int, accountPIN string) bool {
	account, err := atm.accountDB().Get(accountID)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if account.PIN != accountPIN {
		return false
	}

	atm.Session.Authorize(accountID)

	return true
}

func (atm *ATM) Withdraw(accountID int, amount int) (bool, error) {
	overdrawn := false
	err := atm.Session.Valid(accountID)
	if err != nil {
		return overdrawn, err
	}

	account, err := atm.accountDB().Get(accountID)
	// If account does not exist then return error
	if err != nil {
		return overdrawn, err
	}

	if atm.balance() == 0 {
		return overdrawn, ErrWithdrawATMNoFunds
	}

	// Amount is not multiple of 20
	if amount%20 != 0 {
		return overdrawn, ErrWithdrawAmountNoMultipleOf20
	}

	if amount > int(atm.balance()) {
		return overdrawn, ErrWithdrawATMInsufficientFunds
	}

	if account.Balance < 0 {
		return overdrawn, ErrWithdrawAccountOverdrawn
	}

	newBalance := account.Balance - float64(amount)
	if newBalance < 0 {
		overdrawn = true
		newBalance -= 5
	}

	transaction := Transaction{
		AccountID: account.AccountID,
		DateTime:  time.Now().Unix(),
		Amount:    float64(amount * -1),
		Balance:   newBalance,
	}
	atm.transactionDB().Set(transaction)

	account.Balance = newBalance
	atm.accountDB().Set(account)
	atm.ATMBalance = atm.ATMBalance - float64(amount)
	return overdrawn, nil
}

func (atm *ATM) Deposit(accountID int, amount float64) error {
	err := atm.Session.Valid(accountID)
	if err != nil {
		return err
	}

	account, err := atm.accountDB().Get(accountID)
	if err != nil {
		return err
	}

	newBalance := account.Balance + amount
	transaction := Transaction{
		AccountID: account.AccountID,
		DateTime:  time.Now().Unix(),
		Amount:    amount,
		Balance:   newBalance,
	}
	atm.transactionDB().Set(transaction)

	account.Balance = newBalance
	atm.accountDB().Set(account)

	return nil
}

func (atm *ATM) Balance(accountID int) (float64, error) {
	err := atm.Session.Valid(accountID)
	if err != nil {
		return 0.00, err
	}

	account, err := atm.accountDB().Get(accountID)
	return account.Balance, err
}

func (atm *ATM) History(accountID int) []Transaction {
	err := atm.Session.Valid(accountID)
	if err != nil {
		return []Transaction{}
	}
	transactions, err := atm.transactionDB().Get(accountID)
	if err != nil {
		return []Transaction{}
	}

	return transactions
}

func (atm *ATM) Logout() error {
	return atm.Session.LogOut()
}

func (atm *ATM) End() {
	os.Exit(0)
}
