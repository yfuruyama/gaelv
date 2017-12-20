# gaelv
gaelv is a log viewer for Google App Engine local development.

<img alt="screenshot" src="https://raw.github.com/addsict/gaelv/master/img/screenshot.png" width="600">

## Install
```
go get -u github.com/addsict/gaelv/...
```

## How to use

1. Run your app server (`dev_appserver.py`) with `--logs_path=</path/to/log.db>` option.
```
dev_appserver.py app.yaml --logs_path=/tmp/gaelog.db
```

2. Run `gaelv` with same `--logs_path` option.
```
gaelv --logs_path=/tmp/gaelog.db
```

3. Open http://localhost:9090/ on your browser.

## Console mode

If you want to see logs in the console instead of a browser, run gaelv with `--console` option.

```
gaelv --logs_path=/tmp/gaelog.db --console
```

<img alt="screenshot2" src="https://raw.github.com/addsict/gaelv/master/img/screenshot2.png" width="500">

## FAQ

### The latest logs doesn't appear immediately.
Unfortunately, the latest logs are buffered in the app engine log service for 5 seconds.  
There is no workaround for it now except modifying sdk source code, so please change the value `_MIN_COMMIT_INTERVAL` to `0` in the file `${SDK_ROOT_PATH}/platform/google_appengine/google/appengine/api/logservice/logservice_stub.py`.  
You can find your sdk root path by: `gcloud info --format="value(installation.sdk_root)`

#### diff
```diff
--- a/logservice_stub.py
+++ b/logservice_stub.py
@@ -86,7 +86,7 @@ class LogServiceStub(apiproxy_stub.APIProxyStub):
-  _MIN_COMMIT_INTERVAL = 5
+  _MIN_COMMIT_INTERVAL = 0
```
