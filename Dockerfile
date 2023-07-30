# Сборка сервера приложения
FROM golang AS appbuilder
ADD . /home/
RUN cd /home/ && go build -o /app ./cmd/main.go 

# Сборка контейнера для приложения
FROM ubuntu:latest
WORKDIR /launch/

# Данные из контейнера с приложением
COPY --from=appbuilder /app /launch/app

RUN ls /launch
EXPOSE 8080
ENTRYPOINT ["/launch/app"]