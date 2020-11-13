import 'package:flutter/material.dart';
// ignore: avoid_web_libraries_in_flutter
import 'dart:html';
import 'dart:core';
import 'package:http/http.dart' as http;
import 'dart:convert';

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
        body: Form(
          child: Scrollbar(
            child: SingleChildScrollView(
              padding: EdgeInsets.all(20),
              child: Center(
                child: Container(
                  width: 500,
                  child: LoginWidget(queryParameters: Uri.base.queryParameters),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class LoginWidget extends StatefulWidget {
  LoginWidget({Key key, this.title, this.queryParameters}) : super(key: key);

  final String title;
  final Map<String, String> queryParameters;

  @override
  _LoginWidgetState createState() => _LoginWidgetState();
}

class _LoginWidgetState extends State<LoginWidget> {
  String username;
  String password;
  TextField usernameTextField;
  TextField passwordTextField;
  bool statusChecked = false;
  bool statusCheckRunning = false;
  bool validate = false;
  bool authenticated = false;

  @override
  Widget build(BuildContext context) {
    if (!statusChecked) {
      if (!statusCheckRunning) {
        statusCheckRunning = true;
        checkStatus();
      }
      return Column(
        children: [Text("checking login status...")],
      );
    }

    if (authenticated) {
      return Column(
        children: [
          Text("your are logged in"),
          SizedBox(height: 24),
          Text(username),
          SizedBox(height: 24),
          ElevatedButton(
            child: Text('Logout'),
            onPressed: logout,
          ),
        ],
      );
    }

    usernameTextField = TextField(
      decoration: InputDecoration(
        filled: true,
        labelText: 'Username',
      ),
      onChanged: (value) {
        username = value;
      },
      enabled: !validate,
    );

    passwordTextField = TextField(
      decoration: InputDecoration(
        filled: true,
        labelText: 'Password',
      ),
      obscureText: true,
      onChanged: (value) {
        password = value;
      },
      enabled: !validate,
    );

    return Column(
      children: [
        usernameTextField,
        SizedBox(height: 24),
        passwordTextField,
        SizedBox(height: 24),
        ElevatedButton(
          child: Text('Sign in'),
          onPressed: validate ? null : login,
        ),
      ],
    );
  }

  login() async {
    setState(() {
      validate = true;
    });

    var response = await backend_post("/login",
        headers: {"Content-Type": "application/json"},
        body: jsonEncode({"username": username, "password": password}));
    if (response.statusCode == 200) {
      setState(() {
        authenticated = true;
        validate = false;
      });
      if (this.widget.queryParameters.containsKey("redirect_uri")) {
        window.location.assign(Uri(
                path: String.fromEnvironment("backend") + "/oauth/authorize",
                queryParameters: this.widget.queryParameters)
            .toString());
      }
    } else {
      setState(() {
        authenticated = false;
        validate = false;
      });
    }
  }

  logout() async {
    var response = await backend_post("/logout");
    if (response.statusCode == 202) {
      setState(() {
        authenticated = false;
        validate = false;
      });
    }
  }

  checkStatus() async {
    var response = await backend_get("/status");

    var data = json.decode(response.body);

    if ((data is Map) &&
        data.containsKey('username') &&
        data['username'] != "") {
      setState(() {
        username = data['username'];
        statusCheckRunning = false;
        statusChecked = true;
        authenticated = true;
      });
    } else {
      setState(() {
        statusCheckRunning = false;
        statusChecked = true;
        authenticated = false;
      });
    }
  }

  String backend_prefix() {
    const backend = String.fromEnvironment("backend", defaultValue: "");
    return backend;
  }

  Future<http.Response> backend_get(String path,
      {Map<String, String> headers}) async {
    return http.get(backend_prefix() + path, headers: headers);
  }

  Future<http.Response> backend_post(String path,
      {Map<String, String> headers, body, Encoding encoding}) async {
    return http.post(backend_prefix() + path,
        headers: headers, body: body, encoding: encoding);
  }
}
