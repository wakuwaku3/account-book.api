package entitieslog

// Lvl はログレベルです
type Lvl int

const (
	// Fatal はアプリケーションの継続が不可能な致命的な障害に適用します
	Fatal Lvl = iota
	// Error はリクエストの継続が不可能な致命的な障害に適用します
	Error
	// Warn は警告です
	Warn
	// Info は情報です
	Info
	// Debug はデバッグ時のみ出力する情報です
	Debug
	// Trace はデバッグ時のみ出力する追跡情報です
	Trace
)
