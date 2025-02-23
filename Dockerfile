FROM rabbitmq:3.13-management
RUN rabbitmq-plugins enable rabbitmq_stomp

EXPOSE 5672 15672
COPY rabbit.sh /usr/local/bin/rabbit.sh
COPY server /usr/local/bin/server
RUN chmod +x /usr/local/bin/rabbit.sh /usr/local/bin/server
CMD ["/bin/bash", "-c", "/usr/local/bin/rabbit.sh start & /usr/local/bin/server"]