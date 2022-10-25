---
marp: true
---

# 2022年度研究プロジェクト 前半発表

## リトリーバル方式の雑談システムの試作

情報通信系3年 河合輝良 (担当教員: 船越孝太郎)

---

## リトリーバル方式

- 現在の状態に基づいて、システムの発話DBの中から応答を選んで返す方式
- ここでは、直前の入力発話を現在の状態とみなす
  - 過去の入力は考慮しない

---

## 発話DB

- ある発話とその発話に対する応答のペアを記録したもの
- 対話破綻検出チャレンジ(2015年)の雑談対話コーパスを使用
  - <https://sites.google.com/site/dialoguebreakdowndetection/chat-dialogue-corpus>

zipファイルの中身↓ (重要箇所のみ)

```text
chat-dialogue-corpus
├── json
│  ├── init100
│  │  ├── 1407219916.log.json
│  │  └── ...
│  └── rest1046
│     ├── 1404365812.log.json
│     └── ...
└── ...
```

---

## 発話DB - コーパスの加工

各JSONファイルの中身↓ (重要箇所のみ)

```json
{
  "turns": [
    { "utterance": "今日は最高気温36度だって。暑いねえ" },
    { "utterance": "最高気温は１７度が予想されます？？" },
    { "utterance": "いやいや猛暑ですよ" },
    ...
  ]
}
```

発話とその応答がペアになるように加工↓ (実装ではJSONではなくGoのmapを使用)

```json
{
  "今日は最高気温36度だって。暑いねえ": "最高気温は１７度が予想されます？？",
  "最高気温は１７度が予想されます？？": "いやいや猛暑ですよ",
  ...
}
```

---

## 雑談システムの実装

- 応答選択手法(発話に対して最適な応答をDBから選択する手法)が必要
- どのように選択すればよいか？ = どのような評価指標を用いるか？
  - 完全一致
  - 編集距離
  - TF\*IDF指標
  - word2vec
  - BERT
  - ...

---

## 応答選択手法 - 完全一致

- 発話とDBの中の発話が完全一致しているものを選択
- なければ`"I don't know."`

---

## 応答選択手法 - 編集距離

- 編集距離: ある文字列から別の文字列への変換手順の最小回数
  - 1回の変換には挿入、削除、置換の3つの操作のいずれかを用いる
  - 「けんぷろ」と「すたんぷ」の編集距離は3
    - けんぷろ→けんぷ(削除)→たんぷ(置換)→すたんぷ(挿入)
- 発話とDBの中の発話の編集距離が最小になるものを選択

---

## 応答選択手法 - TF\*IDF指標

文書A: 「犬/と/猫/なら/犬/派/です」文書B: 「私/は/人間/です」

- TF\*IDF: TF(Term Frequency)とIDF(Inverse Document Frequency)の積
  - TF\*IDF(A, "犬") = 0.286 \* 0.176 = 0.0503
  - **その文書に出現する頻度が高く、かつ他の文書には出現しない単語は重要**
- TF: (文書中の単語`w`の出現回数) / (文書中の単語の総数)
  - TF(A, "犬") = 2 / 7 = 0.286
  - その文書に出現する頻度が高い単語は重要
- IDF: log((1+総文書数) / (1+単語`w`を含む文書数))
  - IDF([A,B], "です") = log((1+2) / (1+2)) = 0
  - どの文書にも出現する単語は重要ではない
- 発話とDBの中の発話のTF\*IDFのコサイン類似度が最大になるものを選択

---

# 実装

---

## 完成後イメージ

```bash
$ ./research-project1
# ? Which method do you want to use?:
#     Exact Match
#   ▸ Edit Distance
#     TF*IDF
✔ Edit Distance
Input your message: こんにちは
response: こんちわー

Input your message: おはよう、元気ですか？█
response: 元気です

...
```

---

## 方針

- レポジトリ: <https://github.com/Ras96/research-project1>
- Go言語を使用
- [manifoldco/promptui](https://github.com/manifoldco/promptui)を使いインタラクティブなCLIを作成
- 指標の計算には外部パッケージ(ライブラリ)を用いる → 次ページ

---

## 使用パッケージ (Github)

- [agnivade/levenshtein](https://github.com/agnivade/levenshtein)
  - 編集距離の計算
- [ikawaha/kagome](https://github.com/ikawaha/kagome)
  - 日本語文書の形態素解析 (文書を単語に分割)
  - 辞書はmecab-ipadic-2.7.0-20070801を使用
  - TF\*IDF指標の計算に使用
- [wilcosheh/tfidf](https://github.com/wilcosheh/tfidf)
  - TF\*IDF指標の計算
  - 発話とのコサイン類似度の計算

---

## Goによる実装

```go
package main

import "fmt"

func main() {
  // 応答選択の基準となる辞書を作成
  dict := makeResponseDictionary()
  // プロンプトから応答選択手法を決定
  r := selectRetrieverMethodInPrompt(dict)

  for {
    // プロンプトから入力文字列を決定
    req := getRequestMessageInPrompt()
    // 選択した応答選択手法で最も適した応答を出力
    res := r.Retrieve(req) // 選択した手法のRetrieveメソッドが呼び出される
    fmt.Println("response:", res)
  }
}
```

---

## `r.Retrieve(req)`の定義

> retrieve:〈情報を〉引き出す，検索する.
> リトリーバル(retrival)の動詞形
> (出典: <https://ejje.weblio.jp/>)

```go
package retriever

type Retriever interface {
  Retrieve(req string) string
}
```

インターフェイスのメソッドとして定義されている。
各手法で`Retrieve`メソッドを実装することで手法の切り替えが可能になる。

---

## `r.Retrieve(req)` - 完全一致編

```go
package exactmatch

type exactmatchRetriever struct {
  dict retriever.Dictionary
}

func (r *exactmatchRetriever) Retrieve(req string) string {
  if res, ok := r.dict[req]; ok {
    return res
  }

  return "I don't know."
}
```

---

## `r.Retrieve(req)` - 編集距離編

```go
package editdistance

func (r *editDistanceRetriever) Retrieve(req string) string {
  var (
    minDist int = 1e9
    bestRes string
    ref     string
  )

  for k, v := range r.dict {
    d := levenshtein.ComputeDistance(k, req) // 編集距離を計算
    if d < minDist { // 最小値を更新
      minDist = d; bestRes = v; ref = k
    }
  }

  return bestRes
}
```

---

## `r.Retrieve(req)` - TF\*IDF編

```go
package tfidf

func (r *tfIdfRetriever) Retrieve(req string) string {
  reqW := r.f.Cal(req) // 発話のTF*IDFを計算

  maxScore := 0.0
  maxDoc := ""
  for doc := range r.dict {
    docW := r.f.Cal(doc) // 応答のTF*IDFを計算
    score := similarity.Cosine(docW, reqW) // コサイン類似度を計算
    if score > maxScore { // 最大値を更新
      maxScore = score; maxDoc = doc
    }
  }

  return r.dict[maxDoc]
}
```

---

## 課題

- 今回用いた外部パッケージには最終更新が5年以上前のものもあるので自分で実装出来ればなお良かった
  - Goは後方互換性が保たれているので、バージョンが上がっても基本問題なく使うことはできる
- 指導教員からはさらにモダンな方法も紹介されていたが実装が間に合わなかった
  - word2vecやBERTなどのモデルを用いた方法
  - TF\*IDFでコサイン類似度を用いた方法を試せたのはよかった

---

# ご清聴ありがとうございました
