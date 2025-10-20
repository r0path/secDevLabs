# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import uuid
import pickle
import base64
import time
app = Flask(__name__)

# Simple in-memory brute-force protection (best-effort for this sample app):
# - FAILED_LOGIN maps IP -> {count, first, locked_until}
# - MAX_ATTEMPTS within WINDOW_SECONDS triggers a temporary lockout of LOCKOUT_SECONDS
FAILED_LOGIN = {}
MAX_ATTEMPTS = 5
WINDOW_SECONDS = 300
LOCKOUT_SECONDS = 300


@app.route("/")
def ola():
    return render_template('index.html')

@app.route("/admin", methods=['GET','POST'])
def login():
    # Simple in-memory brute-force protection (per-IP)
    global FAILED_LOGIN, MAX_ATTEMPTS, WINDOW_SECONDS, LOCKOUT_SECONDS
    try:
        FAILED_LOGIN
    except NameError:
        FAILED_LOGIN = {}
        MAX_ATTEMPTS = 5
        WINDOW_SECONDS = 300
        LOCKOUT_SECONDS = 300

    if request.method == 'POST':
        import time
        ip = request.headers.get('X-Forwarded-For', request.remote_addr)
        if ip and ',' in ip:
            ip = ip.split(',')[0].strip()
        now = time.time()
        info = FAILED_LOGIN.get(ip, {"count":0, "first":None, "locked_until":0})
        if info.get('locked_until',0) > now:
            return redirect('/admin')
        if info.get('first') is None or (now - info.get('first', now)) > WINDOW_SECONDS:
            info = {"count":0, "first":now, "locked_until":0}

        username = request.values.get('username')
        password = request.values.get('password')

        if username == "admin" and password == "admin":
            # reset failed attempts on successful login
            if ip in FAILED_LOGIN:
                try:
                    del FAILED_LOGIN[ip]
                except Exception:
                    pass
            token = str(uuid.uuid4().hex)
            cookie = { "username":username, "admin":True, "sessionId":token }
            pickle_resultado = pickle.dumps(cookie)
            encodedSessionCookie = base64.b64encode(pickle_resultado)
            resp = make_response(redirect("/user"))
            resp.set_cookie("sessionId", encodedSessionCookie)
            return resp
        else:
            # record failed attempt
            info['count'] = info.get('count',0) + 1
            if info['count'] >= MAX_ATTEMPTS:
                info['locked_until'] = now + LOCKOUT_SECONDS
            FAILED_LOGIN[ip] = info
            return redirect('/admin')

    else:
        return render_template('admin.html')

@app.route("/user", methods=['GET'])
def userInfo():
    cookie = request.cookies.get("sessionId")
    if cookie == None:
        return "Não Autorizado!"
    cookie = pickle.loads(base64.b64decode(cookie))

    return render_template('user.html')
    



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
