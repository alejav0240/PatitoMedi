package com.patitomedi.appointments.dto;

import java.util.UUID;

public class CreateAppointmentRequest {
    public UUID patientId;
    public UUID doctorId;
    public UUID slotId;
    public String notes;
}
