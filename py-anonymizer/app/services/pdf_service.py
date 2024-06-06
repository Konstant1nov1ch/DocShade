import os
import tempfile
from pdfminer.high_level import extract_text
from reportlab.lib.pagesizes import letter
from reportlab.pdfgen import canvas
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.pdfbase import pdfmetrics
import io

def extract_text_from_pdf(pdf_data: bytes) -> str:
    with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as temp_pdf:
        temp_pdf.write(pdf_data)
        temp_pdf.flush()
        text = extract_text(temp_pdf.name)
    os.remove(temp_pdf.name)
    return text

def split_long_lines(lines, max_length=100):
    """Splits long lines into multiple lines with a maximum length of max_length."""
    new_lines = []
    for line in lines:
        while len(line) > max_length:
            split_pos = line.rfind(' ', 0, max_length)
            if split_pos == -1:
                split_pos = max_length
            new_lines.append(line[:split_pos].strip())
            line = line[split_pos:].strip()
        new_lines.append(line)
    return new_lines

def create_pdf_from_text(anonymized_text: str) -> bytes:
    output_directory = tempfile.mkdtemp()
    output_file = os.path.join(output_directory, "output.pdf")

    # Разбиение текста на строки
    lines = anonymized_text.splitlines()

    # Разделение длинных строк на более короткие
    lines = split_long_lines(lines)

    # Создание нового PDF файла
    c = canvas.Canvas(output_file, pagesize=letter)

    # Загрузка и регистрация шрифта
    pdfmetrics.registerFont(TTFont('DejaVuSans', 'DejaVuSans.ttf'))

    # Установка шрифта
    c.setFont("DejaVuSans", 10)

    i = 750
    line_number = 0

    while line_number < len(lines):
        if len(lines) - line_number >= 60:  # Если осталось 60 и более строк
            for line in lines[line_number:line_number + 60]:
                c.drawString(15, i, line.strip())
                line_number += 1
                i -= 12
            c.showPage()
            c.setFont("DejaVuSans", 10)
            i = 750
        else:  # Если осталось меньше 60 строк
            for line in lines[line_number:]:
                c.drawString(15, i, line.strip())
                line_number += 1
                i -= 12
            c.showPage()

    # Сохранение PDF файла
    c.save()

    # Чтение созданного PDF в байты
    with open(output_file, "rb") as f:
        pdf_bytes = f.read()

    # Удаление временного PDF файла и директории
    os.remove(output_file)
    os.rmdir(output_directory)

    return pdf_bytes