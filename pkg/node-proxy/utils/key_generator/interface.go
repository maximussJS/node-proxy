package key_generator

type KeyGenerator interface {
	GenerateJsonRpcKey(jsonrpc, method string, params []any) (string, error)
}
