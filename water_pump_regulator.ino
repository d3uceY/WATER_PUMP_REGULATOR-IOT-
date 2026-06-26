#include <stdint.h>
#define TRIG_PIN 1
#define ECHO_PIN 3

int8_t threshold = 20;
bool pumpOn;


void setup() {
  Serial.begin(115200);
  pinMode(ECHO_PIN, INPUT);
  pinMode(TRIG_PIN, OUTPUT);
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
  long distance = readDistanceCM();

  if (distance <= 2 || distance > 400) {
    delay(100);
    return;
  }

  if (distance >= (20 + threshold)) {
    if (!pumpOn) {  
      pumpOn = true;
      Serial.print("pump on\n distance: ");
      Serial.print(distance);
      Serial.print(" cm \n");
    }
  }
  else if (distance < 20) {
    if (pumpOn) {   
      pumpOn = false;
      Serial.print("pump off\n distance: ");
      Serial.print(distance);
      Serial.print(" cm \n");
    }
  }

  delay(100);
}