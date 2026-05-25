import json
import time
from decimal import Decimal, InvalidOperation

from django.db import connections
from django.http import HttpResponse, JsonResponse
from django.views.decorators.http import require_http_methods

from .event_bus import publish_event
from .metrics import INVOICES_TOTAL, PROVIDER_LATENCY, REFUNDS_TOTAL, REQUESTS_TOTAL, TRANSACTIONS_TOTAL, metrics_bytes, metrics_content_type
from .models import Invoice, Refund, Transaction


def _response(payload, status=200):
    return JsonResponse(payload, status=status, json_dumps_params={"ensure_ascii": False})


def _record(endpoint, method, response):
    REQUESTS_TOTAL.labels(endpoint=endpoint, method=method, status=str(response.status_code)).inc()
    return response


def _json_body(request):
    if not request.body:
        return {}
    return json.loads(request.body.decode("utf-8"))


def _parse_amount(value):
    try:
        return Decimal(str(value))
    except (InvalidOperation, TypeError, ValueError):
        raise ValueError("amount must be a valid number")


@require_http_methods(["POST"])
def invoice_collection(request):
    endpoint = "/api/payments/invoices"
    try:
        data = _json_body(request)
        invoice = Invoice.objects.create(
            appointment_id=data["appointment_id"],
            patient_id=data.get("patient_id", ""),
            doctor_id=data.get("doctor_id", ""),
            amount=_parse_amount(data["amount"]),
            currency=data.get("currency", "USD"),
            metadata=data.get("metadata", {}),
        )
    except (KeyError, TypeError, ValueError, json.JSONDecodeError) as exc:
        return _record(endpoint, request.method, _response({"error": str(exc)}, status=400))

    INVOICES_TOTAL.inc()
    publish_event(
        "invoice-created",
        {"invoice_id": str(invoice.id), "appointment_id": invoice.appointment_id, "amount": str(invoice.amount), "currency": invoice.currency},
    )
    return _record(
        endpoint,
        request.method,
        _response(
            {"id": str(invoice.id), "appointment_id": invoice.appointment_id, "amount": str(invoice.amount), "currency": invoice.currency, "status": invoice.status},
            status=201,
        ),
    )


@require_http_methods(["GET"])
def invoice_detail(request, invoice_id):
    endpoint = "/api/payments/invoices/<id>"
    try:
        invoice = Invoice.objects.get(pk=invoice_id)
    except Invoice.DoesNotExist:
        return _record(endpoint, request.method, _response({"error": "invoice not found"}, status=404))

    return _record(endpoint, request.method, _response({"id": str(invoice.id), "appointment_id": invoice.appointment_id, "patient_id": invoice.patient_id, "doctor_id": invoice.doctor_id, "amount": str(invoice.amount), "currency": invoice.currency, "status": invoice.status, "external_reference": invoice.external_reference}))


@require_http_methods(["POST"])
def transaction_collection(request):
    endpoint = "/api/payments/transactions"
    try:
        data = _json_body(request)
        invoice = Invoice.objects.get(pk=data["invoice_id"])
        outcome = data.get("simulate_outcome", "approved")
        provider_reference = data.get("provider_reference", f"prov_{invoice.id.hex[:12]}")
        with PROVIDER_LATENCY.time():
            transaction = Transaction.objects.create(
                invoice=invoice,
                provider_reference=provider_reference,
                amount=invoice.amount,
                currency=invoice.currency,
                status=Transaction.Status.APPROVED if outcome == "approved" else Transaction.Status.FAILED,
                provider_payload=data.get("provider_payload", {}),
            )
        if transaction.status == Transaction.Status.APPROVED:
            Invoice.objects.filter(pk=invoice.pk).update(status=Invoice.Status.PAID, external_reference=provider_reference)
            publish_event("payment-confirmed", {"transaction_id": str(transaction.id), "invoice_id": str(invoice.id), "amount": str(transaction.amount), "currency": transaction.currency})
        else:
            publish_event("payment-failed", {"transaction_id": str(transaction.id), "invoice_id": str(invoice.id), "amount": str(transaction.amount), "currency": transaction.currency})
    except (KeyError, TypeError, ValueError, json.JSONDecodeError) as exc:
        return _record(endpoint, request.method, _response({"error": str(exc)}, status=400))
    except Invoice.DoesNotExist:
        return _record(endpoint, request.method, _response({"error": "invoice not found"}, status=404))

    TRANSACTIONS_TOTAL.labels(status=transaction.status).inc()
    return _record(endpoint, request.method, _response({"id": str(transaction.id), "invoice_id": str(transaction.invoice_id), "amount": str(transaction.amount), "currency": transaction.currency, "status": transaction.status, "provider_reference": transaction.provider_reference}, status=201))


