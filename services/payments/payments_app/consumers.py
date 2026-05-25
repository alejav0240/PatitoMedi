import json
import logging
import os
import threading
from decimal import Decimal, InvalidOperation

from django.conf import settings

from .event_bus import publish_event
from .models import Invoice

logger = logging.getLogger(__name__)

_consumer_thread = None
_consumer_lock = threading.Lock()


def _parse_amount(value):
    try:
        return Decimal(str(value))
    except (InvalidOperation, TypeError, ValueError):
        return Decimal("0")


def handle_appointment_created(payload):
    appointment_id = str(payload.get("appointment_id", "")).strip()
    if not appointment_id:
        raise ValueError("appointment_id is required")

    invoice, created = Invoice.objects.get_or_create(
        appointment_id=appointment_id,
        defaults={
            "patient_id": str(payload.get("patient_id", "")),
            "doctor_id": str(payload.get("doctor_id", "")),
            "amount": _parse_amount(payload.get("amount", "0")),
            "currency": str(payload.get("currency", "USD")),
            "metadata": payload,
        },
    )

    if not created:
        invoice.patient_id = str(payload.get("patient_id", invoice.patient_id))
        invoice.doctor_id = str(payload.get("doctor_id", invoice.doctor_id))
        invoice.amount = _parse_amount(payload.get("amount", invoice.amount))
        invoice.currency = str(payload.get("currency", invoice.currency))
        invoice.metadata = payload
        if invoice.status == Invoice.Status.CANCELLED:
            invoice.status = Invoice.Status.ISSUED
        invoice.save()

    publish_event(
        "invoice-created",
        {
            "invoice_id": str(invoice.id),
            "appointment_id": invoice.appointment_id,
            "amount": str(invoice.amount),
            "currency": invoice.currency,
        },
    )
    return invoice


def handle_appointment_cancelled(payload):
    appointment_id = str(payload.get("appointment_id", "")).strip()
    if not appointment_id:
        raise ValueError("appointment_id is required")

    try:
        invoice = Invoice.objects.get(appointment_id=appointment_id)
    except Invoice.DoesNotExist:
        return None

    if invoice.status != Invoice.Status.CANCELLED:
        invoice.status = Invoice.Status.CANCELLED
        invoice.save(update_fields=["status", "updated_at"])

    return invoice


def _consume_loop():
    try:
        from kafka import KafkaConsumer
    except Exception as exc:  # pragma: no cover - defensive guard for container startup
        logger.warning("Kafka consumer unavailable: %s", exc)
        return

    consumer = KafkaConsumer(
        "appointment-created",
        "appointment-cancelled",
        bootstrap_servers=settings.KAFKA_BROKERS.split(","),
        value_deserializer=lambda value: json.loads(value.decode("utf-8")),
        auto_offset_reset="latest",
        enable_auto_commit=True,
        group_id=os.environ.get("PAYMENTS_KAFKA_GROUP_ID", "payments-service"),
        client_id="payments-service",
    )

    for message in consumer:
        payload = message.value or {}
        try:
            if message.topic == "appointment-created":
                handle_appointment_created(payload)
            elif message.topic == "appointment-cancelled":
                handle_appointment_cancelled(payload)
        except Exception as exc:  # pragma: no cover - keep consumer alive in Docker
            logger.exception("Failed to process %s event: %s", message.topic, exc)


def start_background_consumer():
    global _consumer_thread

    if not settings.KAFKA_ENABLED or not settings.PAYMENTS_CONSUMER_ENABLED:
        return
    if os.environ.get("RUN_MAIN") != "true":
        return

    with _consumer_lock:
        if _consumer_thread and _consumer_thread.is_alive():
            return
        _consumer_thread = threading.Thread(target=_consume_loop, name="payments-kafka-consumer", daemon=True)
        _consumer_thread.start()