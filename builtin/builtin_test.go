package builtin

import (
	"testing"
)

func TestBuiltinFunctionInterface(t *testing.T) {
	// BuiltinFunction インターフェースが存在することを確認する
	// インターフェースの型だけを確認（nilポインタを避ける）
	var fn BuiltinFunction

	// インターフェースがnilであることを確認（コンパイルが通ることを確認）
	if fn != nil {
		t.Fatal("fn should be nil")
	}
}
