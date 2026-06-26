#include <stdint.h>
#define TRIG_PIN 5
#define ECHO_PIN 6

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

  Serial.print("Distance: ");
  Serial.print(distance);
  Serial.println(" cm");

  delay(100);
}