#-- Base Stage --
FROM debian:stable-slim AS base

ARG DOCKER_USER=default_user

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    gnupg  \
    golang && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd -r "${DOCKER_USER}" && useradd -r -g "${DOCKER_USER}" "${DOCKER_USER}"


#-- Build Stage --
FROM base AS digital-contracting-service

ARG DOCKER_USER=default_user

RUN mkdir -p /home/${DOCKER_USER}
RUN chmod 0777 /home/${DOCKER_USER}

RUN mkdir -p /app
RUN chown $DOCKER_USER:$DOCKER_USER /app

# Create directory for custom CA certificates
RUN mkdir -p /usr/local/share/ca-certificates/custom
RUN chmod 755 /usr/local/share/ca-certificates/custom

WORKDIR /app

USER $DOCKER_USER

COPY --chown=$DOCKER_USER:$DOCKER_USER ./ /app/

# Copy entrypoint script
COPY --chown=$DOCKER_USER:$DOCKER_USER docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

EXPOSE 8991

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["go",  "run" , "./cmd/dcs"]