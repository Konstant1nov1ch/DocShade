FROM python:3.9-slim

# Установка системных зависимостей для компиляции
RUN apt-get update && apt-get install -y \
    build-essential \
    gcc \
    libxml2-dev \
    libxslt1-dev \
    zlib1g-dev \
    libffi-dev \
    libssl-dev \
    libjpeg-dev \
    libfreetype6-dev \
    curl \
    poppler-utils \
    && rm -rf /var/lib/apt/lists/*

# Обновление pip, setuptools и wheel до последних версий
RUN pip install --upgrade pip setuptools wheel

# Установка зависимостей Python
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install python-multipart

# Установка модели spaCy
RUN python -m spacy download ru_core_news_md

WORKDIR /application

# Копирование файлов приложения
COPY ./app /application/app

# Копирование конфигурационного файла
COPY ./app/config.json /application/app/config.json


# Установка переменной окружения для пути к конфигурационному файлу
ENV CONFIG_PATH=./app/config.json

ENV PYTHONPATH /application/app:/application

# Запуск приложения
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
