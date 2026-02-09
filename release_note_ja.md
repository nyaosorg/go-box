リリースノート
=============
( [English](./release_note.md) / **Japanese** )

v3.1.1
------
Feb 10, 2026

- `Meta`+`Ctrl`+`G` と `Meta`+`Ctrl`+`C` も `Ctrl`+`G` と同様に機能するようにした。(#11)

v3.1.0
------
Feb 10, 2026

- go-ttyadapter を v0.3.0 へ更新。`tty8` のかわりに `go-ttyadapter/tty8pe` を使用 (#5)
  - `Escape`キーは入力シーケンスの分断防止のためプリフィックスとして扱われるようになった。今後、キャンセルには `Ctrl`+`G` もしくは `Ctrl`+`C` を使用のこと (#7)
- golang.org/x/sys を v0.29.0 から v0.30.0 へ更新 (#6)
- clipperhouse/uax29 を v2.2.0 から v2.6.0 へ更新 (#6)

v3.0.0
------
Nov 1, 2025

- メジャーバージョンを **v3** に更新。
- Deprecated としていた関数・メソッドを削除。
  - `BoxT`, `Choce`, `ChoiceMulti`, `ChooseMulti`, `Choose`, `PrintNoLastLineFeed`, `Print`
- 関数・メソッド名を整理。
  - `PrintX` → `Print` : 最終行で改行する整列表示関数・メソッド
  - `PrintNoLastLineFeedX` → `Print` : 最終行で改行しない整列表示関数・メソッド
  - `NewBox` → `New` : `Box` 構造体のコンストラクタ
  - `AnsiCutter` → 非公開化 : エスケープシーケンス除去用の正規表現
- ソース構成を整理。
  - `box.go` : `Box` 型関連
  - `main.go` : グローバル定義
  - `internal/ansi` : エスケープシーケンス関連
  - `internal/key` : キーコード関連
- 依存パッケージを最新化。
  - `mattn/go-tty` v0.0.7
  - `mattn/go-colorable` v0.1.14
  - `mattn/go-runewidth` v0.0.19
- パッケージ全体で実質的に利用していなかった `context.Context` パラメータを削除。
- 表示レイアウト機能を `Box` 型から分離し、`grid` サブパッケージの `Grid` 型として独立。
- 端末入力を仮想化し、任意の入力パッケージを利用可能にした。
  `Box.Tty` フィールドに以下の値を設定後、`Init` メソッドを呼び出すことで切り替えられる。
  `New` 関数は従来どおり `tty8` をデフォルトとして使用。
  - `&tty8.Tty{}`  ("[github.com/nyaosorg/go-ttyadapter][adapter]/tty8") → [github.com/mattn/go-tty][go-tty] を利用
  - `&tty10.Tty{}` ("[github.com/nyaosorg/go-ttyadapter][adapter]/tty10") → [golang.org/x/term][xterm] を利用
  - `&auto.Pilot{Text: []string{入力列}}` ("[github.com/nyaosorg/go-ttyadapter][adapter]/auto") → 擬似端末入力
- Windows Terminal が Unicode の *Ambiguous-width* 文字を常に半角として扱う問題に対して、
  `github.com/mattn/go-runewidth` v0.0.17 以降が自動検出に対応したため、当パッケージ側の補正コードを削除。
  （参考: [mattn/go-runewidth#85](https://github.com/mattn/go-runewidth/pull/85)）

[adapter]: https://github.com/nyaosorg/go-ttyadapter
[go-tty]: https://github.com/mattn/go-tty
[xterm]: https://pkg.go.dev/golang.org/x/term
