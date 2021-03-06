# LRP

LRPはgolang製のLive Reload Proxyです。

|OS        |Status    |
|----------|----------|
|Mac|OK|
|Windows|OK|
|Linux|?|

# モチベーション
+ Live Reloadを簡単にしたい
    + 特定の言語に依存せず、汎用的に使える
    + 複数のプログラミング言語を考慮できる
    + コンパイルが必要な言語は自動的に再コンパイルできる
+ golangで何かを作ってみたかった
    + 初golangアプリケーション
+ 各Live Reloadツールのメリット・デメリットを検討した結果、機能をいいとこ取りすると便利になりそうだった
    + BrowserSync
        + 特定言語に依存していない
        + プロキシサーバとして動く
        + コマンドの起動を制御できない
    + LiveReloadX
        + 特定言語に依存していない
        + staticが便利
        + プロキシサーバとして動く
        + コマンドの起動を制御できない
    + goemon
        + 特定言語に依存していない
        + Yamlで定義できる
        + コマンドの起動を制御できる
        + プロキシサーバは含まれていなかったので、エクステンション等が必要

# 主な特徴
+ プロキシサーバを経由してLive Reloadを行う
    + Live Reload用のスクリプトファイルはプロキシサーバが自動的に付与する
      + ブラウザのエクステンションを必要としない
      + htmlに特定のスクリプトを差し込む必要がない
+ 設定をyamlで定義
+ ファイル監視
    + ファイル変更、ディレクトリ変更を検知してLive Reloadイベントを発火
    + 除外設定も可能
+ コマンド管理
    + Live Reloadイベント発火時にコマンドのrestartが可能
    + 標準出力のキーワードでLive Reloadイベントを発火	
+ 静的Webサーバ
    + 任意のディレクトリを静的Webサーバとして公開可能
    + 任意のディレクトリのファイルを優先させるスタブ機能として使用可能
+ ポーリング機能は今のところない
    + ローカルで動かすことを目的としているので優先度が低い
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

# Yaml Example
実行するためには`lrp.yml`の定義が必須です。

```
proxy:
  scheme: "https"
  host: "localhost:9000"
  staticPath: ./
source:
  scheme: "https"
  host: "localhost:8080"
tasks:
  web:
    aggregateTimeout: 300
    commands:
      go:
        executes:
          - "go clean"
          - "go build -o main"
          - "./main"
        needsRestart: true
    monitor:
      paths:
        - ./view
  js:
    commands:
      webpack:
        executes:
          - npm start --prefix ./public/assets
        needsRestart: false
        watchStdouts:
          - bundle.js
      test:
        executes: 
          - npm test --prefix ./public/assets
        needsRestart: true
    monitor:
      paths:
        - ./public/assets
      ignores:
        - node_modules
```

