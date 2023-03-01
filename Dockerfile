FROM alpine
ENV HOME=/
COPY /dist/scan_compare-linux-amd64 /bin/scan_compare
RUN adduser -D app
USER app
ENTRYPOINT ["/bin/scan_compare"]