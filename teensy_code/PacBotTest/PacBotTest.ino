#include <Wire.h>
#include <VL6180X.h>

// Motor pins
const unsigned int MOTORDIR_PINS[4] = {5, 7, 23, 15};
const unsigned int MOTORPWM_PINS[4] = {6, 8, 22, 14};

VL6180X left_sensor;
VL6180X right_sensor;
VL6180X front_sensor;
VL6180X back_sensor;

// Sensor pins
const int left_sensor_pin = 9;
const int right_sensor_pin = 10;
// We arbitrarily chose pins 11 and 12 for the other two sensors
const int front_sensor_pin = 11;
const int back_sensor_pin = 12;


#define LED_PIN LED_BUILTIN

void init_sensors() {
  pinMode(left_sensor_pin, OUTPUT);
  pinMode(right_sensor_pin, OUTPUT);
  pinMode(front_sensor_pin, OUTPUT);
  pinMode(back_sensor_pin, OUTPUT);
  digitalWrite(left_sensor_pin, LOW);
  digitalWrite(right_sensor_pin, LOW);
  digitalWrite(front_sensor_pin, LOW);
  digitalWrite(back_sensor_pin, LOW);

  Wire.begin();
  digitalWrite(left_sensor_pin, HIGH);
  delay(50);
  left_sensor.init();
  left_sensor.configureDefault();
  left_sensor.setTimeout(500);
  left_sensor.setAddress(0x54);
  

  digitalWrite(right_sensor_pin, HIGH);
  delay(50);
  right_sensor.init();
  right_sensor.configureDefault();
  right_sensor.setTimeout(500);
  right_sensor.setAddress(0x56);

  digitalWrite(front_sensor_pin, HIGH);
  delay(50);
  front_sensor.init();
  front_sensor.configureDefault();
  front_sensor.setTimeout(500);
  front_sensor.setAddress(0x58);

  digitalWrite(back_sensor_pin, HIGH);
  delay(50);
  back_sensor.init();
  back_sensor.configureDefault();
  back_sensor.setTimeout(500);
  back_sensor.setAddress(0x60);
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
  const int leftsensorValue = left_sensor.readRangeSingleMillimeters();
  const int rightsensorValue = right_sensor.readRangeSingleMillimeters();
  const int frontsensorValue = front_sensor.readRangeSingleMillimeters();
  const int backsensorValue = back_sensor.readRangeSingleMillimeters();

  
  for (int i = 0; i<4; i++) {
    digitalWrite(MOTORDIR_PINS[i], LOW); // direction of motor
    analogWrite(MOTORPWM_PINS[i], 200); // speed of motor - 0 to 255 speed

    delay(1000);

    digitalWrite(MOTORDIR_PINS[i], HIGH); // Reverse 

    delay(1000);

    analogWrite(MOTORPWM_PINS[i], 0);

    delay(500);
  }
}