# 設定例
+ [Local Web Server](https://github.com/oneut/lrp/tree/master/examples/local-web-server)
+ [Local Web Server with webpack](https://github.com/oneut/lrp/tree/master/examples/local-web-server-with-webpack)

# オプション
## proxy
Live Reloadのプロキシサーバを管理します。

## proxy.scheme
URLスキームを定義できます。任意の設定です。デフォルトは`http`です。

```
proxy:
  scheme: "https"
```

## proxy.host
Live Reloadのプロキシサーバのドメイン、ホスト名、ポートを指定します。デフォルトは`:9000`です。

```
proxy:
  host: ":9000"
```

```
proxy:
  host: "localhost:9000"
```

## proxy.staticPath
プロキシサーバ経由で静的ファイルにアクセスできます。任意の設定です。

```
proxy:
  staticPath: ./static
```

もし、proxy.staticPathとsource.hostで同じパスが存在した場合は、staticPathが優先されます。    
そのため、Webサーバに対してアクセスさえできれば、一部のファイルだけ`staticPath`に存在するローカルファイルに差し替えることができます。

## proxy.browserOpen
ブラウザの自動オープンを制御します。デフォルトは自動でブラウザがオープンします。
ブラウザを自動でオープンしたくない場合は、`none`を設定してください。

```
proxy:
  browserOpen: "none"
```

## source
Live Reloadを行うWebサーバを管理します。

## source.scheme
URLスキームを定義できます。任意の設定です。デフォルトは`http`です。

```
source:
  scheme: "https"
```

## source.host
Live Reloadを行うWebサーバのドメイン、ホスト名、ポートを指定します。
設定は任意です。

```
source:
  host: ":8080"
```

```
source:
  host: "localhost:8080"
```

`source.host`を設定した場合、ブラウザにアクセスした時には`proxy.host`の値に自動的に置換されます。

## source.replaces
プロキシを経由してアクセスしているサイトのHTMLデータを条件を指定して置換します。  
複数の設定が可能です。

## source.replaces.-.search
検索条件になるキーワードを指定します。

```
source:
  host: "localhost:8080"
  replaces:
    - search: "//cdn.example.com"
      replace: "//localhost:9000"
```

## source.replaces.-.replace
検索条件に対して置換されるキーワードを指定します

```
source:
  host: "localhost:8080"
  replaces:
    - search: "//cdn.example.com"
      replace: "//localhost:9000"
```

## source.replaces.-.regexp
検索条件に正規表現の使用を指定します。デフォルトは`false`です。

```
source:
  host: "localhost:8080"
  replaces:
    - search: "abc(.+)"
      replace: "$1"
      regexp: true
```

## tasks
Live Reloadを行うタスクを管理します。タスク名は任意で設定できます。
タスクごとに実行するコマンドとファイル監視を設定します。
コマンドとファイル監視は、タスク単位で制御され、複数のタスクを設定できます。

```
tasks:
  web:
    aggregateTimeout: 300
    commands:
      go:
        executes:
          - "go clean"
          - "go build -o main"
          - "./main"
        needsRestart: true
    monitor:
      paths:
        - ./view
```

## tasks.{task_name}.aggregateTimeout
最初のファイルが変更されてからの遅延を設定できます。
この設定により、この期間に行われた他の変更を1回の再構築に集約できます。
設定はミリ秒単位です。デフォルトは`300`です。

```
tasks:
  web:
    aggregateTimeout: 300
```

## tasks.{task_name}.commands
Live Reloadを行う際のコマンドを管理できます。コマンド名は任意で設定でき、複数設定が可能です。
コマンドごとに再起動、標準出力の監視を行います。

```
tasks:
  web:
    commands:
      go:
        executes: 
          - "go clean"
          - "go build -o main"
          - "./main"
      npm:
        executes:
          - npm start
```

## tasks.{task_name}.commands.{command_name}.executes
実行したいコマンドを設定します。複数設定可能です。タスクは逐次処理です。

```
tasks:
  web:
    commands:
      go:
        executes:
          - "go clean"
          - "go build -o main"
          - "./main"
```

## tasks.{task_name}.commands.{command_name}.needsRestart
ファイル監視で変更が起きた際にコマンドを再実行したい場合に設定します。デフォルトは`false`です。

```
tasks:
  web:
    commands:
      go:
        executes:
          - "go clean"
          - "go build -o main"
          - "./main"
        needsRestart: true
```

## tasks.{task_name}.commands.{command_name}.watchStdouts
コマンド実行時の標準出力を監視します。
設定したキーワードによってLive Reloadイベントを発火します。
`needsRestart=true`で再起動した際に、`watchStdouts`で指定したキーワードがあると無限ループします。気をつけましょう。

```
tasks:
  web:
    commands:
      npm:
        executes:
          - npm start
        needsRestart: false
        watchStdouts:
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
      ignores:
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

## tasks.{task_name}.monitor.ignores
除外するファイル、ディレクトリを指定します。
設定できるパターンは部分一致、パス指定が行なえます。
**また、*、?等のワイルドカードが使えますが、実装を調整中です。**

```
tasks:
  web:
    monitor:
      paths:
        - ./app
      ignores:
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
+ README.mdを英語にする。
+ 取得したHTMLの特定キーワードの置換処理
