#include <Wire.h>
#include <VL6180X.h>

// Motor pins
// current order for the motor pins
// 5, 7, 23, 15 - DIR
// 6, 8, 22, 14 - PWM
const unsigned int MOTORDIR_PINS[4] = {5, 7, 23, 15};
const unsigned int MOTORPWM_PINS[4] = {6, 8, 22, 14};
const unsigned int gpio_pin_1 = 2; 
const unsigned int gpio_pin_2 = 3; 
const unsigned int gpio_pin_3 = 4; 
const unsigned int X = 10;
VL6180X sensor1;
VL6180X sensor2;
VL6180X sensor3;
VL6180X sensor4;

// Sensor pins
const int sensor1_pin = 9;
const int sensor2_pin = 10;
// We arbitrarily chose pins 11 and 12 for the other two sensors
const int sensor3_pin = 11;
const int sensor4_pin = 12;


#define LED_PIN LED_BUILTIN

#define MAX_DISTANCE 10 /*max distance = sum of two sensor distances when not tilted*/

void init_sensors() {
  pinMode(sensor1_pin, OUTPUT);
  pinMode(sensor2_pin, OUTPUT);
  pinMode(sensor3_pin, OUTPUT);
  pinMode(sensor4_pin, OUTPUT);
  digitalWrite(sensor1_pin, LOW);
  digitalWrite(sensor2_pin, LOW);
  digitalWrite(sensor3_pin, LOW);
  digitalWrite(sensor4_pin, LOW);

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

  digitalWrite(sensor3_pin, HIGH);
  delay(50);
  sensor3.init();
  sensor3.configureDefault();
  sensor3.setTimeout(500);
  sensor3.setAddress(0x58);

  digitalWrite(sensor4_pin, HIGH);
  delay(50);
  sensor4.init();
  sensor4.configureDefault();
  sensor4.setTimeout(500);
  sensor4.setAddress(0x60);
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
  Serial.begin(9600);
  pinMode(LED_PIN, OUTPUT);
  pinMode(gpio_pin_1, INPUT);
  pinMode(gpio_pin_2, INPUT);
  pinMode(gpio_pin_3, INPUT);
  init_sensors();
  init_motors();
  digitalWrite(LED_PIN, HIGH);
}

int is_tilted(VL6180X *left_sensor, VL6180X *right_sensor){
  //int ldist = left_sensor.readRangeSingleMillimeters();
  //int rdist = right_sensor.readRangeSingleMillimeters();
  // Serial.println("is_tilted");
  if((left_sensor->readRangeSingleMillimeters() + right_sensor->readRangeSingleMillimeters()) > MAX_DISTANCE) { //is tilted
    return 1;
  } else { //not tilted
    return 0;
  }
}

void correct_tilt(int top_right, int bottom_right, int top_left, int bottom_left, VL6180X *left_sensor, VL6180X *right_sensor) {
  //if tilted to the right
  if (right_sensor->readRangeSingleMillimeters() < left_sensor->readRangeSingleMillimeters()) {
    //move top wheel to the left and bottom wheel to the right
    // digitalWrite(MOTORDIR_PINS[top], HIGH);
    // digitalWrite(MOTORDIR_PINS[bottom], LOW);
    analogWrite(MOTORPWM_PINS[top_right], 200);
    analogWrite(MOTORPWM_PINS[bottom_right], 200);
    analogWrite(MOTORPWM_PINS[top_left], 150);
    analogWrite(MOTORPWM_PINS[bottom_left], 150);
  }
  //if tilted to the left
  if (left_sensor->readRangeSingleMillimeters() > right_sensor->readRangeSingleMillimeters()){
    //move top wheel to the right and bottom wheel to the left
      // digitalWrite(MOTORDIR_PINS[top], LOW);
      // digitalWrite(MOTORDIR_PINS[bottom], HIGH);
    analogWrite(MOTORPWM_PINS[top_left], 200);
    analogWrite(MOTORPWM_PINS[bottom_left], 200);
    analogWrite(MOTORPWM_PINS[top_right], 150);
    analogWrite(MOTORPWM_PINS[bottom_right], 150);
  }

  analogWrite(MOTORPWM_PINS[top_left], 200);
  analogWrite(MOTORPWM_PINS[bottom_left], 200);
  analogWrite(MOTORPWM_PINS[top_right], 200);
  analogWrite(MOTORPWM_PINS[bottom_right], 200);
}

