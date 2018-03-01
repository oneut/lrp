# **WIP**

Exampleで動作はするようにしていますが、まだ作り込みが足りません

# LRP

LRPはgolang製のLive Reload Proxyです。

|OS        |Status    |
|----------|----------|
|Mac|OK|
|Windows|OK|
|Linux|?|

# モチベーション
+ livereloadを簡単にしたい
+ 汎用的なlivereload環境を作りたい
    + 特定の言語に依存しないLiveReload環境を作りたい
    + 複数の言語を組み合わせた環境に対応したい
+ コンパイルが必要な言語も自動的に再コンパイルしたい
+ 設定を簡単にしたい
+ golangで何かを作ってみたかった
    + 初golang
+ mattnさん作のgoemonがあるが、プロキシサーバは含まれていなかったので作ってみることにした

# 特徴
+ プロキシサーバを経由してLiveReloadを行う
    + ブラウザに対してエクステンションを必要としない
    + htmlに特定のスクリプトを差し込む必要がない
+ 設定をyamlで定義できる
+ ファイル監視ができる
    + 変更を検知してLireReloadイベントを発火が可能
    + ファイル作成
    + ファイル更新
    + ファイル削除
    + ディレクトリ作成
    + ディレクトリ削除
    + 除外設定
+ コマンドを管理できる
    + restartを簡単に行える
    + 標準出力のワードでLiveReloadイベントの発火が可能	
+ ポーリング機能はない
    + ローカルで動かすことを目的としているので必要がない
    + 必要そうなら実装を検討する
        + 候補: [radovskyb/watcher](https://github.com/radovskyb/watcher)

# Install
```
go get -u github.com/oneut/lrp
```

# 実行
```
lrp start
```

# YAML Example
実行するためには`lrp.yml`の定義が必須です。

```
proxy_host: "localhost:9000"
source_host: "localhost:8080"
tasks:
  web:
    aggregate_timeout: 300
    commands:
      go:
        execute: go run main.go
        needs_restart: true
    monitor:
      paths:
        - ./view
  js:
    commands:
      webpack:
        execute: npm start --prefix ./public/assets
        needs_restart: false
        watch_stdout:
          - bundle.js
      test:
        execute: npm test --prefix ./public/assets
        needs_restart: true
    monitor:
      paths:
        - ./public/assets
      ignore:
        - node_modules
```

# オプション
完成したら書く

# todo
+ test
+ リファクタリング
    + ファイル分割したい
    + ファイルや変数名の命名が雑
    + テストを考慮した実装になっていない
+ gRPCを使用した変更検知
+ `write tcp 127.0.0.1:9000->127.0.0.1:49720: write: broken pipe`がgo-livereloadのwebsocketで発生している・・・。
    + livereload自体を自作しないと制御できないので悩む
+ monitor定義で、イベント発火したファイル名を動的に取得
    + テストで特定ファイルのテスト実行が簡単にできそう
    + 正規表現でマッチ出来たら良さそう
