package com.patitomedi.appointments.dto;

import java.time.OffsetDateTime;
import java.util.UUID;

public class CreateSlotRequest {
    public UUID doctorId;
    public OffsetDateTime startsAt;
    public OffsetDateTime endsAt;
}
