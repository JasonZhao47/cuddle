FROM ubuntu:20.04
COPY cuddle /app/cuddle
WORKDIR /app
CMD ["/app/cuddle"]