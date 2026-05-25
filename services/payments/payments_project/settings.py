import os
from pathlib import Path

BASE_DIR = Path(__file__).resolve().parent.parent

SECRET_KEY = os.environ.get("DJANGO_SECRET_KEY", "dev-secret-key")
DEBUG = os.environ.get("DJANGO_DEBUG", "1") == "1"

ALLOWED_HOSTS = ["*"]
APPEND_SLASH = False

INSTALLED_APPS = [
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.staticfiles",
    "payments_app.apps.PaymentsAppConfig",
]

MIDDLEWARE = [
    "django.middleware.common.CommonMiddleware",
]

ROOT_URLCONF = "payments_project.urls"

TEMPLATES = []

WSGI_APPLICATION = "payments_project.wsgi.application"

DATABASES = {
    "default": {
        "ENGINE": os.environ.get("PAYMENTS_DB_ENGINE", "django.db.backends.postgresql"),
        "NAME": os.environ.get("PAYMENTS_DB_NAME", "payments"),
        "USER": os.environ.get("PAYMENTS_DB_USER", "payments_app"),
        "PASSWORD": os.environ.get("PAYMENTS_DB_PASSWORD", "payments_pass"),
        "HOST": os.environ.get("PAYMENTS_DB_HOST", "payments-db"),
        "PORT": os.environ.get("PAYMENTS_DB_PORT", "5432"),
    }
}

STATIC_URL = "/static/"
DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

KAFKA_BROKERS = os.environ.get("KAFKA_BROKERS", "kafka:9092")
KAFKA_ENABLED = os.environ.get("KAFKA_ENABLED", "true").lower() in {"1", "true", "yes", "on"}
PAYMENTS_CONSUMER_ENABLED = os.environ.get("PAYMENTS_CONSUMER_ENABLED", "true").lower() in {"1", "true", "yes", "on"}
