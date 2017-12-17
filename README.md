# gaelv
gaelv is a Log Viewer for Google App Engine local development.

## SYNOPSIS
```sh
# run your app server
dev_appserver.py app.yaml --logs_path=/tmp/gaelog.db

# run log viewer
gaelv --logs_path=/tmp/gaelog.db
```

And open browser: http://localhost:9090/
