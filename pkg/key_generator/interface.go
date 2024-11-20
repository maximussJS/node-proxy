package key_generator

type IKeyGenerator interface {
	GenerateJsonRpcKey(jsonrpc, method string, params []any) (string, error)
}
