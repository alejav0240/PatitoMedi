from django.apps import AppConfig


class PaymentsAppConfig(AppConfig):
    default_auto_field = "django.db.models.BigAutoField"
    name = "payments_app"
    verbose_name = "Payments"

    def ready(self):
        from .consumers import start_background_consumer

        start_background_consumer()