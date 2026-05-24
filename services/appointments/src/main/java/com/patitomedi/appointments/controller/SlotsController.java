package com.patitomedi.appointments.controller;

import com.patitomedi.appointments.dto.CreateSlotRequest;
import com.patitomedi.appointments.dto.SlotDto;
import com.patitomedi.appointments.service.AppointmentService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/slots")
public class SlotsController {

    private final AppointmentService service;

    public SlotsController(AppointmentService service) {
        this.service = service;
    }

    @GetMapping
    public List<SlotDto> list(
            @RequestParam(required = false) UUID doctorId,
            @RequestParam(defaultValue = "true") boolean available) {
        return service.listSlots(doctorId, available);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public SlotDto create(@RequestBody CreateSlotRequest req) {
        return service.createSlot(req);
    }
}
