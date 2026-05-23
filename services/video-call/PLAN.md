# Plan De Desarrollo: Video Call Service

## Objetivo

Implementar senalizacion WebRTC, salas por cita y estado temporal de llamadas con Redis.

## Fases

1. Crear proyecto Go con servidor WebSocket.
2. Implementar rooms por `appointmentId`.
3. Manejar mensajes `join-room`, `offer`, `answer`, `ice-candidate` y `leave-room`.
4. Persistir estado temporal en Redis con TTL.
5. Integrar configuracion ICE con coturn.
6. Publicar eventos de inicio, cierre y falla de llamada.

## Entregables Minimos

- Servicio reemplazando el stub en compose mediante `build: ./services/video-call`.
- Endpoint WebSocket.
- Rooms con dos o mas participantes.
- Redis para presencia y estado temporal.
- Eventos `call-started` y `call-ended`.
- `/health` y `/metrics`.

## Dependencias

- Redis.
- Kafka.
- coturn.
- Appointments Service para citas confirmadas.
- User Service para identidad.

## Interfaces Esperadas

Ruta gateway: `/ws/video`.

Mensajes WebSocket:

- `join-room`
- `leave-room`
- `offer`
- `answer`
- `ice-candidate`
- `call-ended`

## Eventos Kafka

Publica:

- `call-started`
- `call-ended`
- `call-failed`
- `call-recording-ready`

Consume:

- `appointment-confirmed`
- `appointment-cancelled`

## Criterios De Aceptacion

- Dos participantes pueden entrar a la misma sala.
- El servicio reenvia offer, answer e ICE candidates al peer correcto.
- Redis limpia salas inactivas con TTL.
- El cierre de llamada publica `call-ended`.
- `/health` y `/metrics` responden.

## Orden Recomendado

1. Bootstrap Go WebSocket.
2. Rooms en memoria.
3. Redis y TTL.
4. coturn e ICE config.
5. Kafka y observabilidad.
