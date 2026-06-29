FROM node:20-alpine AS build
WORKDIR /app
ENV NODE_OPTIONS=--max-old-space-size=256
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
# Copy built frontend
COPY --from=build /app/dist /usr/share/nginx/html
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
