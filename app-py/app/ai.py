import os
from dotenv import load_dotenv
from intents import Intent

load_dotenv()

AI_PROVIDER = os.getenv("AI_PROVIDER", "gemini").lower()

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
    if AI_PROVIDER == "gemini":
        return classify_intent_gemini(user_message)
    return Intent.UNKNOWN

def classify_intent_gemini(user_message: str) -> Intent:
    from google import genai

    client = genai.Client(api_key=os.getenv("GEMINI_API_KEY"))
    prompt = _build_prompt(user_message)

    response = client.models.generate_content(
        model="gemini-2.5-flash",
        contents=prompt,
    )

    intent_str = (response.text or "").strip().lower()
    return Intent(intent_str) if intent_str in Intent._value2member_map_ else Intent.UNKNOWN
