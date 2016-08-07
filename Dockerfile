FROM scratch
EXPOSE 1337
ADD stringservice /
CMD ["/stringservice"]
