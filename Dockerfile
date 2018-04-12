FROM 100.125.16.65:20202/hwcse/as-go:1.8.5

COPY ./test06 /home
COPY ./conf /home/conf
RUN chmod +x /home/test06

CMD ["/home/test06"]