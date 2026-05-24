package com.patitomedi.appointments.repository;

import com.patitomedi.appointments.entity.Appointment;
import org.springframework.data.jpa.repository.JpaRepository;
import java.util.List;
import java.util.UUID;

public interface AppointmentRepository extends JpaRepository<Appointment, UUID> {

    List<Appointment> findByPatientIdOrderByStartsAtAsc(UUID patientId);

    List<Appointment> findByDoctorIdOrderByStartsAtAsc(UUID doctorId);
}
