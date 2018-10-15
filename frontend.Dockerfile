FROM nginx:alpine
RUN rm -r /usr/share/nginx/html/
ADD aimpanel-master-frontend /usr/share/nginx/html
RUN ls -alh /usr/share/nginx/html