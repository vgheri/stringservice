FROM gliderlabs/alpine
RUN apk-install bash
EXPOSE 1337
COPY stringservice /
