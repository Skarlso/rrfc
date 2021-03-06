FROM golang

# Build the Go application
RUN echo "2.0" > /root/version
RUN go get -d -v github.com/Skarlso/rrfc/...
WORKDIR /go/src/github.com/Skarlso/rrfc
RUN go build

# Create the main image from postgres base
FROM ubuntu:18.04
LABEL Author="Gergely Brautigam"
# Install nginx and setup HTTPS
RUN apt-get update
RUN apt-get install tzdata
ENV TZ Europe/Budapest
RUN apt-get install -y nginx vim make build-essential git curl iputils-ping
COPY vhost-rrfc /etc/nginx/sites-available/
RUN mkdir /var/www/html/rrfc
RUN mkdir /var/www/html/rrfc/list
COPY site/ /var/www/html/rrfc/
RUN ln -s /etc/nginx/sites-available/vhost-rrfc /etc/nginx/sites-enabled/rrfc
RUN rm /etc/nginx/sites-enabled/default
WORKDIR /var/www/html/rrfc
# Copy over necessary files
COPY --from=0 /go/src/github.com/Skarlso/rrfc/rrfc .
COPY rfc.template .
COPY previous.template .
COPY index.template .
RUN chmod +x rrfc

EXPOSE 443

VOLUME [ "/etc/letsencrypt/" ]

COPY start.sh /root/
CMD ["/root/start.sh"]
