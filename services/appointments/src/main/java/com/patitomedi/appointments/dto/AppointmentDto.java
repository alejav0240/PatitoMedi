package com.patitomedi.appointments.dto;

import java.time.OffsetDateTime;
import java.util.UUID;

public class AppointmentDto {
    public UUID id;
    public UUID patientId;
    public UUID doctorId;
    public UUID slotId;
    public OffsetDateTime startsAt;
    public OffsetDateTime endsAt;
    public String status;
    public String notes;
    public OffsetDateTime createdAt;

    public static AppointmentDto from(com.patitomedi.appointments.entity.Appointment a) {
        AppointmentDto dto = new AppointmentDto();
        dto.id = a.getId();
        dto.patientId = a.getPatientId();
        dto.doctorId = a.getDoctorId();
        dto.slotId = a.getSlot() != null ? a.getSlot().getId() : null;
        dto.startsAt = a.getStartsAt();
        dto.endsAt = a.getEndsAt();
        dto.status = a.getStatus();
        dto.notes = a.getNotes();
        dto.createdAt = a.getCreatedAt();
        return dto;
    }
}
