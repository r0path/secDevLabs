#!/usr/bin/env python
# -*- coding: utf-8 -*-
from flask import Flask, request, make_response, render_template, redirect, Markup
from model.password import Password
from model.db import DataBase
import base64
import os
import json
import hashlib
import uuid
import time
from functools import wraps


app = Flask(__name__)
database = DataBase(os.environ.get('A2_DATABASE_HOST'),
                    os.environ.get('A2_DATABASE_USER'),
                    os.environ.get('A2_DATABASE_PASSWORD'),
                    os.environ.get('A2_DATABASE_NAME'))

# Simple in-memory brute-force protections.
# Note: This uses in-memory storage and will not persist across processes.
# For production environments behind multiple workers or hosts, replace with
# a centralized store (e.g. Redis) and integrate with a proper rate-limiting
# or account lockout mechanism.
MAX_ATTEMPTS = 5
WINDOW_SECONDS = 300  # sliding window for counting attempts (5 minutes)
LOCKOUT_SECONDS = 300  # lockout duration after MAX_ATTEMPTS exceeded (5 minutes)

failed_attempts_ip = {}       # ip -> [timestamps]
failed_attempts_user = {}     # username -> [timestamps]
locked_until_user = {}        # username -> timestamp when lockout ends



def login_admin_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        cookie = request.cookies.get("sessionId", "")
        cookie = base64.b64decode(cookie).decode("utf-8")
        cookie_separado = cookie.split('.')
        if(len(cookie_separado) != 2):
            return "Invalid cookie!"
        hash_cookie = hashlib.sha256(cookie_separado[0].encode('utf-8')).hexdigest()
        if (hash_cookie != cookie_separado[1]):
            return redirect("/login")
        j = json.loads(cookie_separado[0])
        if j.get("permissao") != 1:
            return "You don't have permission to access this route. You are not an admin. \n"
        return f(*args, **kwargs)
    return decorated_function


def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        cookie = request.cookies.get("sessionId", "")
        cookie = base64.b64decode(cookie).decode("utf-8")
        cookie_separado = cookie.split('.')
        if(len(cookie_separado) != 2):
            return "Invalid cookie! \n"
        hash_cookie = hashlib.sha256(cookie_separado[0].encode('utf-8')).hexdigest()
        if (hash_cookie != cookie_separado[1]):
            return redirect("/login")
        return f(*args, **kwargs)
    return decorated_function


@app.route("/", methods=['GET'])
def home():
    return render_template('index.html')


@app.route("/register", methods=['GET', 'POST'])
def register():
    if request.method == 'GET':
        return render_template('register.html')

    if request.method == 'POST':
        form_username = request.form.get('username', "")
        form_password = request.form.get('password', "")
        form_password2 = request.form.get('password2', "")

        if form_username == "" or form_password == "":
            return "Error! You have to pass username and password! \n"
        elif form_password != form_password2:
            return "Error! Passwords must be the same! \n"

        guid = str(uuid.uuid4())
        password = Password(form_password, form_username, guid)
        hashed_password = password.get_hashed_password()
        message, success = database.insert_user(guid, form_username, hashed_password)
        if success:
            return render_template('login.html')
        return "Error: account creation failed \n"


@app.route("/login", methods=['GET', 'POST'])
def login():
    if request.method == 'GET':
        return render_template('login.html')

    if request.method == 'POST':
        form_username = request.form.get('username', "")
        form_password = request.form.get('password', "")
        if form_username == "" or form_password == "":
            return "Error! You have to pass username and password! \n"

        # Client IP (respect X-Forwarded-For when behind a proxy)
        client_ip = request.headers.get('X-Forwarded-For', request.remote_addr)
        if client_ip:
            client_ip = client_ip.split(',')[0].strip()

        now = time.time()

        def _prune(timestamps):
            return [ts for ts in timestamps if now - ts <= WINDOW_SECONDS]

        # Prune old IP attempts and enforce IP rate limit
        ip_attempts = failed_attempts_ip.get(client_ip, [])
        ip_attempts = _prune(ip_attempts)
        failed_attempts_ip[client_ip] = ip_attempts
        if len(ip_attempts) >= MAX_ATTEMPTS:
            return make_response("Too many requests from your IP. Try again later.\n", 429)

        # Check if account is locked
        locked_until = locked_until_user.get(form_username)
        if locked_until and now < locked_until:
            return make_response("Too many failed attempts. Account temporarily locked. Try again later.\n", 429)

        result, success = database.get_user(form_username)
        if not success or result is None:
            # Record failed attempt
            failed_attempts_ip.setdefault(client_ip, []).append(now)
            failed_attempts_user.setdefault(form_username, []).append(now)
            # Prune and check user attempts
            user_attempts = _prune(failed_attempts_user.get(form_username, []))
            failed_attempts_user[form_username] = user_attempts
            if len(user_attempts) >= MAX_ATTEMPTS:
                locked_until_user[form_username] = now + LOCKOUT_SECONDS
            return "Login failed! \n"

        password = Password(form_password, form_username, result[2])
        if not password.validate_password(result[0]):
            # Record failed attempt
            failed_attempts_ip.setdefault(client_ip, []).append(now)
            failed_attempts_user.setdefault(form_username, []).append(now)
            user_attempts = _prune(failed_attempts_user.get(form_username, []))
            failed_attempts_user[form_username] = user_attempts
            if len(user_attempts) >= MAX_ATTEMPTS:
                locked_until_user[form_username] = now + LOCKOUT_SECONDS
            return "Login failed! \n"

        # Successful login -> clear failed attempts and lockout
        failed_attempts_user.pop(form_username, None)
        failed_attempts_ip.pop(client_ip, None)
        locked_until_user.pop(form_username, None)

        cookie_dic = {"permissao": result[1], "username": form_username}
        cookie = json.dumps(cookie_dic)
        hash_cookie = hashlib.sha256(cookie.encode('utf-8')).hexdigest()
        cookie_done = '.'.join([cookie,hash_cookie])
        cookie_done = base64.b64encode(str(cookie_done).encode("utf-8"))
        resp = make_response("Logged in!")
        resp.set_cookie("sessionId", cookie_done)
        return resp


@app.route("/admin", methods=['GET'])
@login_admin_required
def admin():
    return "You are an admin! \n"


@app.route("/user", methods=['GET'])
@login_required
def userInfo():
    return "You are an user! \n"


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=10002)