void correct_drift(int top_right, int bottom_right, int top_left, int bottom_left, VL6180X *left_sensor, VL6180X *right_sensor) {
  // Serial.println("correct_drift");
  // if(right_sensor->readRangeSingleMillimeters() == left_sensor->readRangeSingleMillimeters()) {
  //   analogWrite(MOTORPWM_PINS[bottom], 0);
  //   analogWrite(MOTORPWM_PINS[top], 0);
  // } 
  if (right_sensor->readRangeSingleMillimeters() > left_sensor->readRangeSingleMillimeters()) {
    // digitalWrite(MOTORDIR_PINS[top], HIGH); //moves top motor to the left
    // analogWrite(MOTORPWM_PINS[top], 200); //moves motors
    // digitalWrite(MOTORDIR_PINS[bottom], HIGH);
    // analogWrite(MOTORPWM_PINS[bottom], 200);
    analogWrite(MOTORPWM_PINS[top_left], 200);
    analogWrite(MOTORPWM_PINS[bottom_left], 150);
    analogWrite(MOTORPWM_PINS[top_right], 150);
    analogWrite(MOTORPWM_PINS[bottom_right], 200);
  }
  if (left_sensor->readRangeSingleMillimeters() > right_sensor->readRangeSingleMillimeters()) {
    //move it to the right
    // digitalWrite(MOTORDIR_PINS[top], LOW);
    // analogWrite(MOTORPWM_PINS[top], 200);
    // digitalWrite(MOTORDIR_PINS[bottom], LOW);
    // analogWrite(MOTORPWM_PINS[bottom], 200);
    analogWrite(MOTORPWM_PINS[top_left], 150);
    analogWrite(MOTORPWM_PINS[bottom_left], 200);
    analogWrite(MOTORPWM_PINS[top_right], 200);
    analogWrite(MOTORPWM_PINS[bottom_right], 150);
  }
  analogWrite(MOTORPWM_PINS[top_left], 200);
  analogWrite(MOTORPWM_PINS[bottom_left], 200);
  analogWrite(MOTORPWM_PINS[top_right], 200);
  analogWrite(MOTORPWM_PINS[bottom_right], 200);
}

// current order of sensors: top, bottom, left, right
// but this is moving two adjacent motors
void straight_old(int top, int bottom, int left, int right, VL6180X *left_sensor, VL6180X *right_sensor, int absolute_direction) {
  //go forward
  //if direction is forward or right
  // Serial.println("straight");
  if(absolute_direction == 2 || absolute_direction == 3){
    // testing top and bottom
    digitalWrite(MOTORDIR_PINS[left], LOW);
    digitalWrite(MOTORDIR_PINS[right], LOW);
    analogWrite(MOTORPWM_PINS[left], 200);
    analogWrite(MOTORPWM_PINS[right], 200);
  } else {
    digitalWrite(MOTORDIR_PINS[left], HIGH);
    digitalWrite(MOTORDIR_PINS[right], HIGH);
    analogWrite(MOTORPWM_PINS[left], 200);
    analogWrite(MOTORPWM_PINS[right], 200);
  }


  
}

void corner_wall(int left, int right, VL6180X *top_sensor, int absolute_direction)
{
  //Serial.println("corner_wall");
  //this function moves the robot forward in a corner case
  //0.5 in = 12.7 mm
  if(absolute_direction == 2 || absolute_direction == 3){
    while(top_sensor->readRangeSingleMillimeters() >= 12.7) {
    digitalWrite(MOTORDIR_PINS[left], LOW);
    digitalWrite(MOTORDIR_PINS[right], LOW);
    analogWrite(MOTORPWM_PINS[left], 200);
    analogWrite(MOTORPWM_PINS[right], 200);
    }
  }else {
    while(top_sensor->readRangeSingleMillimeters() >= 12.7) {
    digitalWrite(MOTORDIR_PINS[left], HIGH);
    digitalWrite(MOTORDIR_PINS[right], HIGH);
    analogWrite(MOTORPWM_PINS[left], 200);
    analogWrite(MOTORPWM_PINS[right], 200);
    }
  }
}

// this is relative forward
void straight(int top_left, int top_right, int bottom_left, int bottom_right, VL6180X *left_sensor, VL6180X *right_sensor, int absolute_direction){
  // 
  int top_left_dir = 0;
  int top_right_dir = 0;
  int bottom_left_dir = 0;
  int bottom_right_dir = 0;
  // 0 = not moving, 1 = left, 2 = right, 3 = forward, 4 = backward
  // logic for CW vs CCW for the motors based on the absolute direction
  if (absolute_direction == 1) { // left
    top_left_dir = 1;
    top_right_dir = 1;
  } else if (absolute_direction == 2) { // right
    bottom_left_dir = 1;
    bottom_right_dir = 1;
  } else if (absolute_direction == 3) { // forward 
    top_right_dir = 1;
    bottom_right_dir = 1;
  } else if (absolute_direction == 4) { // backward
    top_left_dir = 1;
    bottom_left_dir = 1;    
  }
  // direction of motors
  digitalWrite(MOTORDIR_PINS[top_left], top_left_dir);
  digitalWrite(MOTORDIR_PINS[top_right], top_right_dir);
  digitalWrite(MOTORDIR_PINS[bottom_left], bottom_left_dir);
  digitalWrite(MOTORDIR_PINS[bottom_right], bottom_right_dir);
  // speed of motors
  analogWrite(MOTORPWM_PINS[top_left], 200);
  analogWrite(MOTORPWM_PINS[top_right], 200);
  analogWrite(MOTORPWM_PINS[bottom_left], 200);
  analogWrite(MOTORPWM_PINS[bottom_right], 200);

  if(is_tilted(left_sensor, right_sensor) == 1) { /*is tilted*/
    correct_tilt(top_right, bottom_right, top_left, bottom_left, left_sensor, right_sensor);
  } else {
    correct_drift(top_right, bottom_right, top_left, bottom_left, left_sensor, right_sensor);
  }

}

