#!/usr/bin/env python
# -*- coding: utf-8 -*-
from functools import wraps
import uuid
import datetime
from flask import (
    Flask,
    render_template,
    request,
    redirect,
    flash,
    make_response,
    session
)
from util.init_db import init_db
from flask.logging import default_handler
from flask_bootstrap import Bootstrap
from model.password import Password
from model.db import DataBase
import logging
import os
import time
from threading import Lock

from flask_cors import CORS, cross_origin
from model.db import DataBase

app = Flask(__name__)
bootstrap = Bootstrap(app)

# Simple in-memory rate limiter for login to mitigate brute force attacks.
# NOTE: This is a per-process in-memory limiter and is suitable for single-process
# deployments or testing. For production use across multiple processes/hosts,
# replace with a distributed store (e.g., Redis) and a robust rate-limiting
# middleware.
RATE_LIMIT_MAX_ATTEMPTS = 5
RATE_LIMIT_WINDOW = 300  # seconds
LOCKOUT_TIME = 600  # seconds
_login_attempts = {}
_login_lock = Lock()

app.config.from_pyfile('config.py')


def generate_csrf_token():
    '''
        Generate csrf token and store it in session
    '''
    if '_csrf_token' not in session:
        session['_csrf_token'] = str(uuid.uuid4())
    return session.get('_csrf_token')


app.jinja_env.globals['csrf_token'] = generate_csrf_token

@app.before_request
def csrf_protect():
    '''
        CSRF PROTECION
    '''
    if request.method == "POST":
        token_csrf = session.get('_csrf_token')
        form_token = request.form.get('_csrf_token')
        if not token_csrf or str(token_csrf) != str(form_token):
            return "ERROR: Wrong value for csrf_token"

def login_required(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        if 'username' not in session:
            flash('oops, session expired', "danger")
            return redirect('/login')
        return f(*args, **kwargs)
    return decorated_function

@app.route('/', methods=['GET'])
def root():
    return redirect('/login')

@app.route('/logout', methods=['GET'])
@login_required
def logout():
    session.clear()
    return redirect('/login')

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form.get('username').encode('utf-8')
        username_str = username.decode('utf-8', errors='ignore') if username else ''
        client_ip = request.remote_addr or 'unknown'
        key = f"{username_str}:{client_ip}"
        now = time.time()
        # Check rate limit / lockout
        with _login_lock:
            entry = _login_attempts.get(key)
            if entry:
                if entry.get('locked_until', 0) > now:
                    # Locked out - generic error
                    flash("Usuario ou senha incorretos", "danger")
                    return render_template('login.html')
                # Reset if window expired
                if now - entry.get('first_seen', now) > RATE_LIMIT_WINDOW:
                    entry = {'attempts': 0, 'first_seen': now, 'locked_until': 0}
                    _login_attempts[key] = entry
            else:
                _login_attempts[key] = {'attempts': 0, 'first_seen': now, 'locked_until': 0}

        psw = Password(request.form.get('password').encode('utf-8'))
        user_password, success = database.get_user_password(username)
        valid = success and user_password != None and psw.validate_password(str(user_password[0]))
        if not valid:
            with _login_lock:
                e = _login_attempts.get(key)
                if e is None:
                    e = {'attempts': 1, 'first_seen': now, 'locked_until': 0}
                    _login_attempts[key] = e
                else:
                    # Reset attempts if window expired
                    if now - e.get('first_seen', now) > RATE_LIMIT_WINDOW:
                        e['attempts'] = 1
                        e['first_seen'] = now
                        e['locked_until'] = 0
                    else:
                        e['attempts'] = e.get('attempts', 0) + 1
                        if e['attempts'] >= RATE_LIMIT_MAX_ATTEMPTS:
                            e['locked_until'] = now + LOCKOUT_TIME
            # Generic error message to avoid username enumeration
            flash("Usuario ou senha incorretos", "danger")
            return render_template('login.html')

        # Successful login: clear any recorded attempts
        with _login_lock:
            if key in _login_attempts:
                try:
                    del _login_attempts[key]
                except Exception:
                    pass

        session['username'] = username
        return redirect('/home')
    else:
        return render_template('login.html')

@app.route('/register', methods=['GET', 'POST'])
def newuser():
    if request.method == 'POST':
        username = request.form.get('username').encode('utf-8')
        psw1 = request.form.get('password1').encode('utf-8')
        psw2 = request.form.get('password2').encode('utf-8')

        if psw1 == psw2:
            psw = Password(psw1)
            hashed_psw = psw.get_hashed_password()
            message, success = database.insert_user(username, hashed_psw)
            if success == 1:
                flash("Novo usuario adicionado!", "primary")
                return redirect('/login')
            else:
                flash(message, "danger")
                return redirect('/register')

        flash("Passwords must be the same!", "danger")
        return redirect('/register')
    else:
        return render_template('register.html')

@app.route('/home', methods=['GET'])
@login_required
def home():
    return render_template('index.html')

@app.route('/coupon', methods=['GET', 'POST'])
@login_required
def cupom():
    if request.method == 'POST':
        coupon = request.form.get('coupon')
        rows, success = database.get_game_coupon(coupon, session.get('username'))
        if not success or rows == None or rows == 0:
            flash("Cupom invalido", "danger")
            return render_template('coupon.html')
        game, success = database.get_game(coupon, session.get('username'))
        if not success or game == None:
            flash("Cupom invalido", "danger")
            return render_template('coupon.html')
        flash("Voce ganhou {}".format(game[0]), "primary")
        return render_template('coupon.html')
    else:
        return render_template('coupon.html')

if __name__ == '__main__':
    dbEndpoint = os.environ.get('MYSQL_ENDPOINT')
    dbUser = os.environ.get('MYSQL_USER')
    dbPassword = os.environ.get('MYSQL_PASSWORD')
    dbName = os.environ.get('MYSQL_DB')
    database = DataBase(dbEndpoint, dbUser, dbPassword, dbName)
    init_db(database)
    app.run(host='0.0.0.0',port=10010, debug=True)
