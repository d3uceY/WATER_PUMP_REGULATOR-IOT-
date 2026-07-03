# water pump regulator

so basically i built a water tank monitor that uses an ESP32 and an ultrasonic sensor to check if the tank is full or empty, then sends an MQTT message so a Go server can do something about it. right now the Go server sends the alert to WhatsApp and Telegram, because apparently one notification app was not enough.

---

## how it works

the ESP32 reads the distance from the water surface using a HC-SR04 ultrasonic sensor. if the water is low, it publishes to `message_pump_on` (turn pump on). if the water is high enough, it publishes to `message_pump_off` (turn pump off).

the Go server runs an embedded MQTT broker so the ESP32 connects directly to your machine. the server subscribes to both topics and forwards the messages to WhatsApp using the Meta Cloud API, and to Telegram using a bot you make yourself.

when the broker receives `message_pump_on`, it sends `pump is turned on`. when it receives `message_pump_off`, it sends `pump is turned off`.

i also added a retry system on the Arduino side because the water turbulence was making the pump toggle on and off rapidly. not the cleanest fix but it works.

---

## the parts

- **ESP32** — the microcontroller
- **HC-SR04** — ultrasonic distance sensor (TRIG on pin 1, ECHO on pin 3)
- a pump or relay or whatever you're trying to control

---

## pin diagram

this is the wiring i used in the sketch. nothing fancy, just don't mix up trig and echo or you will be debugging the wrong thing for no reason. I did that

```text
ESP32                         HC-SR04
-----                         -------
5V / VIN  ------------------> VCC
GND      ------------------> GND
GPIO 1   ------------------> TRIG
GPIO 3   <------------------ ECHO
```
I did not do this, btw because i am lazy and i fried some pins😂😂😂
small warning: the HC-SR04 echo pin can output 5V, and most ESP32 C3 pins want 3.3V. use a voltage divider on `ECHO` if your sensor board does not already handle that. something like this works: 

```text
HC-SR04 ECHO ---- 1k resistor ----+---- ESP32 GPIO 3
                                  |
                              2k resistor
                                  |
                                 GND
```

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
- Telegram Bot API — for sending Telegram messages to whoever has chatted with your bot
- `godotenv` — loads credentials from a `.env` file so i don't have to hardcode them

---

## setup

### 1. clone the repo and go into the broker folder

```bash
cd mqtt_broker
```

### 2. create a `.env` file

copy the example file so you don't have to freestyle the env names:

```bash
cp .env.example .env
```

on Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

then fill in the values with your actual stuff:

```
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_ACCESS_TOKEN=your_access_token
WHATSAPP_RECIPIENT_NUMBER=recipient_phone_number_with_country_code

TELEGRAM_BOT_TOKEN=your_telegram_bot_token
```

do not commit your real `.env` file. keep the token stuff private unless you enjoy pain.

### 3. set up WhatsApp Cloud API

get the WhatsApp values from [Meta's developer platform](https://developers.facebook.com). the flow is basically:

1. create or open a Meta app
2. add WhatsApp to the app
3. go to **WhatsApp > API Setup**
4. copy the phone number ID into `WHATSAPP_PHONE_NUMBER_ID`
5. copy the access token into `WHATSAPP_ACCESS_TOKEN`
6. set `WHATSAPP_RECIPIENT_NUMBER` to the phone number you want to notify, including country code

for testing, Meta might make you add the recipient phone number as a test recipient first. if messages are not sending, check that before blaming the code.

### 4. make your own Telegram bot

this part is actually pretty chill:

1. open Telegram and search for `@BotFather`
2. send `/newbot`
3. choose a bot name and username
4. copy the bot token into `TELEGRAM_BOT_TOKEN`
5. open your new bot in Telegram and send it any message, like `start`

that last step matters because the app gets chat IDs from Telegram updates. if nobody has messaged the bot yet, the Go server has no chat IDs to send to, so it just has nobody to bother.

if you want the bot to notify a group, add the bot to the group and send a message in that group so Telegram creates an update for it. after that, the bot can pick up the group chat id too.

### 5. find your machine's local IP

on Windows:
```bash
ipconfig
```
look for the IPv4 address under your WiFi adapter. it usually starts with `192.168.x.x`.

### 6. update the Arduino sketch

in `water_pump_regulator.ino`, set your WiFi credentials and your machine's IP:

```cpp
const char *ssid = "your_wifi_name";
const char *password = "your_wifi_password";
const char *mqtt_server = "192.168.x.x"; // your machine's IP
```

make sure your machine and the ESP32 are on the same WiFi network or this won't work, broski.

### 7. open port 1883 on windows firewall

this one got me. the ESP32 kept failing to connect with `rc=-2` and i could not figure out why for a while. turns out Windows Firewall was just blocking the port the whole time, lol.

open PowerShell **as administrator** (right-click → run as administrator) and run:

```powershell
New-NetFirewallRule -DisplayName "MQTT Broker" -Direction Inbound -Protocol TCP -LocalPort 1883 -Action Allow
```

you only need to do this once. if you skip this step, the ESP32 will not connect, full stop.

### 8. run the Go server

```bash
go run .
```

### 9. flash the ESP32

open `water_pump_regulator.ino` in the Arduino IDE and upload it.

---

## project structure

```
.
├── mqtt_broker/          # the Go server
│   ├── main.go
│   ├── .env.example
│   └── internal/
│       ├── config/       # loads .env
│       ├── mqtt/         # embedded broker + subscriber + notification router thing
│       ├── telegram/     # sends Telegram bot messages
│       └── whatsapp/     # sends WhatsApp messages
└── water_pump_regulator/
    └── water_pump_regulator.ino  # the Arduino sketch
```

