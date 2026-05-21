FROM nginx:1.27-alpine

RUN apk add --no-cache gettext

COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

COPY index.html styles.css app.js /usr/share/nginx/html/

EXPOSE 80
ENTRYPOINT ["/docker-entrypoint.sh"]
