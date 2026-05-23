# Video Call Service

Servicio responsable de señalización WebRTC, salas de atención y sesiones activas de video llamada.

## Stack previsto

- Lenguaje: Golang.
- Protocolo público: WebSocket.
- Infraestructura WebRTC: coturn para STUN/TURN.
- Estado temporal: Redis.
- Eventos: Apache Kafka.

## Responsabilidades

- Crear y administrar salas de video llamada.
- Manejar señalización WebRTC: `offer`, `answer` e ICE candidates.
- Coordinar participantes por cita médica.
- Mantener presencia y estado temporal de llamadas en Redis.
- Publicar eventos de inicio, cierre y fallas de llamada.

## API esperada

Ruta por gateway: `/ws/video`.

Mensajes WebSocket previstos:

- `join-room`: unir participante a una sala.
- `leave-room`: abandonar sala.
- `offer`: enviar oferta WebRTC.
- `answer`: enviar respuesta WebRTC.
- `ice-candidate`: intercambiar candidato ICE.
- `call-ended`: finalizar llamada.

## Datos temporales

Redis debe almacenar información efímera:

- `room:{appointmentId}`: estado de sala y participantes.
- `peer:{userId}`: conexión activa de un participante.
- `call:{callId}`: metadatos de llamada activa.

Las claves deben usar TTL para evitar sesiones colgadas.

## Eventos Kafka

Publica:

- `call-started`
- `call-ended`
- `call-failed`
- `call-recording-ready`

Consume:

- `appointment-confirmed`
- `appointment-cancelled`

## Integración coturn

El servicio debe entregar a los clientes la configuración ICE con servidores STUN/TURN. En desarrollo local, coturn queda disponible en `localhost:3478`.

## Observabilidad

Debe exponer `/metrics` para Prometheus con métricas de salas activas, peers conectados, duración de llamadas, fallas ICE y cierres inesperados.
