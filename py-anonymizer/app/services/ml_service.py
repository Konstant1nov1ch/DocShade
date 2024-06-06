from .pdf_service import extract_text_from_pdf, create_pdf_from_text
from presidio_analyzer import AnalyzerEngine, PatternRecognizer, RecognizerResult, EntityRecognizer
from presidio_analyzer.nlp_engine import NlpEngineProvider, NlpArtifacts
from presidio_anonymizer import AnonymizerEngine
from presidio_anonymizer.entities import OperatorConfig, RecognizerResult, EngineResult
from natasha import Segmenter, NewsEmbedding, NewsMorphTagger, NewsSyntaxParser, NewsNERTagger, MorphVocab, Doc, NamesExtractor, DatesExtractor, MoneyExtractor, AddrExtractor
import spacy
import logging

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

# Инициализация spaCy
nlp_spacy = spacy.load("ru_core_news_md")

# Настройка компонентов Natasha
segmenter = Segmenter()
embeddings = NewsEmbedding()
morph_tagger = NewsMorphTagger(embeddings)
syntax_parser = NewsSyntaxParser(embeddings)
ner_tagger = NewsNERTagger(embeddings)
morph_vocab = MorphVocab()

# Настройка экстракторов Natasha
names_extractor = NamesExtractor(morph_vocab)
dates_extractor = DatesExtractor(morph_vocab)
money_extractor = MoneyExtractor(morph_vocab)
addr_extractor = AddrExtractor(morph_vocab)

# Конфигурация NLP-движка spaCy через NlpEngineProvider
configuration = {
    "nlp_engine_name": "spacy",
    "models": [{"lang_code": "ru", "model_name": "ru_core_news_md"}]
}

provider = NlpEngineProvider(nlp_configuration=configuration)
nlp_engine = provider.create_engine()

# Конфигурация AnalyzerEngine с использованием настроенного NLP-движка
analyzer = AnalyzerEngine(nlp_engine=nlp_engine, supported_languages=["ru"])
anonymizer = AnonymizerEngine()

# Создание кастомной функции для замены
def mask_text(entity_text):
    return '*' * (len(entity_text)//2)

# Настройка операторов для различных типов PII
operators = {
    "PERSON": OperatorConfig("custom", {"lambda": mask_text}),
    "LOCATION": OperatorConfig("custom", {"lambda": mask_text}),
    "ORGANIZATION": OperatorConfig("custom", {"lambda": mask_text}),
    "DATE_TIME": OperatorConfig("custom", {"lambda": mask_text}),
    "PHONE_NUMBER": OperatorConfig("custom", {"lambda": mask_text}),
    "EMAIL_ADDRESS": OperatorConfig("custom", {"lambda": mask_text})
}

def anonymize_text_presidio(text):
    # Анализ текста для обнаружения PII
    results = analyzer.analyze(text=text, entities=list(operators.keys()), language="ru")

    # Анонимизация текста на основе обнаруженных PII и операторов
    anonymized_result = anonymizer.anonymize(text=text, analyzer_results=results, operators=operators)
    
    return anonymized_result.text

class Span:
    def __init__(self, start, stop, entity_type):
        self.start = start
        self.stop = stop
        self.type = entity_type

def anonymize_text_natasha(text):
    doc = Doc(text)
    doc.segment(segmenter)
    doc.tag_ner(ner_tagger)

    spans = [Span(span.start, span.stop, span.type) for span in doc.spans if span.type in {"PER", "ORG", "DATE"}]

    # Использование дополнительных экстракторов Natasha
    spans.extend([Span(match.start, match.stop, "PER") for match in dates_extractor(text)])
    spans.extend([Span(match.start, match.stop, "DATE") for match in dates_extractor(text)])
    spans.extend([Span(match.start, match.stop, "MONEY") for match in money_extractor(text)])

    # Обработка текстов с помощью spaCy
    doc_spacy = nlp_spacy(text)
    for ent in doc_spacy.ents:
        if ent.label_.upper() in {"PER", "ORG", "DATE"}:
            spans.append(Span(ent.start_char, ent.end_char, ent.label_.upper()))

    spans = sorted(spans, key=lambda span: span.start)
    
    # Анонимизация текста
    anonymized_text = []
    current_pos = 0
    for span in spans:
        anonymized_text.append(text[current_pos:span.start])
        length = span.stop - span.start
        replacement = '*' * (length//2)
        anonymized_text.append(replacement)
        current_pos = span.stop
    
    anonymized_text.append(text[current_pos:])
    return ''.join(anonymized_text)

def anonymize_document(pdf_data: bytes) -> bytes:
    # Извлечение текста из PDF
    text = extract_text_from_pdf(pdf_data)
    
    # Анонимизация текста с использованием Presidio
    anonymized_text_presidio = anonymize_text_presidio(text)
    
    # Дополнительная анонимизация текста с использованием Natasha
    final_anonymized_text = anonymize_text_natasha(anonymized_text_presidio)
    
    # Создание нового PDF с анонимизированным текстом
    anonymized_pdf = create_pdf_from_text(final_anonymized_text)
    
    return anonymized_pdf
