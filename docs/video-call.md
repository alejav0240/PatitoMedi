# Video Call Service

Servicio de señalización WebRTC para video llamadas médicas.

**Stack:** Go · WebSocket (`gorilla/websocket`) · Redis · Kafka · coturn (STUN/TURN)

**Gateway base:** `/ws/video` (expuesto a través de Kong)

Estado: ✅ Implementado (servicio real disponible en Compose)

## Resumen

Este microservicio gestiona la señalización WebRTC entre dos participantes (paciente y médico) dentro de una sala asociada a una cita (`appointmentId`). No transmite media — actúa como relay de mensajes SDP e ICE candidates entre peers.

Dentro del contenedor el servidor escucha en el puerto `8080`. Kong enruta `ws://localhost:8000/ws/video` hacia el servicio con `strip_path: false`.

## Protocolo WebSocket

Todos los mensajes son JSON sobre una conexión WebSocket.

### Conectar

```
ws://localhost:8000/ws/video
```

Headers requeridos por el upgrade:
```
Upgrade: websocket
Connection: Upgrade
```

### Tipos de mensajes

| Tipo | Dirección | Descripción |
|------|-----------|-------------|
| `join-room` | Cliente → Servidor | Unirse a la sala de una cita |
| `leave-room` | Cliente → Servidor | Abandonar la sala |
| `offer` | Cliente → Servidor | Oferta SDP WebRTC (enrutada al peer destino) |
| `answer` | Cliente → Servidor | Respuesta SDP WebRTC (enrutada al peer destino) |
| `ice-candidate` | Cliente ↔ Servidor | Candidato ICE (enrutado al peer destino) |
| `call-ended` | Cliente → Servidor | Finalizar llamada |
| `error` | Servidor → Cliente | Error de señalización |

### Flujo típico

```
Paciente                    Servidor                    Médico
   |                           |                           |
   |-- join-room ------------->|                           |
   |                           |<---------- join-room -----|
   |                           |  (publica call-started)   |
   |-- offer (to: doctorId) -->|-- offer ----------------->|
   |                           |<-- answer (to: patientId)-|
   |<-- answer ----------------|                           |
   |-- ice-candidate --------->|-- ice-candidate --------->|
   |<-- ice-candidate ----------|<---------- ice-candidate-|
   |-- call-ended ------------>|                           |
   |                           |  (publica call-ended)     |
```

### Mensajes de ejemplo

**join-room**
```json
{
  "type": "join-room",
  "appointmentId": "00000000-0000-0000-0000-000000000001",
  "userId": "00000000-0000-0000-0000-000000000002"
}
```

**offer**
```json
{
  "type": "offer",
  "appointmentId": "00000000-0000-0000-0000-000000000001",
  "sdp": "v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\n...",
  "to": "00000000-0000-0000-0000-000000000003"
}
```

**answer**
```json
{
  "type": "answer",
  "appointmentId": "00000000-0000-0000-0000-000000000001",
  "sdp": "v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\n...",
  "to": "00000000-0000-0000-0000-000000000002"
}
```

**ice-candidate**
```json
{
  "type": "ice-candidate",
  "appointmentId": "00000000-0000-0000-0000-000000000001",
  "candidate": "candidate:1 1 UDP 2130706431 192.168.1.1 54321 typ host",
  "sdpMid": "0",
  "sdpMLineIndex": 0,
  "to": "00000000-0000-0000-0000-000000000003"
}
```

**call-ended**
```json
{
  "type": "call-ended",
  "appointmentId": "00000000-0000-0000-0000-000000000001"
}
```

**error** (servidor → cliente)
```json
{
  "type": "error",
  "message": "not in a room"
}
```

## Endpoints HTTP

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/ws/video` | WebSocket upgrade |
| `GET` | `/health` | Health check |
| `GET` | `/metrics` | Métricas Prometheus |

**Health check** (directo al contenedor):
```bash
docker exec patitomedi-video-call-service wget -qO- http://localhost:8080/health
# {"service":"video-call","status":"ok"}
```

## Redis — estado temporal

Las salas activas se persisten en Redis con TTL de 2 horas.

| Clave | Tipo | Contenido |
|-------|------|-----------|
| `room:{appointmentId}` | Hash | `status`, `started_at` |
| `room:{appointmentId}:participants` | Set | userIds conectados |

```bash
# Inspeccionar sala activa
docker exec patitomedi-redis redis-cli HGETALL room:00000000-0000-0000-0000-000000000001
docker exec patitomedi-redis redis-cli SMEMBERS room:00000000-0000-0000-0000-000000000001:participants
```

## Kafka — eventos publicados

| Topic | Cuándo | Payload |
|-------|--------|---------|
| `call-started` | Cuando el 2do peer entra a la sala | `appointmentId`, `participants`, `startedAt` |
| `call-ended` | Cuando un peer envía `call-ended` o se desconecta y la sala queda vacía | `appointmentId`, `initiatorId`, `endedAt` |
| `call-failed` | Error crítico de señalización | `appointmentId`, `reason`, `failedAt` |

```bash
# Consumir eventos en tiempo real
docker exec patitomedi-kafka \
  /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 --topic call-started --from-beginning
```

## Métricas Prometheus

| Métrica | Tipo | Descripción |
|---------|------|-------------|
| `video_active_rooms` | Gauge | Salas activas en este momento |
| `video_active_participants` | Gauge | Participantes conectados en total |
| `video_calls_total` | Counter | Total de llamadas iniciadas |

## STUN/TURN — coturn

El servidor coturn corre en el puerto `3478` (TCP y UDP). Configuración en `infra/coturn/turnserver.conf`.

Credenciales locales:
- **realm:** `patitomedi.local`
- **user:** `patito` / `patito_turn_secret`

Configuración ICE para el cliente WebRTC:
```json
{
  "iceServers": [
    { "urls": "stun:localhost:3478" },
    {
      "urls": "turn:localhost:3478",
      "username": "patito",
      "credential": "patito_turn_secret"
    }
  ]
}
```

## Variables de entorno

| Variable | Default | Descripción |
|----------|---------|-------------|
| `PORT` | `8080` | Puerto del servidor |
| `REDIS_ADDR` | `localhost:6379` | Dirección Redis |
| `KAFKA_BROKERS` | `localhost:9092` | Brokers Kafka (separados por coma) |

## Probar con Insomnia

1. Crear dos requests de tipo **WebSocket** en Insomnia.
2. Ambas apuntan a `ws://localhost:8000/ws/video`.
3. Conectar ambas y enviar `join-room` con el mismo `appointmentId` pero distinto `userId`.
4. Desde la primera conexión enviar `offer` con `"to": "<userId de la segunda>"`.
5. Verificar que la segunda conexión recibe el mensaje.

## Pruebas automatizadas

```bash
python3 scripts/test_video_call.py
# 10/10 tests passed
```

## Notas

- El servicio no valida que el `appointmentId` exista en el Appointments Service (validación futura).
- El enrutamiento de mensajes usa el campo `to` (userId del peer destino). Si `to` está ausente, el mensaje se hace broadcast a todos los peers de la sala excepto el emisor.
- Kong tiene `strip_path: false` para esta ruta — el upstream recibe el path completo `/ws/video`.
