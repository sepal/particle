String aString = "My name is particle";
double aDouble = 3.14;
int anInt = 0b1010011010;

void setup() {
    Particle.variable("aString", aString);
    Particle.variable("aDouble", aDouble);
    Particle.variable("anInt", anInt);
}

void loop() {
    delay(1000);
}