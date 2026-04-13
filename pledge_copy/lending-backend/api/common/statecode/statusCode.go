package statecode

const (
	CommonSuccess       = 0
	CommonErrServerErr  = 10001
	ParameterEmptyErr   = 10002
	ChainIdEmpty        = 10003
	ChainIdErr          = 10004
)

const LangEn = 1

func GetMsg(code int, lang int) string {
	switch code {
	case CommonSuccess:
		return "success"
	case CommonErrServerErr:
		return "server error"
	case ParameterEmptyErr:
		return "parameter empty"
	case ChainIdEmpty:
		return "chain id empty"
	case ChainIdErr:
		return "chain id error"
	default:
		return "unknown"
	}
}
