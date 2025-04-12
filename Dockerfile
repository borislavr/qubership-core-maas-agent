FROM ghcr.io/netcracker/qubership/core-base:main-20250411083149-49

COPY --chown=10001:0 --chmod=755 maas-agent-service/maas-agent-service /app/maas-agent
COPY --chown=10001:0 maas-agent-service/application.yaml /app/
COPY --chown=10001:0 maas-agent-service/policies.conf /app/

CMD ["/app/maas-agent"]