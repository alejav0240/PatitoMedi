# Plan De Clientes

## Objetivo

Crear clientes que consuman la plataforma de telemedicina, empezando por web y dejando mobile para una fase posterior.

## Fases

1. Definir flujo MVP web.
2. Implementar autenticacion y sesion.
3. Implementar agenda, busqueda de medicos y creacion de citas.
4. Implementar pago simulado.
5. Implementar historial medico.
6. Implementar sala de video llamada con WebRTC.
7. Replicar flujos principales en mobile.

## Entregables Minimos

- Web app con login, registro y sesion.
- Pantalla de medicos y disponibilidad.
- Creacion y confirmacion de cita.
- Vista de historial medico.
- Entrada a sala de video por cita confirmada.

## Dependencias

- User Service para login y perfiles.
- Appointments Service para agenda.
- Payments Service para pago.
- Medical History Service para historial.
- Video Call Service para WebSocket y senalizacion.

## Interfaces Esperadas

- REST: `/api/users`, `/api/appointments`, `/api/payments`.
- GraphQL: `/graphql/medical-history`.
- WebSocket: `/ws/video`.

## Eventos Kafka Relevantes

Los clientes no consumen Kafka directamente. Deben reflejar estados expuestos por APIs que reaccionan a eventos internos.

## Criterios De Aceptacion

- El usuario puede completar el flujo MVP desde navegador.
- Las llamadas API pasan por gateway.
- Los errores de autenticacion, agenda y pago se muestran de forma clara.
- La pantalla de video intercambia mensajes WebRTC basicos.

## Orden Recomendado

1. Web MVP.
2. Integracion WebRTC.
3. Ajustes UX.
4. Mobile MVP.
