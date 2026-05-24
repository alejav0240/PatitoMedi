package com.patitomedi.appointments.controller;

import com.patitomedi.appointments.dto.AppointmentDto;
import com.patitomedi.appointments.dto.CreateAppointmentRequest;
import com.patitomedi.appointments.dto.RescheduleRequest;
import com.patitomedi.appointments.service.AppointmentService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/")
public class AppointmentsController {

    private final AppointmentService service;

    public AppointmentsController(AppointmentService service) {
        this.service = service;
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public AppointmentDto create(@RequestBody CreateAppointmentRequest req) {
        return service.createAppointment(req);
    }

    @GetMapping("{id}")
    public AppointmentDto get(@PathVariable UUID id) {
        return service.getAppointment(id);
    }

    @PatchMapping("{id}/confirm")
    public AppointmentDto confirm(@PathVariable UUID id) {
        return service.confirm(id);
    }

    @PatchMapping("{id}/reschedule")
    public AppointmentDto reschedule(@PathVariable UUID id, @RequestBody RescheduleRequest req) {
        return service.reschedule(id, req);
    }

    @PatchMapping("{id}/cancel")
    public AppointmentDto cancel(@PathVariable UUID id) {
        return service.cancel(id);
    }

    @GetMapping("patients/{patientId}")
    public List<AppointmentDto> byPatient(@PathVariable UUID patientId) {
        return service.listByPatient(patientId);
    }

    @GetMapping("doctors/{doctorId}")
    public List<AppointmentDto> byDoctor(@PathVariable UUID doctorId) {
        return service.listByDoctor(doctorId);
    }
}
