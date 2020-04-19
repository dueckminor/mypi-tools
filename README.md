# MyPi-Tools

![JS](https://github.com/dueckminor/mypi-tools/workflows/JS/badge.svg)
![Go](https://github.com/dueckminor/mypi-tools/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dueckminor/mypi-tools)](https://goreportcard.com/report/github.com/dueckminor/mypi-tools)

## Tools

### mypi-setup

Runs on your PC/Notebook and creates an SD-Card for a automated installation
of Alpine-Linux on a Raspberry-PI. It's not necessary to connect a keyboard
or display to your Raspberry-PI. All you need is a network cable and a power
supply.

When you start `mypi-setup` a browser window will open. Just follow the
instructions and you will get a SD-Card ready for your Raspberry-PI. After
you have inserted this SD-Card to your Raspberry-PI, `mypi-setup` will connect
to it and completes the installation.

### mypi-admin

The admin UI running on your Raspberry-PI.

### mypi-auth

The user authentication service. Provides a Login-UI.

### mypi-router

Allows to add authentication in front of apps.

### mypi-videostream

Allows to access web-cams

### mypi-owntracks

Reacts on owntracks events. It opens my gate when I come home.

## Used Ports

| Component   | Port | Webpack-Port |
|-------------|------|--------------|
| admin       | 9000 | 9001         |
| auth        | 9100 | 9101         |
| router      | 9200 | N/A          |
| videostream | 9300 | 9301         |
| owntracks   | 9400 | N/A          |
| setup       | 9500 | 9501         |
| control     | 9600 | 9601         |

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
