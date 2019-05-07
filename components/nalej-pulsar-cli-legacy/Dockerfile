###
## This image is going to be ignored until we have a solution for the cross-compilation problem.
##


FROM nalejdev.azurecr.io/juan/mypulsar:latest AS pulsar

FROM scratch AS base

COPY --from=pulsar /usr/local/lib/libpulsar.so /usr/local/lib/.
COPY --from=pulsar /usr/local/lib/libpulsar.a /usr/local/lib/.

# Now just add the binary
COPY nalej-pulsar-cli /nalej/
ENTRYPOINT ["./nalej/nalej-pulsar-cli"]


