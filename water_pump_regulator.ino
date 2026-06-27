#include <stdint.h>
#include <WiFi.h>
#include <PubSubClient.h>
#define TRIG_PIN 1
#define ECHO_PIN 3

// wifi information and shit
const char *password = "lhs814185";
const char *ssid = "MTN-LHS";

// MQTT broker
// i hope i do not forget to make a server lol
// i will use this for the test
const char *mqtt_server = "broker.emqx.io";

WiFiClient espClient;
PubSubClient client(espClient);

int8_t threshold = 20;
int8_t maxRetries = 8;
int8_t retryCount = 0;
uint8_t delayMillis = 200;
bool pumpOn;

void reconnectMQTT() {

  while (!client.connected()) {

    Serial.print("Connecting to MQTT...");

    String clientId = "ESP32-" + String((uint32_t)ESP.getEfuseMac(), HEX);

    if (client.connect(clientId.c_str())) {

      Serial.println("Connected!");
    } else {

      Serial.print("Failed. rc=");
      Serial.print(client.state());
      Serial.println(" retrying...");
      delay(2000);
    }
  }
}

void connectWIFI() {
  // wifi stuff
  WiFi.begin(ssid, password);
  Serial.print("Connecting");

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println();
  Serial.println("Wifi is connected, broski!");
  Serial.print("IP Address: ");
  Serial.println(WiFi.localIP());
}

void setup() {
  Serial.begin(115200);

  connectWIFI();

  client.setServer(mqtt_server, 1883);

  pinMode(ECHO_PIN, INPUT);
  pinMode(TRIG_PIN, OUTPUT);
  Serial.print("it has started");
}

long readDistanceCM() {
  digitalWrite(TRIG_PIN, LOW);
  delayMicroseconds(2);

  digitalWrite(TRIG_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIG_PIN, LOW);

  long duration = pulseIn(ECHO_PIN, HIGH);

  // had to look this shit up, forgot my physics
  long distance = duration * 0.0343 / 2;

  return distance;
}

void loop() {

  // reconnect MQTT
  if (!client.connected()) {
    reconnectMQTT();
  }

  client.loop();

  client.publish("deuce/esp32/test", "Hello from ESP32!");

  Serial.println("Message Published!");

  delay(5000);

  long distance = readDistanceCM();

  if (distance <= 2 || distance > 400) {
    delay(delayMillis);
    return;
  }

  if (distance >= (20 + threshold)) {
    retryCount = 0;
    if (!pumpOn) {
      pumpOn = true;
      Serial.print("pump on\n distance: ");
      Serial.print(distance);
      Serial.print(" cm \n");
    }
  } else if (distance < 20) {
    retryCount += 1;
    if (retryCount < maxRetries) {
      Serial.print("retrying... ");
      Serial.print(retryCount);
      Serial.print(" \n");
      delay(delayMillis);
      return;
    } else {
      retryCount = 0;
    }
    if (pumpOn) {
      pumpOn = false;
      Serial.print("pump off\n distance: ");
      Serial.print(distance);
      Serial.print(" cm \n");
    }
  }

  delay(delayMillis);
}