# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import uuid
import base64
import json
from flask import session
app = Flask(__name__)
app.secret_key = uuid.uuid4().hex


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
            session["username"] = username
            session["admin"] = True
            session["sessionId"] = token
            return redirect("/user")

        else:
            return redirect("/admin")

    else:
        return render_template('admin.html')

@app.route("/user", methods=['GET'])
def userInfo():
    if not session.get("sessionId"):
        return "Não Autorizado!"

    return render_template('user.html')
    



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
