FROM nginx:1.27-alpine

COPY dev/ui/nginx.conf /etc/nginx/conf.d/default.conf
COPY dev/ui/index.html dev/ui/styles.css dev/ui/app.js /usr/share/nginx/html/

EXPOSE 80
