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
    + 初golangアプリケーション
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
## proxy_host
Live Reloadのプロキシサーバのドメイン、ホスト名、ポートを指定します。デフォルトは`localhost:9000`です。

+ `proxy_host: :9000`
+ `proxy_host: localhost:9000`

## source_host
Live Reloadを行うWebサーバのドメイン、ホスト名、ポートを指定します。デフォルトは`:8080`です。

+ `source_host: :8080`
+ `source_host: localhost:8080`

## tasks
Live Reloadを行うタスクを管理します。タスク名は任意で設定できます。
タスクごとに実行するコマンドとファイル監視を設定します。
コマンドとファイル監視は、タスク単位で制御されます。

```
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
```

## tasks.{task_name}.aggregate\_timeout
最初のファイルが変更されてからの遅延を設定できます。
この設定により、この期間に行われた他の変更を1回の再構築に集約できます。
設定はミリ秒単位です。デフォルトは`300`です。

```
tasks:
  web:
    aggregate_timeout: 300
```

## tasks.{task_name}.commands
Live Reloadを行う際のコマンドを管理できます。コマンド名は任意で設定でき、複数設定が可能です。
コマンドごとに再起動、標準出力の監視を行います。

```
tasks:
  web:
    commands:
      go:
        execute: go build
      npm:
        execute: npm start
```

## tasks.{task_name}.commands.{command_name}.execute
実行したいコマンドを設定します。

```
tasks:
  web:
    commands:
      go:
        execute: go build
```

## tasks.{task_name}.commands.{command_name}.needs_restart
ファイル監視で変更が起きた際にコマンドを再実行したい場合に設定します。デフォルトは`false`です。

```
tasks:
  web:
    commands:
      go:
        execute: go build
        needs_restart: true
```

## tasks.{task_name}.commands.{command_name}.watch_stdout
コマンド実行時の標準出力を監視します。
設定したキーワードによってLive Reloadイベントを発火します。
`needs_restart=true`で再起動した際に、`watch_stdout`で指定したキーワードがあると無限ループします。気をつけましょう。

```
tasks:
  web:
    commands:
      npm:
        execute: npm start
        needs_restart: false
        watch_stdout:
          - bundle.js
```

## tasks.{task_name}.monitor
ファイル監視を管理します。

```
tasks:
  web:
    monitor:
      paths:
        - ./view
      ignore:
        - node_modules
```

## tasks.{task_name}.monitor.paths
ファイルを監視するディレクトリを指定します。
ディレクトリは複数設定できます。

```
tasks:
  web:
    monitor:
      paths:
        - ./app
        - ./views
```

## tasks.{task_name}.monitor.ignore
除外するファイル、ディレクトリを指定します。
設定できるパターンは部分一致、パス指定が行なえます。
**また、*、?等のワイルドカードが使えますが、実装を調整中です。**

```
tasks:
  web:
    monitor:
      paths:
        - ./app
      ignore:
        - node_modules
        - ./node_modules
```



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
