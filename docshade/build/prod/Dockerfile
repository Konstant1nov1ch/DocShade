# Stage 1: Build the React app
FROM node:16-alpine AS build

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и package-lock.json в рабочую директорию
COPY ./package*.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем все файлы и папки в рабочую директорию
COPY . .

# Устанавливаем аргумент сборки
ARG REACT_APP_BACKEND_HOST

# Определяем переменную окружения для использования в процессе сборки
ENV REACT_APP_BACKEND_HOST=$REACT_APP_BACKEND_HOST

# Собираем приложение для production
RUN npm run build

# Stage 2: Serve the React app with http-server
FROM node:16-alpine

WORKDIR /app

# Копируем собранное приложение из предыдущего этапа
COPY --from=build /app/build /app/build

# Устанавливаем http-server
RUN npm install -g http-server

# Запускаем сервер для обслуживания статических файлов
CMD ["http-server", "build"]
