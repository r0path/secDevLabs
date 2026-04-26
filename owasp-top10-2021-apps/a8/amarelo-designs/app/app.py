# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import uuid
import base64
from itsdangerous import URLSafeSerializer, BadSignature
app = Flask(__name__)
app.secret_key = app.config.get("SECRET_KEY", "change-me-in-production")
serializer = URLSafeSerializer(app.secret_key, salt="sessionId")


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
            cookie = { "username": username, "admin": True, "sessionId": token }
            encodedSessionCookie = serializer.dumps(cookie)
            resp = make_response(redirect("/user"))
            resp.set_cookie("sessionId", encodedSessionCookie, httponly=True, samesite="Lax")
            return resp

        else:
            return redirect("/admin")

    else:
        return render_template('admin.html')

@app.route("/user", methods=['GET'])
def userInfo():
    cookie = request.cookies.get("sessionId")
    if cookie is None:
        return "Não Autorizado!"
    try:
        cookie = serializer.loads(cookie)
    except BadSignature:
        return "Não Autorizado!"

    return render_template('user.html')
    



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