@require_http_methods(["GET"])
def transaction_detail(request, transaction_id):
    endpoint = "/api/payments/transactions/<id>"
    try:
        transaction = Transaction.objects.get(pk=transaction_id)
    except Transaction.DoesNotExist:
        return _record(endpoint, request.method, _response({"error": "transaction not found"}, status=404))

    return _record(endpoint, request.method, _response({"id": str(transaction.id), "invoice_id": str(transaction.invoice_id), "amount": str(transaction.amount), "currency": transaction.currency, "status": transaction.status, "provider_reference": transaction.provider_reference}))


@require_http_methods(["POST"])
def refund_collection(request):
    endpoint = "/api/payments/refunds"
    try:
        data = _json_body(request)
        transaction = Transaction.objects.get(pk=data["transaction_id"])
        amount = _parse_amount(data.get("amount", transaction.amount))
        refund = Refund.objects.create(
            transaction=transaction,
            amount=amount,
            reason=data.get("reason", ""),
            status=Refund.Status.PROCESSED if transaction.status == Transaction.Status.APPROVED else Refund.Status.FAILED,
            provider_reference=data.get("provider_reference", f"refund_{transaction.id.hex[:12]}"),
            provider_payload=data.get("provider_payload", {}),
        )
    except (KeyError, TypeError, ValueError, json.JSONDecodeError) as exc:
        return _record(endpoint, request.method, _response({"error": str(exc)}, status=400))
    except Transaction.DoesNotExist:
        return _record(endpoint, request.method, _response({"error": "transaction not found"}, status=404))

    if refund.status == Refund.Status.PROCESSED:
        Transaction.objects.filter(pk=transaction.pk).update(status=Transaction.Status.REFUNDED)
    REFUNDS_TOTAL.inc()
    publish_event("refund-created", {"refund_id": str(refund.id), "transaction_id": str(transaction.id), "amount": str(refund.amount), "status": refund.status})
    return _record(endpoint, request.method, _response({"id": str(refund.id), "transaction_id": str(transaction.id), "amount": str(refund.amount), "status": refund.status, "provider_reference": refund.provider_reference}, status=201))


@require_http_methods(["POST"])
def provider_webhook(request):
    endpoint = "/api/payments/webhooks/provider"
    try:
        data = _json_body(request)
    except json.JSONDecodeError as exc:
        return _record(endpoint, request.method, _response({"error": str(exc)}, status=400))

    event_type = data.get("event_type")
    payload = data.get("payload", {})
    started = time.monotonic()

    if event_type in {"payment.succeeded", "payment.confirmed"}:
        transaction_id = payload.get("transaction_id")
        if transaction_id:
            Transaction.objects.filter(pk=transaction_id).update(status=Transaction.Status.APPROVED, provider_payload=payload)
            transaction = Transaction.objects.filter(pk=transaction_id).first()
            if transaction:
                Invoice.objects.filter(pk=transaction.invoice_id).update(status=Invoice.Status.PAID)
        publish_event("payment-confirmed", payload)
    elif event_type == "payment.failed":
        transaction_id = payload.get("transaction_id")
        if transaction_id:
            Transaction.objects.filter(pk=transaction_id).update(status=Transaction.Status.FAILED, provider_payload=payload)
        publish_event("payment-failed", payload)
    elif event_type == "refund.succeeded":
        refund_id = payload.get("refund_id")
        if refund_id:
            Refund.objects.filter(pk=refund_id).update(status=Refund.Status.PROCESSED, provider_payload=payload)
        publish_event("refund-created", payload)

    PROVIDER_LATENCY.observe(time.monotonic() - started)
    return _record(endpoint, request.method, _response({"accepted": True, "event_type": event_type}, status=202))


def health_view(request):
    endpoint = "/health"
    try:
        with connections["default"].cursor() as cursor:
            cursor.execute("SELECT 1")
            cursor.fetchone()
    except Exception as exc:
        return _record(endpoint, request.method, _response({"status": "degraded", "error": str(exc)}, status=503))
    return _record(endpoint, request.method, _response({"status": "ok", "database": "up"}))


def metrics_view(request):
    endpoint = "/metrics"
    response = HttpResponse(metrics_bytes(), content_type=metrics_content_type())
    return _record(endpoint, request.method, response)