package com.patitomedi.appointments.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.CONFLICT)
public class SlotUnavailableException extends RuntimeException {
    public SlotUnavailableException() { super("slot is not available"); }
}
