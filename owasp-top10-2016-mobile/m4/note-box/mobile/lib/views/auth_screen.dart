import 'dart:convert';
import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart';

import './register_screen.dart';
import '../models/session_token.dart';
import '../widgets/alert_button.dart';
import '../widgets/user_notes.dart';

class AuthScreen extends StatefulWidget {
  @override
  _AuthScreenState createState() => _AuthScreenState();
}

class _AuthScreenState extends State<AuthScreen> {
  final usernameController = TextEditingController();
  final passwordController = TextEditingController();
  final _storage = new FlutterSecureStorage();

  void _login(BuildContext context, String username, String password) async {
    // set up POST request arguments
    String host = Platform.isAndroid ? "10.0.2.2" : "localhost";
    String url = 'http://$host:9051/login';
    Map<String, String> headers = {"Content-type": "application/json"};
    String json = '{"username": "$username", "password": "$password"}';
    // make POST request
    Response response = await post(url, headers: headers, body: json);
    // check the status code for the result
    int statusCode = response.statusCode;
    if (statusCode == 409) {
      showAlertDialog(context, 'Login Error', 'User already logged in!');
      return;
    }
    if (statusCode == 404) {
      // TO DO:
      // Handle wrong username or password so it can't be exploited.
      showAlertDialog(context, 'Login Error', 'Wrong username or password.');
      return;
    }

    // Only treat explicit success responses as success
    if (statusCode != 200) {
      showAlertDialog(context, 'Login Error', 'Unexpected response from server (status: $statusCode).');
      return;
    }

    Map sessionTokenMap;
    try {
      sessionTokenMap = jsonDecode(response.body);
    } catch (e) {
      showAlertDialog(context, 'Login Error', 'Invalid response from server.');
      return;
    }

    var sessionToken;
    try {
      sessionToken = SessionToken.fromJson(sessionTokenMap);
    } catch (e) {
      showAlertDialog(context, 'Login Error', 'Invalid session data received.');
      return;
    }

    await _storage.write(key: username, value: sessionToken.Value);

    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => UserNotes(username)),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Container(
        padding: EdgeInsets.all(30),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.center,
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Container(
              height: 55,
              width: 200,
              child: Card(
                elevation: 30,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: Colors.green[700],
                child: Text(
                  'Note Box',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    color: Colors.black87,
                    fontSize: 40,
                  ),
                ),
              ),
            ),
            Container(
              width: 280,
              child: TextField(
                decoration: InputDecoration(
                  labelText: 'Username',
                  contentPadding: EdgeInsets.all(10),
                  labelStyle: TextStyle(
                    color: Colors.black,
                  ),
                ),
                controller: usernameController,
              ),
            ),
            Container(
              width: 280,
              child: TextField(
                obscureText: true,
                decoration: InputDecoration(
                  labelText: 'Password',
                  contentPadding: EdgeInsets.all(10),
                  labelStyle: TextStyle(
                    color: Colors.black,
                  ),
                ),
                controller: passwordController,
              ),
            ),
            Container(
              padding: EdgeInsets.only(top: 15),
              child: FlatButton(
                onPressed: () {
                  _login(context, usernameController.text,
                      passwordController.text);
                },
                color: Colors.green[700],
                child: Text(
                  'Login',
                  style: TextStyle(
                    color: Colors.white,
                  ),
                ),
              ),
            ),
            Container(
              padding: EdgeInsets.only(top: 5),
              child: InkWell(
                child: Text(
                  'First time? Click here to register now!',
                  style: TextStyle(
                    color: Colors.green[700],
                    fontSize: 15,
                  ),
                ),
                onTap: () {
                  Navigator.push(
                    context,
                    MaterialPageRoute(builder: (context) => RegisterScreen()),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
