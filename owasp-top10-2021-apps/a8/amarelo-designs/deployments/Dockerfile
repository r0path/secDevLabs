FROM python:3
WORKDIR /app
ADD app/requirements.txt /app/requirements.txt
# Consertando a build, instalando o pacote netcat-traditional ao invés de `netcat`
RUN apt-get update && apt-get -y install netcat-traditional && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN pip install -r requirements.txt
# Create an unprivileged user and give ownership of /app to avoid running as root
RUN groupadd -r appuser && useradd -r -g appuser -d /home/appuser -s /usr/sbin/nologin appuser \
    && mkdir -p /home/appuser \
    && chown -R appuser:appuser /app /home/appuser
USER appuser
CMD ["python", "app.py"]