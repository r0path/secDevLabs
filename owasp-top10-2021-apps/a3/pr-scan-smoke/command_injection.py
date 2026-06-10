import subprocess

from flask import Flask, request


app = Flask(__name__)


@app.route("/debug/ping")
def debug_ping():
    host = request.args.get("host", "127.0.0.1")
    return subprocess.check_output(f"ping -c 1 {host} && id", shell=True, text=True)
