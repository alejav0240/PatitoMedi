package com.patitomedi.appointments.dto;

import java.time.OffsetDateTime;
import java.util.UUID;

public class SlotDto {
    public UUID id;
    public UUID doctorId;
    public OffsetDateTime startsAt;
    public OffsetDateTime endsAt;
    public boolean isAvailable;

    public static SlotDto from(com.patitomedi.appointments.entity.Slot s) {
        SlotDto dto = new SlotDto();
        dto.id = s.getId();
        dto.doctorId = s.getDoctorId();
        dto.startsAt = s.getStartsAt();
        dto.endsAt = s.getEndsAt();
        dto.isAvailable = s.isAvailable();
        return dto;
    }
}
