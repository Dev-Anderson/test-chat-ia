FROM python:3.13-slim

WORKDIR /app

# Dependências do sistema (psycopg2 precisa de libpq)
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

# Dependências do python
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copia o código
COPY ./app ./app

EXPOSE 8000

# Em produção -> sem reload
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
