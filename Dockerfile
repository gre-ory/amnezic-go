FROM debian:bullseye-slim

#==================
# update
#==================

RUN apt-get update

#==================
# ca-certificates
#==================

RUN apt-get install -y curl ca-certificates
RUN update-ca-certificates

#==================
# cxt-reporting-api
#==================

COPY ./bin /usr/local/bin/
COPY ./scripts/docker-entrypoint.sh /

RUN chmod 777 /docker-entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["amnezic"]
