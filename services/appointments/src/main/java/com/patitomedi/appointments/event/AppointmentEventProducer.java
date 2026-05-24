package com.patitomedi.appointments.event;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

import java.util.Map;

@Component
public class AppointmentEventProducer {

    private static final Logger log = LoggerFactory.getLogger(AppointmentEventProducer.class);

    private final KafkaTemplate<String, Object> kafka;

    public AppointmentEventProducer(KafkaTemplate<String, Object> kafka) {
        this.kafka = kafka;
    }

    public void publish(String topic, Object payload) {
        kafka.send(topic, payload)
             .whenComplete((result, ex) -> {
                 if (ex != null) log.error("Failed to publish event to {}: {}", topic, ex.getMessage());
             });
    }
}
