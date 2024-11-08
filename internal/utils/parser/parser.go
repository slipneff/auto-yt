package parser

import (
	"bufio"
	"os"
	"strings"
)

type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Recovery string `json:"recovery"`
}

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

func (a *Accounts) Add(account Account) {
	a.Accounts = append(a.Accounts, account)
}

func (a *Accounts) Parse(data string) Account {
	d := strings.Split(data, ":")
	return Account{Email: d[0], Password: d[1], Recovery: d[2]}
}


func ReadFile(file string) (*Accounts, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	accounts := &Accounts{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		accounts.Add(accounts.Parse(scanner.Text()))
	}
	return accounts, nil
}
