import json
import logging

from django.conf import settings

logger = logging.getLogger(__name__)

_producer = None


def _get_producer():
    global _producer
    if _producer is not None:
        return _producer
    if not settings.KAFKA_ENABLED:
        return None
    try:
        from kafka import KafkaProducer

        _producer = KafkaProducer(
            bootstrap_servers=settings.KAFKA_BROKERS.split(","),
            value_serializer=lambda value: json.dumps(value).encode("utf-8"),
            acks="all",
            retries=1,
        )
    except Exception as exc:
        logger.warning("Kafka producer unavailable: %s", exc)
        _producer = None
    return _producer


def publish_event(topic, payload):
    producer = _get_producer()
    if producer is None:
        return False
    try:
        producer.send(topic, payload)
        producer.flush(timeout=2)
        return True
    except Exception as exc:
        logger.warning("Kafka publish failed for %s: %s", topic, exc)
        return False