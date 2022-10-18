package retriever

// リトリーバル方式の雑談システムをインターフェイスで抽象化する
type Retriever interface {
	Retrieve(dict map[string]string, req string) string
}
