# water pump regulator

so basically i built a water tank monitor that uses an ESP32 and an ultrasonic sensor to check if the tank is full or empty, then sends an MQTT message so a Go server can do something about it (in this case, notify via WhatsApp). pretty simple stuff.

---

## how it works

the ESP32 reads the distance from the water surface using a HC-SR04 ultrasonic sensor. if the water is low, it publishes to `message_pump_on` (turn pump on). if the water is high enough, it publishes to `message_pump_off` (turn pump off).

the Go server runs an embedded MQTT broker so the ESP32 connects directly to your machine. the server subscribes to both topics and can forward the messages to WhatsApp using the Meta Cloud API.

i also added a retry system on the Arduino side because the water turbulence was making the pump toggle on and off rapidly. not the cleanest fix but it works.

---

## the parts

- **ESP32** — the microcontroller
- **HC-SR04** — ultrasonic distance sensor (TRIG on pin 1, ECHO on pin 3)
- a pump or relay or whatever you're trying to control

---

## tools i had to install

### Arduino side
- [Arduino IDE](https://www.arduino.cc/en/software)
- **WiFi** — comes with the ESP32 board package
- **PubSubClient** — MQTT client library for Arduino. install it from the library manager in the IDE

### Go server side
you need Go installed. then run this in the `mqtt_broker` folder:

```bash
go mod tidy
```

that pulls everything. the main packages are:
- `mochi-mqtt/server` — this is the embedded MQTT broker, so you don't need to install mosquitto or anything external
- `paho.mqtt.golang` — the Go MQTT client that connects to the broker
- `gowhatsapp` — for sending WhatsApp messages via the Meta Cloud API
- `godotenv` — loads credentials from a `.env` file so i don't have to hardcode them

---

## setup

### 1. clone the repo and go into the broker folder

```bash
cd mqtt_broker
```

### 2. create a `.env` file

```
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_ACCESS_TOKEN=your_access_token
```

get these from [Meta's developer platform](https://developers.facebook.com). create an app, add WhatsApp, go to API Setup.

### 3. find your machine's local IP

on Windows:
```bash
ipconfig
```
look for the IPv4 address under your WiFi adapter. it usually starts with `192.168.x.x`.

### 4. update the Arduino sketch

in `water_pump_regulator.ino`, set your WiFi credentials and your machine's IP:

```cpp
const char *ssid = "your_wifi_name";
const char *password = "your_wifi_password";
const char *mqtt_server = "192.168.x.x"; // your machine's IP
```

make sure your machine and the ESP32 are on the same WiFi network or this won't work, broski.

### 5. open port 1883 on windows firewall

this one got me. the ESP32 kept failing to connect with `rc=-2` and i could not figure out why for a while. turns out Windows Firewall was just blocking the port the whole time, lol.

open PowerShell **as administrator** (right-click → run as administrator) and run:

```powershell
New-NetFirewallRule -DisplayName "MQTT Broker" -Direction Inbound -Protocol TCP -LocalPort 1883 -Action Allow
```

you only need to do this once. if you skip this step, the ESP32 will not connect, full stop.

### 6. run the Go server

```bash
go run .
```

### 6. flash the ESP32

open `water_pump_regulator.ino` in the Arduino IDE and upload it.

---

## project structure

```
.
├── mqtt_broker/          # the Go server
│   ├── main.go
│   ├── .env
│   └── internal/
│       ├── config/       # loads .env
│       ├── mqtt/         # embedded broker + subscriber
│       └── whatsapp/     # sends WhatsApp messages
└── water_pump_regulator/
    └── water_pump_regulator.ino  # the Arduino sketch
```

