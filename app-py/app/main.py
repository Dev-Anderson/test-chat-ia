from fastapi import FastAPI, Depends
from pydantic import BaseModel
from sqlalchemy.orm import Session

from database import Base, engine, get_db
from intents import Intent
from ai import classify_intent
import crud
from models import ClassSchedule

Base.metadata.create_all(bind=engine)

app = FastAPI(title="CrossFit Box Chatbot - Postgres")

class ChatRequest(BaseModel):
    user_id: str
    message: str

class PlanCreateRequest(BaseModel):
    name: str
    description: str | None = None
    price: float
    active: bool = True

class ScheduleCreateRequest(BaseModel):
    weekday: str
    start_time: str # "18:00"
    end_time: str | None = None
    modality: str | None = None
    coach: str | None = None
    active: bool = True

def format_plans(plans) -> str:
    if not plans:
        return "Ainda não temos planos cadastrados."
    lines = ["💳 Planos disponíveis:"]
    for p in plans:
        desc = f" — {p.description}" if p.description else ""
        lines.append(f"- {p.name}: R$ {float(p.price):.2f}/mês{desc}")
    return "\n".join(lines)

def format_schedules(schedules) -> str:
    if not schedules:
        return "Ainda não temos horários cadastrados."
    lines = ["⏰ Horários disponíveis:"]
    for s in schedules:
        end = f" às {s.end_time.strftime('%H:%M')}" if s.end_time else ""
        mod = f" ({s.modality})" if s.modality else ""
        coach = f" - Coach {s.coach}" if s.coach else ""
        lines.append(f"- {s.weekday.value.upper()} {s.start_time.strftime('%H:%M')}{end}{mod}{coach}")
    return "\n".join(lines)

def handle_intent(intent: Intent, db: Session):
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
        return {"reply": "🏋️ Somos um Box focado em performance, saúde e comunidade. Quer ver horários ou planos?"}

    if intent == Intent.PLANS:
        plans = crud.list_plans(db)
        return {"reply": format_plans(plans)}

    if intent == Intent.SCHEDULE:
        schedules = crud.list_schedules(db)
        return {"reply": format_schedules(schedules)}

    if intent == Intent.BOOK_CLASS:
        # reserva fica pra próximo fluxo como você pediu
        return {"reply": "✅ Legal! Em breve vamos ativar reservas. Por enquanto posso te mostrar horários e planos 🙂"}

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

# ===== ADMIN API (simples) =====
@app.post("/admin/plans")
def create_plan(req: PlanCreateRequest, db: Session = Depends(get_db)):
    plan = crud.create_plan(db, req.name, req.description, req.price, req.active)
    return {"id": plan.id}

@app.get("/admin/plans")
def get_plans(db: Session = Depends(get_db)):
    plans = crud.list_plans(db)
    return plans

@app.post("/admin/schedules")
def create_schedule(req: ScheduleCreateRequest, db: Session = Depends(get_db)):
    # weekday precisa bater com enum: mon/tue/wed/thu/fri/sat/sun
    from datetime import datetime
    start_t = datetime.strptime(req.start_time, "%H:%M").time()
    end_t = datetime.strptime(req.end_time, "%H:%M").time() if req.end_time else None

    schedule = ClassSchedule(
        weekday=req.weekday,
        start_time=start_t,
        end_time=end_t,
        modality=req.modality,
        coach=req.coach,
        active=req.active,
    )
    created = crud.create_schedule(db, schedule)
    return {"id": created.id}

@app.get("/admin/schedules")
def get_schedules(db: Session = Depends(get_db)):
    return crud.list_schedules(db)

# ===== CHAT =====
@app.post("/chat")
def chat(req: ChatRequest, db: Session = Depends(get_db)):
    intent = classify_intent(req.message)
    return {"intent": intent, **handle_intent(intent, db)}
