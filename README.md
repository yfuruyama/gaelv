gaemonitor or gae-log-viewer or gaelogviewr
===

SYNOPSIS
---
```sh
# app server
dev_appserver.py app.yaml --logs_path=/tmp/gaelog.db

# log viewer
gaemonitor --logs_path=/tmp/gaelog.db
```

And open browser: http://localhost:92384/

```sh
# to console
gaemonitor --logs_path=/tmp/gaelog.db --no-ur
```

Web UI
---

* Send log line to web browser through Server Sent Events
