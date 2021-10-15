# Alpine and 1.16 in go because is the "newer"
FROM golang:1.16-alpine AS builder

ARG MONGODB_URI
ARG DB
ARG COLLECTION
ENV MONGODB_URI $MONGODB_URI
ENV DB $DB
ENV COLLECTION $COLLECTION

WORKDIR /app

# Set the work destination


#copy the files
COPY app .


RUN go mod download
RUN go mod tidy

# creating the server object
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /server .

#Exposing the port
EXPOSE 443

CMD [ "/server" ]