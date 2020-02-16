FROM nginx:alpine
RUN rm /usr/share/nginx/html/index.html \
 && rm /usr/share/nginx/html/50x.html
ADD aimpanel-master-frontend /usr/share/nginx/html
RUN ls -alh /usr/share/nginx/html
