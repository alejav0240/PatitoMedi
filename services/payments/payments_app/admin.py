from django.contrib import admin

from .models import Invoice, Refund, Transaction


@admin.register(Invoice)
class InvoiceAdmin(admin.ModelAdmin):
    list_display = ("id", "appointment_id", "amount", "currency", "status", "created_at")
    search_fields = ("id", "appointment_id", "patient_id", "doctor_id")
    list_filter = ("status", "currency")


@admin.register(Transaction)
class TransactionAdmin(admin.ModelAdmin):
    list_display = ("id", "invoice", "amount", "currency", "status", "created_at")
    search_fields = ("id", "provider_reference", "invoice__appointment_id")
    list_filter = ("status", "currency")


@admin.register(Refund)
class RefundAdmin(admin.ModelAdmin):
    list_display = ("id", "transaction", "amount", "status", "created_at")
    search_fields = ("id", "provider_reference", "transaction__id")
    list_filter = ("status",)