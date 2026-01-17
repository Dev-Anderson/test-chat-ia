from app.intents import Intent
from app.settings import AI_PROVIDER, GEMINI_API_KEY

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

    # (Se quiser depois eu adiciono OpenAI aqui também)
    return Intent.UNKNOWN

def classify_intent_gemini(user_message: str) -> Intent:
    if not GEMINI_API_KEY:
        return Intent.UNKNOWN

    from google import genai

    client = genai.Client(api_key=GEMINI_API_KEY)
    prompt = _build_prompt(user_message)

    response = client.models.generate_content(
        model="gemini-2.5-flash",
        contents=prompt,
    )

    intent_str = (response.text or "").strip().lower()
    return Intent(intent_str) if intent_str in Intent._value2member_map_ else Intent.UNKNOWN
