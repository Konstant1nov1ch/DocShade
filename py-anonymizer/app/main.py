from fastapi import FastAPI
import uvicorn
from app.api.routes import router
from app.core.config import load_config
from app.utils.logger import setup_logger

# Загрузка конфигурации
config = load_config()

# Настройка логирования
logger = setup_logger("py-anonymizer")
logger.setLevel(config.get("log_level", "DEBUG"))

app = FastAPI(
    title="Document Anonymization Service",
    description="API for anonymizing PDF documents",
    version="1.0.0"
)

app.include_router(router)

@app.on_event("startup")
async def startup_event():
    # Инициализация действий при старте сервиса, если необходимо
    logger.info("Starting up the ML anonymizer service")

@app.on_event("shutdown")
async def shutdown_event():
    # Действия при завершении работы сервиса, если необходимо
    logger.info("Shutting down the ML anonymizer service")

if __name__ == "__main__":
    uvicorn.run(app, host=config["server"]["host"], port=config["server"]["port"], log_level="debug")