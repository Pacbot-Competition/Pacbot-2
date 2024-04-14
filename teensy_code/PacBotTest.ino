#include <Wire.h>
#include <VL6180X.h>

// Motor pins
const unsigned int MOTORDIR_PINS[4] = {5, 7, 23, 15};
const unsigned int MOTORPWM_PINS[4] = {6, 8, 22, 14};

VL6180X sensor1;
VL6180X sensor2;

// Sensor pins
const int sensor1_pin = 9;
const int sensor2_pin = 10;


#define LED_PIN LED_BUILTIN

void init_sensors() {
  pinMode(sensor1_pin, OUTPUT);
  pinMode(sensor2_pin, OUTPUT);
  digitalWrite(sensor1_pin, LOW);
  digitalWrite(sensor2_pin, LOW);

  Wire.begin();
  digitalWrite(sensor1_pin, HIGH);
  delay(50);
  sensor1.init();
  sensor1.configureDefault();
  sensor1.setTimeout(500);
  sensor1.setAddress(0x54);
  

  digitalWrite(sensor2_pin, HIGH);
  delay(50);
  sensor2.init();
  sensor2.configureDefault();
  sensor2.setTimeout(500);
  sensor2.setAddress(0x56);
}

void init_motors() {
  for (int i = 0; i<4; i++) {
    pinMode(MOTORDIR_PINS[i], OUTPUT);
    pinMode(MOTORPWM_PINS[i], OUTPUT);

    digitalWrite(MOTORDIR_PINS[i], LOW);
    analogWrite(MOTORPWM_PINS[i], 0);
  }
}

void setup()
{
  pinMode(LED_PIN, OUTPUT);
  init_sensors();
  init_motors();
  digitalWrite(LED_PIN, HIGH);
}

void loop()
{
  const int sensor1Value = sensor1.readRangeSingleMillimeters();
  const int sensor2Value = sensor2.readRangeSingleMillimeters();

  for (int i = 0; i<4; i++) {
    digitalWrite(MOTORDIR_PINS[i], LOW);
    analogWrite(MOTORPWM_PINS[i], 500);

    delay(1000);

    digitalWrite(MOTORDIR_PINS[i], HIGH); // Reverse

    delay(1000);

    analogWrite(MOTORPWM_PINS[i], 0);

    delay(500);
  }
}