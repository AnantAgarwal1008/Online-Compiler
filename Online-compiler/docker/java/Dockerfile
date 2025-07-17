FROM openjdk:latest

WORKDIR /app

ENTRYPOINT ["/bin/bash", "-c", "javac /app/Main.java && java Main < /app/input.txt"]
