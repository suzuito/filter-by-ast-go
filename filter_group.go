package main

// TODO: ↑ぶっちゃけパッケージ名が微妙。後で変えたい。

// FilterGroupは、トピックに対するフィルタ。
type FilterGroup interface {
	Eval(*Post) (bool, error)
	String() string
}
