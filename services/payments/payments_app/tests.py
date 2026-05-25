from django.test import TestCase
from unittest.mock import patch

from .consumers import handle_appointment_cancelled, handle_appointment_created
from .models import Invoice, Transaction


class PaymentsApiTests(TestCase):
    def test_health(self):
        response = self.client.get("/health")
        self.assertEqual(response.status_code, 200)

    def test_metrics(self):
        response = self.client.get("/metrics")
        self.assertEqual(response.status_code, 200)
        self.assertIn(b"payments_http_requests_total", response.content)

    def test_invoice_transaction_and_refund_flow(self):
        invoice_response = self.client.post(
            "/api/payments/invoices",
            data='{"appointment_id":"apt-1","patient_id":"pat-1","doctor_id":"doc-1","amount":"120.50","currency":"PEN"}',
            content_type="application/json",
        )
        self.assertEqual(invoice_response.status_code, 201)
        invoice_id = invoice_response.json()["id"]

        transaction_response = self.client.post(
            "/api/payments/transactions",
            data=f'{{"invoice_id":"{invoice_id}","simulate_outcome":"approved"}}',
            content_type="application/json",
        )
        self.assertEqual(transaction_response.status_code, 201)
        transaction_id = transaction_response.json()["id"]

        refund_response = self.client.post(
            "/api/payments/refunds",
            data=f'{{"transaction_id":"{transaction_id}","reason":"customer_request"}}',
            content_type="application/json",
        )
        self.assertEqual(refund_response.status_code, 201)

        invoice = Invoice.objects.get(pk=invoice_id)
        transaction = Transaction.objects.get(pk=transaction_id)
        self.assertEqual(invoice.status, Invoice.Status.PAID)
        self.assertEqual(transaction.status, Transaction.Status.REFUNDED)

    @patch("payments_app.consumers.publish_event")
    def test_appointment_events_create_and_cancel_invoice(self, _publish_event):
        invoice = handle_appointment_created(
            {
                "appointment_id": "apt-100",
                "patient_id": "pat-100",
                "doctor_id": "doc-100",
                "amount": "75.00",
                "currency": "PEN",
            }
        )

        self.assertEqual(invoice.appointment_id, "apt-100")
        self.assertEqual(invoice.status, Invoice.Status.ISSUED)

        cancelled_invoice = handle_appointment_cancelled({"appointment_id": "apt-100"})

        self.assertIsNotNone(cancelled_invoice)
        self.assertEqual(cancelled_invoice.status, Invoice.Status.CANCELLED)