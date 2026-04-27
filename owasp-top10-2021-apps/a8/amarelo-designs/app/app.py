# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import uuid
import json
import base64
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
    
        admin_user = os.environ.get("ADMIN_USERNAME", "")
        admin_pass = os.environ.get("ADMIN_PASSWORD", "")
        if not admin_user or not admin_pass:
            return redirect("/admin")
        if username == admin_user and password == admin_pass:
            token = str(uuid.uuid4().hex)
            cookie = { "username":username, "admin":True, "sessionId":token }
            encodedSessionCookie = base64.b64encode(json.dumps(cookie).encode('utf-8'))
            resp = make_response(redirect("/user"))
            resp.set_cookie("sessionId", encodedSessionCookie)
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

    return render_template('user.html')
    



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
