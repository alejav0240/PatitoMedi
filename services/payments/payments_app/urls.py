from django.urls import path

from . import views

urlpatterns = [
    path("invoices", views.invoice_collection),
    path("invoices/<uuid:invoice_id>", views.invoice_detail),
    path("transactions", views.transaction_collection),
    path("transactions/<uuid:transaction_id>", views.transaction_detail),
    path("refunds", views.refund_collection),
    path("webhooks/provider", views.provider_webhook),
]