# MyPi-API

## Used Ports

| Component   | Port | Webpack-Port |
|-------------|------|--------------|
| admin       | 9000 | 9001         |
| auth        | 9100 | 9101         |
| router      | 9200 | N/A          |
| videostream | 9300 | 9301         |
| owntracks   | 9400 | N/A          |

## Dev-Environment

```bash
sshfs pi@rpi.fritz.box:/opt/mypi /opt/mypi
./scripts/prepare_debug
```

```txt
Developer-PC                 | Raspberry-PI
-----------------------------+----------------------
                             |
+--------------+             | 443  +-------+
| Browser      |-------------|--->--| ngnix |----+
+--------------+             |      +-------+    |
                             |                   |
+-----------+  8080   +------+------+            |
| MyPI-API  |----<----|  SSH-Tunnel |-----<------+
| (GO)      |---->----|             |-->--+
+-----------+  11111  +-------------+     |
      |                      |            |
      | 9100                 |            |
      |                      |  +--------------+
+------------|               |  | Docker-API   |
| MyPI-Admin |               |  +--------------+
| (JS)       |               |
+------------+               |
```
