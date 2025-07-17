FROM gcc:latest

WORKDIR /app

ENTRYPOINT ["/bin/bash", "-c", "g++ /app/main.cpp -o /app/main.out && /app/main.out < /app/input.txt"]
