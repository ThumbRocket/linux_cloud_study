FROM python:3.8-slim-buster

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

RUN apt-get update \
  && apt-get install -y gcc libpq-dev python-dev vim \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /usr/src/app
COPY requirements.txt ./
RUN pip install -r requirements.txt
COPY . .

EXPOSE 6000

#CMD ["python", "manage.py", "runserver", "0.0.0.0:6000"]
CMD ["gunicorn", "--bind", "0.0.0.0:6000", "--workers", "3", "pollme.wsgi:application"]