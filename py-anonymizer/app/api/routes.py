from fastapi import APIRouter, UploadFile, File, HTTPException
from ..services.ml_service import anonymize_document
from fastapi.responses import StreamingResponse
import io

router = APIRouter()

@router.post("/anonymize", summary="Anonymize a PDF document", description="Uploads a PDF document and returns an anonymized version of it")
async def anonymize(file: UploadFile = File(...)):
    if file.content_type != "application/pdf":
        raise HTTPException(status_code=400, detail="Invalid file format. Only PDF is allowed.")
    
    content = await file.read()
    anonymized_pdf = anonymize_document(content)
    
    return StreamingResponse(io.BytesIO(anonymized_pdf), media_type="application/pdf", headers={"Content-Disposition": "attachment; filename=anonymized.pdf"})

@router.get("/health", summary="Health Check", description="Returns the health status of the service")
async def health():
    return {"status": "ok"}