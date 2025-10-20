# coding: utf-8

from flask import Flask, request, make_response, render_template, redirect, flash
import uuid
import pickle
import base64
app = Flask(__name__)


# Set secure response headers
@app.after_request
def set_security_headers(response):
    """Apply common security headers to all responses.
    - Content-Security-Policy (CSP): conservative default limiting sources to self. Adjust if your app loads resources from CDNs.
    - X-Content-Type-Options: nosniff to prevent MIME sniffing.
    - X-Frame-Options: SAMEORIGIN to mitigate clickjacking.
    - Strict-Transport-Security (HSTS): applied only when not in debug and request is secure (or X-Forwarded-Proto indicates https).
    """
    # Content Security Policy: adjust as needed for your application's external resources
    response.headers.setdefault("Content-Security-Policy", "default-src 'self';")
    # Prevent MIME type sniffing
    response.headers.setdefault("X-Content-Type-Options", "nosniff")
    # Prevent clickjacking
    response.headers.setdefault("X-Frame-Options", "SAMEORIGIN")
    # Set HSTS only when not in debug mode and over HTTPS
    try:
        if not app.debug and (request.is_secure or request.headers.get('X-Forwarded-Proto', '') == 'https'):
            response.headers.setdefault("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
    except Exception:
        # If request context not available, skip HSTS
        pass
    return response


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
            pickle_resultado = pickle.dumps(cookie)
            encodedSessionCookie = base64.b64encode(pickle_resultado)
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
    cookie = pickle.loads(base64.b64decode(cookie))

    return render_template('user.html')
    



if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0')
