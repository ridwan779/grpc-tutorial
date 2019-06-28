FROM ubuntu:16.04
RUN apt-get update
RUN apt-get install python3 -y
RUN apt-get install python3-pip -y
WORKDIR /app/server
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt
COPY . .
CMD ["python3", "-u", "app.py"]