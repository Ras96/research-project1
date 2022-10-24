package retriever

type Dictionary map[string]string

// リトリーバル方式の雑談システムをインターフェイスで抽象化する
type Retriever interface {
	Retrieve(req string) string
}
