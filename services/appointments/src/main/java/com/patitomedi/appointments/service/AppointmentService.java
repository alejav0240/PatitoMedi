package com.patitomedi.appointments.service;

import com.patitomedi.appointments.dto.*;
import com.patitomedi.appointments.entity.Appointment;
import com.patitomedi.appointments.entity.Slot;
import com.patitomedi.appointments.event.AppointmentEventProducer;
import com.patitomedi.appointments.exception.NotFoundException;
import com.patitomedi.appointments.exception.SlotUnavailableException;
import com.patitomedi.appointments.repository.AppointmentRepository;
import com.patitomedi.appointments.repository.SlotRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

@Service
public class AppointmentService {

    private final SlotRepository slotRepo;
    private final AppointmentRepository appointmentRepo;
    private final AppointmentEventProducer eventProducer;

    public AppointmentService(SlotRepository slotRepo, AppointmentRepository appointmentRepo,
                              AppointmentEventProducer eventProducer) {
        this.slotRepo = slotRepo;
        this.appointmentRepo = appointmentRepo;
        this.eventProducer = eventProducer;
    }

    // --- Slots ---

    @Transactional
    public SlotDto createSlot(CreateSlotRequest req) {
        Slot slot = new Slot();
        slot.setDoctorId(req.doctorId);
        slot.setStartsAt(req.startsAt);
        slot.setEndsAt(req.endsAt);
        return SlotDto.from(slotRepo.save(slot));
    }

    public List<SlotDto> listSlots(UUID doctorId, boolean onlyAvailable) {
        return slotRepo.findByOptionalDoctorAndAvailability(doctorId, onlyAvailable)
                .stream().map(SlotDto::from).toList();
    }

    // --- Appointments ---

    @Transactional
    public AppointmentDto createAppointment(CreateAppointmentRequest req) {
        Slot slot = slotRepo.findById(req.slotId)
                .orElseThrow(() -> new NotFoundException("slot not found"));
        if (!slot.isAvailable()) throw new SlotUnavailableException();

        slot.setAvailable(false);
        slotRepo.save(slot);

        Appointment appt = new Appointment();
        appt.setPatientId(req.patientId);
        appt.setDoctorId(req.doctorId);
        appt.setSlot(slot);
        appt.setStartsAt(slot.getStartsAt());
        appt.setEndsAt(slot.getEndsAt());
        appt.setNotes(req.notes);

        AppointmentDto dto = AppointmentDto.from(appointmentRepo.save(appt));
        eventProducer.publish("appointment-created", dto);
        return dto;
    }

    public AppointmentDto getAppointment(UUID id) {
        return appointmentRepo.findById(id)
                .map(AppointmentDto::from)
                .orElseThrow(() -> new NotFoundException("appointment not found"));
    }

    @Transactional
    public AppointmentDto confirm(UUID id) {
        Appointment appt = appointmentRepo.findById(id)
                .orElseThrow(() -> new NotFoundException("appointment not found"));
        appt.setStatus("confirmed");
        AppointmentDto dto = AppointmentDto.from(appointmentRepo.save(appt));
        eventProducer.publish("appointment-confirmed", dto);
        return dto;
    }

    @Transactional
    public AppointmentDto reschedule(UUID id, RescheduleRequest req) {
        Appointment appt = appointmentRepo.findById(id)
                .orElseThrow(() -> new NotFoundException("appointment not found"));

        Slot newSlot = slotRepo.findById(req.slotId)
                .orElseThrow(() -> new NotFoundException("slot not found"));
        if (!newSlot.isAvailable()) throw new SlotUnavailableException();

        // Free old slot
        Slot oldSlot = appt.getSlot();
        if (oldSlot != null) {
            oldSlot.setAvailable(true);
            slotRepo.save(oldSlot);
        }

        newSlot.setAvailable(false);
        slotRepo.save(newSlot);

        appt.setSlot(newSlot);
        appt.setStartsAt(newSlot.getStartsAt());
        appt.setEndsAt(newSlot.getEndsAt());
        appt.setStatus("rescheduled");

        AppointmentDto dto = AppointmentDto.from(appointmentRepo.save(appt));
        eventProducer.publish("appointment-rescheduled", dto);
        return dto;
    }

    @Transactional
    public AppointmentDto cancel(UUID id) {
        Appointment appt = appointmentRepo.findById(id)
                .orElseThrow(() -> new NotFoundException("appointment not found"));

        Slot slot = appt.getSlot();
        if (slot != null) {
            slot.setAvailable(true);
            slotRepo.save(slot);
        }

        appt.setStatus("cancelled");
        AppointmentDto dto = AppointmentDto.from(appointmentRepo.save(appt));
        eventProducer.publish("appointment-cancelled", dto);
        return dto;
    }

    public List<AppointmentDto> listByPatient(UUID patientId) {
        return appointmentRepo.findByPatientIdOrderByStartsAtAsc(patientId)
                .stream().map(AppointmentDto::from).toList();
    }

    public List<AppointmentDto> listByDoctor(UUID doctorId) {
        return appointmentRepo.findByDoctorIdOrderByStartsAtAsc(doctorId)
                .stream().map(AppointmentDto::from).toList();
    }
}
