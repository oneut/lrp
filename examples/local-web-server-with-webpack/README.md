# Local Web Server with webpack Example

This is a example that monitors the standard output of webpack and performs Live Reload to the local web server.

![local web server with webpack example](https://raw.githubusercontent.com/oneut/lrp/master/examples/local-web-server-with-webpack/local-web-server-with-webpack.gif)

## Doing
+ Live Reload
    + Source website is local web server.
    + Set webpack task.
        + Command
            + Execute webpack watch.
            + Watch stdouts.
                + `Entrypoint main = bundle.js` 
                + If stdout matches, `lrp` fires Live Reload.
    + Set local web server task.
        + Command
            + Start local web server.
                + use [http-server](https://github.com/indexzero/http-server)

## Ready
```
npm install
```

## Start
```
lrp start
```