void loop()
{
  digitalWrite(MOTORDIR_PINS[0], HIGH); 
  analogWrite(MOTORPWM_PINS[0], 100);

  // const int sensor1Value = sensor1.readRangeSingleMillimeters();
  // const int sensor2Value = sensor2.readRangeSingleMillimeters();
  // const int sensor3Value = sensor3.readRangeSingleMillimeters();
  // const int sensor4Value = sensor4.readRangeSingleMillimeters();

  // 3 bits - 8 possibilities 
  // bit1 bit2 bit3
  // 000 not move (0)
  // 001 left (1)
  // 010 right (2)
  // 011 forward (3)
  // 100 back (4)
  int gpio1_val = digitalRead(gpio_pin_1); 
  int gpio2_val = digitalRead(gpio_pin_2);
  int gpio3_val = digitalRead(gpio_pin_3);
  int bot_direction = 4*gpio1_val + 2*gpio2_val + 1*gpio3_val;

  VL6180X top_sensor;
  VL6180X bottom_sensor;
  VL6180X left_sensor;
  VL6180X right_sensor;

  // indices of motors
  int top_left;
  int top_right;
  int bottom_left;
  int bottom_right;
  
  bot_direction = 2;
  if (bot_direction == 0) {  //not moving
    for (int i = 0; i<4; i++) {
      digitalWrite(MOTORDIR_PINS[i], LOW); 
      analogWrite(MOTORPWM_PINS[i], 0);
    }
  } else if (bot_direction == 1) {  //going left
      top_sensor = sensor4;
      bottom_sensor = sensor2;
      left_sensor = sensor3;
      right_sensor = sensor1;

      top_left = 0;
      top_right = 3;
      bottom_left = 2;
      bottom_right = 1;
    // sensor 2 is at the top, sensor 1 and 3 are sides, sensor 4 is that back
  } else if (bot_direction == 2) { // going right
      top_sensor = sensor2;
      bottom_sensor = sensor4;
      left_sensor = sensor1;
      right_sensor = sensor3;

      top_left = 3;
      top_right = 0;
      bottom_left = 1;
      bottom_right = 2;
    //sensor 4 is at the top, sensor 1 and 3 is at the sides, sensor 2 is at the back
  } else if (bot_direction == 3) { //move forward
      top_sensor = sensor1;
      bottom_sensor = sensor3;
      left_sensor = sensor4;
      right_sensor = sensor2;

      top_left = 1; //top
      top_right = 2; //bottom
      bottom_left = 0; //left
      bottom_right = 3; //right
    // sensor1 is top, sensor 2 and 4 is at the sides, sensor 3 is at the bottom
  } else { // move backward
      top_sensor = sensor3;
      bottom_sensor = sensor1;
      left_sensor = sensor2;
      right_sensor = sensor4;

      top_left = 2;
      top_right = 0;
      bottom_left = 1;
      bottom_right = 3;
  }
  // IF THE top sensor reading is less than X (value to be determined later), don't go straight
  // OR 
  /*
  VL6180X old_top_sensor = top_sensor;
  VL6180X old_right_sensor = right_sensor;
  VL6180X old_bottom_sensor = bottom_sensor;


  if (left_sensor.readRangeSingleMillimeters() < X && right_sensor.readRangeSingleMillimeters() > X && top_sensor.readRangeSingleMillimeters() < X) { // 
      // corner, turn left
      corner_wall(left_motor, right_motor, &top_sensor, bot_direction);
      top_sensor = left_sensor; //L
      right_sensor =  old_top_sensor;
      bottom_sensor = old_right_sensor; //R 
      left_sensor = old_bottom_sensor; // B   
      top_motor = 3;
      bottom_motor = 1;
      left_motor = 2;
      right_motor = 0;
  } else if (right_sensor.readRangeSingleMillimeters() < X && left_sensor.readRangeSingleMillimeters() > X && top_sensor.readRangeSingleMillimeters() < X) {
      // corner, turn right
      corner_wall(left_motor, right_motor, &top_sensor, bot_direction);
      bottom_sensor = left_sensor; //L 
      left_sensor = old_top_sensor;
      right_sensor = old_bottom_sensor; //B
      top_sensor = old_right_sensor; //R
      top_motor = 1;
      bottom_motor = 3;
      left_motor = 0;
      right_motor = 2;
  } 
  */
  straight(top_left, top_right, bottom_left, bottom_right, left_sensor, right_sensor, bot_direction);
  

// for (int i = 0; i<4; i++) {
//   digitalWrite(MOTORDIR_PINS[i], LOW); // direction of motor
//   analogWrite(MOTORPWM_PINS[i], 200); // speed of motor - 0 to 255 speed

  //   delay(1000);

  //   digitalWrite(MOTORDIR_PINS[i], HIGH); // Reverse 

  //   delay(1000);

  //   analogWrite(MOTORPWM_PINS[i], 0);

  //   delay(500);
  // }
}