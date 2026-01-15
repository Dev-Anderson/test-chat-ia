from fastapi import FastAPI
from pydantic import BaseModel
from intents import Intent
from ai import classify_intent

app = FastAPI(title="CrossFit Box Chatbot")

class ChatRequest(BaseModel):
    user_id: str
    message: str

def handle_intent(intent: Intent):
    if intent == Intent.GREETING:
        return {
            "reply": (
                "👋 Olá! Bem-vindo ao Box! 😄\n\n"
                "Como posso ajudar?\n"
                "1) Sobre o Box\n"
                "2) Horários\n"
                "3) Planos\n"
                "4) Reservar uma aula"
            )
        }

    if intent == Intent.ABOUT:
        return {
            "reply": (
                "🏋️ Sobre o Box:\n"
                "- Treinos para iniciantes e avançados\n"
                "- Coaches certificados\n"
                "- Aulas em grupo e acompanhamento\n\n"
                "Quer ver horários, planos ou reservar uma aula?"
            )
        }

    if intent == Intent.SCHEDULE:
        return {
            "reply": (
                "⏰ Horários:\n"
                "- Seg a Sex: 06:00, 07:00, 12:00, 18:00, 19:00, 20:00\n"
                "- Sábado: 09:00 e 10:00\n"
                "- Domingo: fechado\n\n"
                "Deseja reservar uma aula?"
            )
        }

    if intent == Intent.PLANS:
        return {
            "reply": (
                "💳 Planos:\n"
                "- 2x/semana: R$ 149/mês\n"
                "- 3x/semana: R$ 189/mês\n"
                "- Livre: R$ 229/mês\n"
                "- Day Pass: R$ 35\n\n"
                "Quer reservar uma aula experimental grátis?"
            )
        }

    if intent == Intent.BOOK_CLASS:
        return {
            "reply": (
                "✅ Reservar aula:\n"
                "Me diga o melhor dia e horário (ex: 'terça 19h')."
            )
        }

    return {
        "reply": (
            "Não entendi totalmente 😅\n"
            "Você quer:\n"
            "1) Sobre o Box\n"
            "2) Horários\n"
            "3) Planos\n"
            "4) Reservar uma aula"
        )
    }

@app.post("/chat")
def chat(req: ChatRequest):
    intent = classify_intent(req.message)
    return {
        "intent": intent,
        **handle_intent(intent)
    }
