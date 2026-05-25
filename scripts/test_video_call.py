#!/usr/bin/env python3
"""
Integration test for video-call WebSocket service.
Tests: join-room, offer/answer routing, call-ended, error handling.
"""
import json
import threading
import time
import sys
import websocket

APPOINTMENT_ID = "00000000-0000-0000-0000-000000000001"
PATIENT_ID     = "00000000-0000-0000-0000-000000000002"
DOCTOR_ID      = "00000000-0000-0000-0000-000000000003"
WS_URL         = "ws://localhost:8000/ws/video"

results = {}
received = {"patient": [], "doctor": []}
lock = threading.Lock()

def make_ws(name, on_message_cb):
    ws = websocket.WebSocket()
    ws.connect(WS_URL)
    def reader():
        try:
            while True:
                msg = ws.recv()
                with lock:
                    on_message_cb(json.loads(msg))
        except Exception:
            pass
    t = threading.Thread(target=reader, daemon=True)
    t.start()
    return ws

def send(ws, payload):
    ws.send(json.dumps(payload))
    time.sleep(0.2)

def check(name, condition, detail=""):
    status = "PASS" if condition else "FAIL"
    results[name] = status
    print(f"  [{status}] {name}" + (f" — {detail}" if detail else ""))

def run():
    print(f"\nTarget: {WS_URL}\n")

    # --- Test 1: health via Kong (HTTP GET to /ws/video returns Bad Request, not 404) ---
    import urllib.request, urllib.error
    try:
        urllib.request.urlopen("http://localhost:8000/ws/video", timeout=3)
    except urllib.error.HTTPError as e:
        # 400 Bad Request = Kong is routing correctly (missing WS upgrade)
        check("health via Kong routing", e.code == 400, f"HTTP {e.code}")
    except Exception as e:
        check("health via Kong routing", False, str(e))

    # --- Test 2: connect patient ---
    try:
        patient_ws = make_ws("patient", lambda m: received["patient"].append(m))
        check("patient WS connect", True)
    except Exception as e:
        check("patient WS connect", False, str(e))
        sys.exit(1)

    # --- Test 3: connect doctor ---
    try:
        doctor_ws = make_ws("doctor", lambda m: received["doctor"].append(m))
        check("doctor WS connect", True)
    except Exception as e:
        check("doctor WS connect", False, str(e))
        sys.exit(1)

    # --- Test 4: join-room patient ---
    send(patient_ws, {"type": "join-room", "appointmentId": APPOINTMENT_ID, "userId": PATIENT_ID})
    check("patient join-room (no error)", not any(m.get("type") == "error" for m in received["patient"]))

    # --- Test 5: join-room doctor ---
    send(doctor_ws, {"type": "join-room", "appointmentId": APPOINTMENT_ID, "userId": DOCTOR_ID})
    check("doctor join-room (no error)", not any(m.get("type") == "error" for m in received["doctor"]))

    # --- Test 6: offer routing (patient → doctor) ---
    received["doctor"].clear()
    send(patient_ws, {
        "type": "offer",
        "appointmentId": APPOINTMENT_ID,
        "sdp": "v=0\r\no=- 0 0 IN IP4 127.0.0.1\r\ns=-\r\n",
        "to": DOCTOR_ID
    })
    time.sleep(0.3)
    got_offer = any(m.get("type") == "offer" for m in received["doctor"])
    check("offer routed to doctor", got_offer, f"doctor received: {received['doctor']}")

    # --- Test 7: answer routing (doctor → patient) ---
    received["patient"].clear()
    send(doctor_ws, {
        "type": "answer",
        "appointmentId": APPOINTMENT_ID,
        "sdp": "v=0\r\no=- 1 1 IN IP4 127.0.0.1\r\ns=-\r\n",
        "to": PATIENT_ID
    })
    time.sleep(0.3)
    got_answer = any(m.get("type") == "answer" for m in received["patient"])
    check("answer routed to patient", got_answer, f"patient received: {received['patient']}")

    # --- Test 8: ice-candidate routing ---
    received["doctor"].clear()
    send(patient_ws, {
        "type": "ice-candidate",
        "appointmentId": APPOINTMENT_ID,
        "candidate": "candidate:1 1 UDP 2130706431 192.168.1.1 54321 typ host",
        "sdpMid": "0",
        "sdpMLineIndex": 0,
        "to": DOCTOR_ID
    })
    time.sleep(0.3)
    got_ice = any(m.get("type") == "ice-candidate" for m in received["doctor"])
    check("ice-candidate routed to doctor", got_ice)

    # --- Test 9: error on unknown type ---
    received["patient"].clear()
    send(patient_ws, {"type": "unknown-type"})
    time.sleep(0.2)
    got_error = any(m.get("type") == "error" for m in received["patient"])
    check("error on unknown message type", got_error, f"received: {received['patient']}")

    # --- Test 10: call-ended ---
    send(patient_ws, {"type": "call-ended", "appointmentId": APPOINTMENT_ID})
    check("call-ended sent (no crash)", True)

    patient_ws.close()
    doctor_ws.close()

    # Summary
    passed = sum(1 for v in results.values() if v == "PASS")
    total  = len(results)
    print(f"\n{'='*40}")
    print(f"  {passed}/{total} tests passed")
    print(f"{'='*40}\n")
    sys.exit(0 if passed == total else 1)

if __name__ == "__main__":
    run()
