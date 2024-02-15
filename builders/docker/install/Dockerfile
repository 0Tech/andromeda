FROM debian:bookworm AS platform

RUN apt-get update && apt-get upgrade -y && apt-get install -y cmake

FROM platform AS builders

WORKDIR /tmp/builders
COPY install .
RUN cmake -S . -B build -Wno-dev && cmake --build build --target install_builders

FROM platform AS final

ARG UID
ARG GID

RUN addgroup --system --gid ${GID} builder
RUN adduser --system --uid ${UID} --gid ${GID} builder --home /home/builder

COPY --from=builders /usr/local /usr/local/
COPY --from=builders --chown=${UID}:${GID} /root/go /home/builder/go/

ENV PATH /home/builder/go/bin:/usr/local/bin:/usr/bin

WORKDIR /tmp/mods
RUN chown ${UID}:${GID} .
COPY --chown=${UID}:${GID} mods .

USER ${UID}:${GID}
RUN go mod download
