String aString = "My name is particle";
double aDouble = 3.14;
int anInt = 0b1010011010;

int brewCoffee(String command);

void setup() {
    Particle.variable("aString", aString);
    Particle.variable("aDouble", aDouble);
    Particle.variable("anInt", anInt);

    Particle.function("brew", brewCoffee);
}

void loop() {
    delay(1000);
}

int brewCoffee(String command)
{
    // look for the matching argument "coffee" <-- max of 64 characters long
    if (command == "coffee") {
        return 1;
    }
    else return -1;
}