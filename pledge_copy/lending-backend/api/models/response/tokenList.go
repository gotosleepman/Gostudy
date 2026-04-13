package response

import "time"

type TokenList struct {
	Name      string    `json:"name"`
	LogoURI   string    `json:"logoURI"`
	Timestamp time.Time `json:"timestamp"`
	Version   Version   `json:"version"`
	Tokens    []Token   `json:"tokens"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch"`
}

type Token struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
	ChainID  string `json:"chainId"`
	LogoURI  string `json:"logoURI"`
}
