from prometheus_client import CONTENT_TYPE_LATEST, Counter, Histogram, generate_latest


REQUESTS_TOTAL = Counter(
    "payments_http_requests_total",
    "HTTP requests served by the payments service",
    ["endpoint", "method", "status"],
)
INVOICES_TOTAL = Counter("payments_invoices_total", "Invoices created by the payments service")
TRANSACTIONS_TOTAL = Counter("payments_transactions_total", "Transactions processed by the payments service", ["status"])
REFUNDS_TOTAL = Counter("payments_refunds_total", "Refunds processed by the payments service")
PROVIDER_LATENCY = Histogram(
    "payments_provider_latency_seconds",
    "Latency of provider-simulated operations in the payments service",
)


def metrics_bytes():
    return generate_latest()


def metrics_content_type():
    return CONTENT_TYPE_LATEST