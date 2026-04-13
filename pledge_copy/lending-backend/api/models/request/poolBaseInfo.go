package request

type PoolBaseInfo struct {
	ChainId int `form:"chain_id" json:"chain_id" validate:"required"`
}

type PoolDataInfo struct {
	ChainId int `form:"chain_id" json:"chain_id" validate:"required"`
}

type TokenList struct {
	ChainId int `form:"chain_id" json:"chain_id" validate:"required"`
}

type Search struct {
	ChainID         int    `json:"chain_id" validate:"required"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
	LendTokenSymbol string `json:"lend_token_symbol"`
	State           string `json:"state"`
}
