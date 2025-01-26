#include <Wire.h>
#include <VL6180X.h>

// Motor pins
const unsigned int MOTORCW_PINS[4]  = {8, 14, 23, 6};
const unsigned int MOTORCCW_PINS[4] = {7, 15, 22, 5};
const unsigned int gpio_pin_1 = 2; 
const unsigned int gpio_pin_2 = 3; 
const unsigned int gpio_pin_3 = 4; 

#define CW(n, speed) analogWrite(MOTORCW_PINS[n], speed)
#define CCW(n, speed) analogWrite(MOTORCCW_PINS[n], speed)

VL6180X sensor_2_3;
VL6180X sensor_0_2;
VL6180X sensor_0_1;
VL6180X sensor_3_1;

// Sensor pins
const int sensor_2_3_pin = 9;
const int sensor_0_2_pin = 10;
const int sensor_0_1_pin = 11;
const int sensor_3_1_pin = 12;


#define LED_PIN LED_BUILTIN

#define MAX_DISTANCE 10 /*max distance = sum of two sensor distances when not tilted*/

void init_sensors() {
  pinMode(sensor_2_3_pin, OUTPUT);
  pinMode(sensor_0_2_pin, OUTPUT);
  pinMode(sensor_0_1_pin, OUTPUT);
  pinMode(sensor_3_1_pin, OUTPUT);
  digitalWrite(sensor_2_3_pin, LOW);
  digitalWrite(sensor_0_2_pin, LOW);
  digitalWrite(sensor_0_1_pin, LOW);
  digitalWrite(sensor_3_1_pin, LOW);

  Wire.begin();
  digitalWrite(sensor_2_3_pin, HIGH);
  delay(50);
  sensor_2_3.init();
  sensor_2_3.configureDefault();
  sensor_2_3.setTimeout(500);
  sensor_2_3.setAddress(0x54);
  

  digitalWrite(sensor_0_2_pin, HIGH);
  delay(50);
  sensor_0_2.init();
  sensor_0_2.configureDefault();
  sensor_0_2.setTimeout(500);
  sensor_0_2.setAddress(0x56);

  digitalWrite(sensor_0_1_pin, HIGH);
  delay(50);
  sensor_0_1.init();
  sensor_0_1.configureDefault();
  sensor_0_1.setTimeout(500);
  sensor_0_1.setAddress(0x58);

  digitalWrite(sensor_3_1_pin, HIGH);
  delay(50);
  sensor_3_1.init();
  sensor_3_1.configureDefault();
  sensor_3_1.setTimeout(500);
  sensor_3_1.setAddress(0x60);
}

void init_motors() {
  for (int i = 0; i<4; i++) {
    pinMode(MOTORCW_PINS[i], OUTPUT);
    pinMode(MOTORCCW_PINS[i], OUTPUT);

    analogWrite(MOTORCW_PINS[i], 0);
    analogWrite(MOTORCCW_PINS[i], 0);
  }
}

void forward(int speed, int rightBias) {
    CW(0, 0);
    CW(1, 0);
    CW(2, speed);
    CW(3, speed);
    CCW(0, speed);
    CCW(1, speed);
    CCW(2, 0);
    CCW(3, 0);
}

void backward(int speed) {
    CW(0, speed);
    CW(1, speed);
    CW(2, 0);
    CW(3, 0);
    CCW(0, 0);
    CCW(1, 0);
    CCW(2, speed);
    CCW(3, speed);
}

void right(int speed) {
  CW(0, speed);
  CW(1, 0);
  CW(2, speed);
  CW(3, 0);
  CCW(0, 0);
  CCW(1, speed);
  CCW(2, 0);
  CCW(3, speed);
}

void left(int speed) {
  CW(0, 0);
  CW(1, speed);
  CW(2, 0);
  CW(3, speed);
  CCW(0, speed);
  CCW(1, 0);
  CCW(2, speed);
  CCW(3, 0);
}

void stop() {
  for (int i = 0; i<4; i++) {
    CW(i, 0);
    CCW(i, 0);
  }
}

void setup()
{
  Serial.begin(9600);
  pinMode(LED_PIN, OUTPUT);
  pinMode(gpio_pin_1, INPUT);
  pinMode(gpio_pin_2, INPUT);
  pinMode(gpio_pin_3, INPUT);
  init_sensors();
  init_motors();
  digitalWrite(LED_PIN, HIGH);

  Serial.println("Setup complete");
}

void loop() {
  delay(3000);
  forward(120);
  delay(500);
  right(120);
  delay(500);
  backward(120);
  delay(500);
  left(120);
  delay(500);
  stop();
}

