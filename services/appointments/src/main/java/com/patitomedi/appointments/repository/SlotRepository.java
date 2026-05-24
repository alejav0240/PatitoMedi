package com.patitomedi.appointments.repository;

import com.patitomedi.appointments.entity.Slot;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import java.util.List;
import java.util.UUID;

public interface SlotRepository extends JpaRepository<Slot, UUID> {

    List<Slot> findByDoctorIdAndIsAvailableTrue(UUID doctorId);

    List<Slot> findByIsAvailableTrue();

    @Query("SELECT s FROM Slot s WHERE (:doctorId IS NULL OR s.doctorId = :doctorId) AND s.isAvailable = :available")
    List<Slot> findByOptionalDoctorAndAvailability(UUID doctorId, boolean available);
}
