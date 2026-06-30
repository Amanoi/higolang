FROM nginx:alpine
# Pre-built dist/ is copied directly — no Node build needed
COPY dist /usr/share/nginx/html
# Copy nginx config templates
COPY nginx.conf.template /etc/nginx/templates/nginx.conf.template
COPY nginx-http.conf.template /etc/nginx/templates/nginx-http.conf.template
# Copy entrypoint script
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
# Create certbot webroot directory
RUN mkdir -p /var/www/certbot
EXPOSE 80 443
ENTRYPOINT ["/docker-entrypoint.sh"]
