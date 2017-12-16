# SYNOPSIS
```sh
# app server
dev_appserver.py app.yaml --logs_path=/tmp/gaelog.db

# log viewer
gae_log_viewer --logs_path=/tmp/gaelog.db
```

And open browser: http://localhost:92384/

```sh
# to console
gae_log_viewer --logs_path=/tmp/gaelog.db --no-ui
```

# Web UI

* Send log line to web browser through Server Sent Events
