# gaelv
gaelv is a Log Viewer for Google App Engine local development.

## Install
```
go get -u github.com/addsict/gaelv
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

## TODO

* build with static files using go-bindata
* Fix /usr/local/google-cloud-sdk/platform/google_appengine/google/appengine/api/logservice/logservice_stub.py
