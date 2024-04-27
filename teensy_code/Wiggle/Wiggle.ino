#include <Wire.h>

const unsigned int MOTOR_FORWARD_PINS[4] = {5, 7, 23, 14};
const unsigned int MOTOR_BACKWARD_PINS[4] = {6, 8, 22, 15};


const unsigned int gpio_pin_1 = 2; 
const unsigned int gpio_pin_2 = 3; 
const unsigned int gpio_pin_3 = 4; 

int speed = 100;
int previous_direction = 0;
int delay_time = 200; 
int vertical_time = 2000;
int horizontal_time = 500;
int speed_multi
double WIGGLE_CONST = (double)40;

#define LED_PIN LED_BUILTIN

void init_motors() {
  for (int i = 0; i<4; i++) {
    pinMode(MOTOR_FORWARD_PINS[i], OUTPUT);
    pinMode(MOTOR_BACKWARD_PINS[i], OUTPUT);

    analogWrite(MOTOR_FORWARD_PINS[i],rng_babyy());
    analogWrite(MOTOR_BACKWARD_PINS[i],rng_babyy());
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
  analogWrite(MOTOR_FORWARD_PINS[bottom_right],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[top_right],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[top_left],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[top_left], speed*speed_multiplier);
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], speed*speed_multiplier);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], speed);
}

void backward(int top_left, int top_right, int bottom_left, int bottom_right){
  analogWrite(MOTOR_FORWARD_PINS[top_left],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[bottom_left],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[top_right],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[top_right], speed);
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], speed);
  analogWrite(MOTOR_BACKWARD_PINS[top_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], speed);
  
}

void left(int top_left, int top_right, int bottom_left, int bottom_right) {
  analogWrite(MOTOR_FORWARD_PINS[top_left],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[top_right],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[bottom_left], speed);
  analogWrite(MOTOR_FORWARD_PINS[bottom_right], speed);

  analogWrite(MOTOR_BACKWARD_PINS[top_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[top_right], speed);
  
}

void right(int top_left, int top_right, int bottom_left, int bottom_right) {
  analogWrite(MOTOR_FORWARD_PINS[bottom_left],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[bottom_right],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[top_left],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[top_right],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[top_left], speed);
  analogWrite(MOTOR_FORWARD_PINS[top_right], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left], speed);
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right], speed);
}

void stop(int top_left, int top_right, int bottom_left, int bottom_right){
  analogWrite(MOTOR_FORWARD_PINS[top_left],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[top_right],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[bottom_left],rng_babyy());
  analogWrite(MOTOR_FORWARD_PINS[bottom_right],rng_babyy());

  analogWrite(MOTOR_BACKWARD_PINS[top_left],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[top_right],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[bottom_left],rng_babyy());
  analogWrite(MOTOR_BACKWARD_PINS[bottom_right],rng_babyy());
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
  Serial.print("Direction:");
  Serial.println(bot_direction);

  if (bot_direction == 1) { // left
      left(top_left, top_right, bottom_left, bottom_right);
    } else if (bot_direction == 2) { // right
      right(top_left, top_right, bottom_left, bottom_right);
    } else if (bot_direction == 3) { // forward 
      forward(top_left, top_right, bottom_left, bottom_right);
    } else if (bot_direction == 4) { // backward
      backward(top_left, top_right, bottom_left, bottom_right);    
    } else if (bot_direction == 5) { // chaos behavior
      int random_bot_direction = random(1, 5);
      if (bot_direction == 1) { // left
        left(top_left, top_right, bottom_left, bottom_right);
        delay(horizontal_time);
      } else if (bot_direction == 2) { // right
        right(top_left, top_right, bottom_left, bottom_right);
        delay(horizontal_time);
      } else if (bot_direction == 3) { // forward 
        forward(top_left, top_right, bottom_left, bottom_right);
        delay(vertical_time);
      } else if (bot_direction == 4) { // backward
        backward(top_left, top_right, bottom_left, bottom_right);  
        delay(vertical_time);  
      }
    } else {
      stop(top_left, top_right, bottom_left, bottom_right);
    }
}

double rng_babyy()
{
    return (double)rand() / (double)RAND_MAX * WIGGLE_CONST;
}
