import os
from dotenv import load_dotenv
from intents import Intent

load_dotenv()

AI_PROVIDER = os.getenv("AI_PROVIDER", "openai").lower()

def _build_prompt(user_message: str) -> str:
    return f"""
Você é um classificador de intenção para um chatbot de atendimento de um Box de CrossFit.
Retorne APENAS uma das intenções abaixo (sem texto adicional):

- greeting
- about
- schedule
- plans
- book_class
- unknown

Mensagem do usuário: "{user_message}"
"""

def classify_intent(user_message: str) -> Intent:
    if AI_PROVIDER == "openai":
        return classify_intent_openai(user_message)
    elif AI_PROVIDER == "gemini":
        return classify_intent_gemini(user_message)
    else:
        return Intent.UNKNOWN

def classify_intent_openai(user_message: str) -> Intent:
    from openai import OpenAI

    client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

    prompt = _build_prompt(user_message)

    response = client.chat.completions.create(
        model="gpt-4.1-mini",
        messages=[
            {"role": "system", "content": "Responda somente com a intenção."},
            {"role": "user", "content": prompt},
        ],
        temperature=0,
    )

    intent_str = response.choices[0].message.content.strip().lower()
    return Intent(intent_str) if intent_str in Intent._value2member_map_ else Intent.UNKNOWN

def classify_intent_gemini(user_message: str) -> Intent:
    from google import genai
    import os

    client = genai.Client(api_key=os.getenv("GEMINI_API_KEY")) # pega a chave dentro do env para usar

    prompt = _build_prompt(user_message) # pegando o prompt para passar para a IA

    response = client.models.generate_content(
        model="gemini-2.5-flash",
        contents=prompt,
    )

    intent_str = (response.text or "").strip().lower()
    return Intent(intent_str) if intent_str in Intent._value2member_map_ else Intent.UNKNOWN