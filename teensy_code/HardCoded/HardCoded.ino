#include <Wire.h>

const unsigned int MOTOR_FORWARD_PINS[4] = {5, 7, 23, 14};
const unsigned int MOTOR_BACKWARD_PINS[4] = {6, 8, 22, 15};


const unsigned int gpio_pin_1 = 2; 
const unsigned int gpio_pin_2 = 3; 
const unsigned int gpio_pin_3 = 4; 

int speed = 100;
int block_time = 400;
#define LED_PIN LED_BUILTIN


void init_motors() {
  for (int i = 0; i<4; i++) {
    pinMode(MOTOR_FORWARD_PINS[i], OUTPUT);
    pinMode(MOTOR_BACKWARD_PINS[i], OUTPUT);

    analogWrite(MOTOR_FORWARD_PINS[i], 0);
    analogWrite(MOTOR_BACKWARD_PINS[i], 0);
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

void forward(int top_left, int top_right, int bottom_left, int bottom_right) {
  // 
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], 0);
  analogWrite(MOTOR_FORWARD_PINS[top_right], 0);
  analogWrite(MOTOR_BACKWARD_PINS[top_left], 0);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], 0);
  analogWrite(MOTOR_FORWARD_PINS[top_left], speed);
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], speed);
}

void backward(int top_left, int top_right, int bottom_left, int bottom_right){
  analogWrite(MOTOR_FORWARD_PINS[top_left], 0);
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], 0);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], 0);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], 0);
  analogWrite(MOTOR_FORWARD_PINS[top_right], speed);
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], speed);
  analogWrite(MOTOR_BACKWARD_PINS[top_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], speed);
  
}

void left(int top_left, int top_right, int bottom_left, int bottom_right) {
  analogWrite(MOTOR_FORWARD_PINS[top_left], 0);
  analogWrite(MOTOR_FORWARD_PINS[top_right], 0);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], 0);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], 0);
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], speed);
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], speed);

  analogWrite(MOTOR_BACKWARD_PINS[top_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], speed);
  
}

void right(int top_left, int top_right, int bottom_left, int bottom_right) {
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], 0);
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], 0);
  analogWrite(MOTOR_BACKWARD_PINS[top_left], 0);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], 0);
  analogWrite(MOTOR_FORWARD_PINS[top_left], speed);
  analogWrite(MOTOR_FORWARD_PINS[top_right], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], speed);
}

void stop(int top_left, int top_right, int bottom_left, int bottom_right){
  analogWrite(MOTOR_FORWARD_PINS[top_left], 0);
  analogWrite(MOTOR_FORWARD_PINS[top_right], 0);
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], 0);
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], 0);

  analogWrite(MOTOR_BACKWARD_PINS[top_left], 0);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], 0);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], 0);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], 0);
}


void loop() {
  // put your main code here, to run repeatedly:
  // Serial.println("HELLO");
  // int i = 2;

  // analogWrite(MOTOR_BACKWARD_PINS[i], 0);
  // analogWrite(MOTOR_FORWARD_PINS[i], 200);
    
  // delay(500);
  // analogWrite(MOTOR_BACKWARD_PINS[i], 0);
  // analogWrite(MOTOR_FORWARD_PINS[i], 0);

  // delay(500);

  
  int top_left = 0;
  int top_right = 1;
  int bottom_left = 2;
  int bottom_right = 3;
  
  left(top_left, top_right, bottom_left, bottom_right);
  delay(7*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  backward(top_left, top_right, bottom_left, bottom_right);
  delay(3*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  left(top_left, top_right, bottom_left, bottom_right);
  delay(5*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  backward(top_left, top_right, bottom_left, bottom_right);
  delay(3*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  right(top_left, top_right, bottom_left, bottom_right);
  delay(25*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  forward(top_left, top_right, bottom_left, bottom_right);
  delay(3*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

left(top_left, top_right, bottom_left, bottom_right);
  delay(5*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  forward(top_left, top_right, bottom_left, bottom_right);
  delay(25*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

right(top_left, top_right, bottom_left, bottom_right);
  delay(5*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

backward(top_left, top_right, bottom_left, bottom_right);
  delay(4*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

left(top_left, top_right, bottom_left, bottom_right);
  delay(25*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  forward(top_left, top_right, bottom_left, bottom_right);
  delay(4*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  right(top_left, top_right, bottom_left, bottom_right);
  delay(5*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  backward(top_left, top_right, bottom_left, bottom_right);
  delay(19*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);

  right(top_left, top_right, bottom_left, bottom_right);
  delay(6*block_time);
  stop(top_left, top_right, bottom_left, bottom_right);


  delay(10000000);
}
