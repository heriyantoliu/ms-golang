FROM rabbitmq:3.8

RUN rabbitmq-plugins enable --offline rabbitmq_management
RUN rabbitmq-plugins enable --offline rabbitmq_mqtt

EXPOSE 15672
EXPOSE 5672