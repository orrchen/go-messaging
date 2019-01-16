FROM debian
MAINTAINER "Orr Chen"
WORKDIR /app
# Now just add the binary
ADD app/tcp-server.linux /app/
ADD config /app/
EXPOSE 8080 3000
CMD ["./tcp-server.linux","-config=config.yml"]
