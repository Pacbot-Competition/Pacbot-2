#include <Wire.h>

const unsigned int MOTORDIR_PINS[4] = {5, 7, 23, 15};
const unsigned int MOTORPWM_PINS[4] = {6, 8, 22, 14};
const unsigned int gpio_pin_1 = 2; 
const unsigned int gpio_pin_2 = 3; 
const unsigned int gpio_pin_3 = 4; 
#define LED_PIN LED_BUILTIN


void init_motors() {
  for (int i = 0; i<4; i++) {
    pinMode(MOTORDIR_PINS[i], OUTPUT);
    pinMode(MOTORPWM_PINS[i], OUTPUT);

    digitalWrite(MOTORDIR_PINS[i], LOW);
    analogWrite(MOTORPWM_PINS[i], 0);
  }
}

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  pinMode(LED_PIN, OUTPUT);
  pinMode(gpio_pin_1, INPUT);
  pinMode(gpio_pin_2, INPUT);
  pinMode(gpio_pin_3, INPUT);
  init_motors();
  digitalWrite(LED_PIN, HIGH);
}

void straight(int top_left, int top_right, int bottom_left, int bottom_right, int absolute_direction){
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

}

void loop() {
  // put your main code here, to run repeatedly:
  int top_left = 0;
  int top_right = 1;
  int bottom_left = 2;
  int bottom_right = 3;

  int gpio1_val = digitalRead(gpio_pin_1); 
  int gpio2_val = digitalRead(gpio_pin_2);
  int gpio3_val = digitalRead(gpio_pin_3);
  int bot_direction = 4*gpio1_val + 2*gpio2_val + 1*gpio3_val;

  straight(top_left, top_right, bottom_left, bottom_right, bot_direction);
}
