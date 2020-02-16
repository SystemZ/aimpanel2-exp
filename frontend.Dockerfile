FROM nginx:alpine
ADD aimpanel-master-frontend /usr/share/nginx/html
RUN ls -alh /usr/share/nginx/html
