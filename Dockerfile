FROM debian:latest
RUN mkdir /app

EXPOSE 3031

COPY nsi_survey_server /app/server

ENTRYPOINT [ "/app/server" ]
