package builtin

// BuiltinFunction は組み込み関数のインターフェース
type BuiltinFunction interface {
	Name() string
	Execute(input map[string]interface{}) (string, error)
}

// Registry は組み込み関数を管理するレジストリ
type Registry struct {
	functions map[string]BuiltinFunction
}

// NewRegistry は新しいRegistryを作成し、組み込み関数を登録
func NewRegistry() *Registry {
	r := &Registry{
		functions: make(map[string]BuiltinFunction),
	}

	// 組み込み関数を登録
	r.Register(ContextPercentFunction{})
	r.Register(ContextUsedFunction{})

	return r
}

// Register は組み込み関数を登録
func (r *Registry) Register(fn BuiltinFunction) {
	r.functions[fn.Name()] = fn
}

// Get は名前で組み込み関数を取得
func (r *Registry) Get(name string) BuiltinFunction {
	return r.functions[name]
}
