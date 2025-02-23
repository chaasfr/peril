FROM rabbitmq:3.13-management
RUN rabbitmq-plugins enable rabbitmq_stomp

COPY server /usr/local/bin/server
RUN chmod +x /usr/local/bin/server
CMD ["/bin/bash", "-c", "/usr/local/bin/server"]