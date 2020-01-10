# MyPi-API

## Used Ports

| Component   | Port | Webpack-Port |
|-------------|------|--------------|
| admin       | 9000 | 9001         |
| auth        | 9100 | 9101         |
| router      | 9200 | N/A          |
| videostream | 9300 | 9301         |
| owntracks   | 9400 | N/A          |
| setup       | 9500 | 9501         |

## Dev-Environment

To allow debugging of the GoLang and Web Applications you could use do the
following:

1. mount your `mypi` installation of your Raspberry-PI in the same directory
as it is remote:

```bash
sshfs pi@rpi.fritz.box:/opt/mypi /opt/mypi
```

2. call the `prepare_debug` script:

```bash
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
+-------------+  8080   +----+--------+          |
| MyPI-Admin  |----<----|  SSH-Tunnel |-----<----+
| (GO)        |---->----|             |-->--+
+-------------+  11111  +-------------+     |
      |                      |              |
      | 9100                 |              |
      |                      |  +--------------+
+------------|               |  | Docker-API   |
| MyPI-Admin |               |  +--------------+
| (JS)       |               |
+------------+               |
```
