FROM golang

# Build the Go application
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
RUN apt-get install -y nginx vim make build-essential git php php-fpm
RUN git clone https://github.com/letsencrypt/letsencrypt /opt/letsencrypt
COPY vhost-rrfc /etc/nginx/sites-available
RUN mkdir /var/www/html/rrfc
RUN mkdir /var/www/html/rrfc/list
COPY index.php /var/www/html/rrfc
RUN ln -s /etc/nginx/sites-available/vhost-rrfc /etc/nginx/sites-enabled/rrfc
RUN rm /etc/nginx/sites-enabled/default
RUN service nginx start
WORKDIR /opt/letsencrypt
RUN ./certbot-auto --authenticator standalone --installer nginx -d 'rrfc.app' --email skarlso777@gmail.com --agree-tos -n --no-verify-ssl --redirect
WORKDIR /var/www/html/rrfc
# Copy over necessary files
COPY --from=0 /go/src/github.com/Skarlso/rrfc/rrfc .
RUN chmod +x rrfc

EXPOSE 80
