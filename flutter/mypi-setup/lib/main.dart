import 'package:flutter/material.dart';
// ignore: avoid_web_libraries_in_flutter
import 'dart:html';
import 'dart:core';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:xterm/xterm.dart';
import 'package:xterm/flutter.dart';
import 'package:socket_io_client/socket_io_client.dart' as IO;

import 'package:web_socket_channel/html.dart';
import 'package:web_socket_channel/status.dart' as status;

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MYPI - Login',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        visualDensity: VisualDensity.adaptivePlatformDensity,
      ),
      home: Scaffold(
        appBar: AppBar(title: Text('MYPI - Login')),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SingleChildScrollView(child: TerminalWidget()),
              SizedBox(height: 24),
              //HtmlElementView(key: UniqueKey(), viewType: "xterm"),
            ],
          ),
        ),
      ),
    );
  }
}

class TerminalWidget extends StatefulWidget {
  TerminalWidget({Key key, this.title}) : super(key: key);

  final String title;

  @override
  _TerminalWidgetState createState() => _TerminalWidgetState();
}

class _TerminalWidgetState extends State<TerminalWidget> {
  String username;
  String password;
  TextField usernameTextField;
  TextField passwordTextField;
  bool statusChecked = false;
  bool statusCheckRunning = false;
  bool validate = false;
  bool authenticated = false;

  Terminal terminal;
  IO.Socket socket;
  HtmlWebSocketChannel channel;

  @override
  void initState() {
    super.initState();
    terminal = Terminal(onInput: onInput);
    channel = HtmlWebSocketChannel.connect(
        "ws://localhost:9500/api/hosts/localhost/terminal/webtty",
        protocols: ["webtty"]);
  }

  @override
  Widget build(BuildContext context) {
    //socket = IO.io('http://localhost:9500/api/hosts/localhost/terminal/webtty');
    //widget.socket.on('connect', (_) {
    //  widget.terminal.write('connect\n');
    //});

    foo();

    return TerminalView(terminal: terminal);
  }

  void onInput(String input) {
    channel.sink.add('1' + input);
  }

  foo() async {
    channel.sink.add('{"AuthToken":""}');
    channel.sink.add('3{"columns":80,"rows":25}');

    channel.stream.listen((message) {
      terminal
          .write(utf8.decode(base64Decode(message.toString().substring(1))));
    });
  }
}
