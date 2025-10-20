# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import uuid
import base64
import json
import hmac
import hashlib
import os
app = Flask(__name__)


@app.route("/")
def ola():
    return render_template('index.html')

@app.route("/admin", methods=['GET','POST'])
def login():
    if request.method == 'POST':
        username = request.values.get('username')
        password = request.values.get('password')
    
        if username == "admin" and password == "admin":
            token = str(uuid.uuid4().hex)
            cookie = { "username":username, "admin":True, "sessionId":token }
            # Use a signed JSON cookie instead of insecure pickled data
            secret = os.environ.get('SESSION_SIGNING_KEY', 'default_change_me')
            json_cookie = json.dumps(cookie, separators=(',',':')).encode('utf-8')
            signature = hmac.new(secret.encode('utf-8'), json_cookie, hashlib.sha256).hexdigest()
            encodedSessionCookie = base64.b64encode(json_cookie).decode('utf-8') + '.' + signature
            resp = make_response(redirect("/user"))
            resp.set_cookie("sessionId", encodedSessionCookie, httponly=True, samesite='Lax')
            return resp

        else:
            return redirect("/admin")

    else:
        return render_template('admin.html')

@app.route("/user", methods=['GET'])
def userInfo():
    cookie = request.cookies.get("sessionId")
    if cookie == None:
        return "Não Autorizado!"
    try:
        # Expect format: base64(json).signature
        if '.' not in cookie:
            return "Não Autorizado!"
        encoded_json, signature = cookie.rsplit('.', 1)
        json_bytes = base64.b64decode(encoded_json)
        secret = os.environ.get('SESSION_SIGNING_KEY', 'default_change_me')
        expected_sig = hmac.new(secret.encode('utf-8'), json_bytes, hashlib.sha256).hexdigest()
        if not hmac.compare_digest(expected_sig, signature):
            return "Não Autorizado!"
        cookie = json.loads(json_bytes.decode('utf-8'))
    except Exception:
        return "Não Autorizado!"

    return render_template('user.html')
    



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
