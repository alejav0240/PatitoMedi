from django.urls import include, path

from payments_app import views as payments_views

urlpatterns = [
    path("", include("payments_app.urls")),
    path("health", payments_views.health_view),
    path("metrics", payments_views.metrics_view),
]
