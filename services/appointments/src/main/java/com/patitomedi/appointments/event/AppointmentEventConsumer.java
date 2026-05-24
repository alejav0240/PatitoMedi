package com.patitomedi.appointments.event;

import com.patitomedi.appointments.repository.AppointmentRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import org.springframework.transaction.annotation.Transactional;

import java.util.Map;
import java.util.UUID;

@Component
public class AppointmentEventConsumer {

    private static final Logger log = LoggerFactory.getLogger(AppointmentEventConsumer.class);

    private final AppointmentRepository appointmentRepo;

    public AppointmentEventConsumer(AppointmentRepository appointmentRepo) {
        this.appointmentRepo = appointmentRepo;
    }

    @KafkaListener(topics = "payment-confirmed")
    @Transactional
    public void onPaymentConfirmed(Map<String, Object> event) {
        UUID appointmentId = extractAppointmentId(event);
        if (appointmentId == null) return;
        appointmentRepo.findById(appointmentId).ifPresent(appt -> {
            appt.setStatus("confirmed");
            appointmentRepo.save(appt);
            log.info("Appointment {} confirmed via payment", appointmentId);
        });
    }

    @KafkaListener(topics = "payment-failed")
    @Transactional
    public void onPaymentFailed(Map<String, Object> event) {
        UUID appointmentId = extractAppointmentId(event);
        if (appointmentId == null) return;
        appointmentRepo.findById(appointmentId).ifPresent(appt -> {
            appt.setStatus("cancelled");
            appointmentRepo.save(appt);
            log.info("Appointment {} cancelled due to payment failure", appointmentId);
        });
    }

    @KafkaListener(topics = "call-ended")
    public void onCallEnded(Map<String, Object> event) {
        log.info("call-ended event received: {}", event);
    }

    private UUID extractAppointmentId(Map<String, Object> event) {
        try {
            Object id = event.get("appointment_id");
            return id != null ? UUID.fromString(id.toString()) : null;
        } catch (IllegalArgumentException e) {
            log.warn("Invalid appointment_id in event: {}", event);
            return null;
        }
    }
}
