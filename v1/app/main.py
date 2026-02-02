import time
from datetime import datetime

from fastapi import FastAPI, Depends
from pydantic import BaseModel
from sqlalchemy import text
from sqlalchemy.exc import OperationalError
from sqlalchemy.orm import Session

from app.database import Base, engine, get_db
from app.intents import Intent
from app.ai import classify_intent
from app import crud
from app.models import ClassSchedule, Weekday
from fastapi.middleware.cors import CORSMiddleware


app = FastAPI(title="CrossFit Box Chatbot - Python + Postgres")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ============================================================
# STARTUP: aguarda Postgres ficar disponível e cria tabelas
# ============================================================
@app.on_event("startup")
def on_startup():
    retries = 20
    delay_seconds = 2

    for attempt in range(1, retries + 1):
        try:
            with engine.connect() as conn:
                conn.execute(text("SELECT 1"))
            break
        except OperationalError:
            print(f"[startup] aguardando Postgres... tentativa {attempt}/{retries}")
            time.sleep(delay_seconds)

    # quando o banco estiver acessível, cria tabelas
    Base.metadata.create_all(bind=engine)
    print("[startup] Postgres OK, tabelas garantidas ✅")


# ============================================================
# REQUEST MODELS
# ============================================================
class ChatRequest(BaseModel):
    user_id: str
    message: str


class PlanCreateRequest(BaseModel):
    name: str
    description: str | None = None
    price: float
    active: bool = True


class ScheduleCreateRequest(BaseModel):
    weekday: str               # mon/tue/wed/thu/fri/sat/sun
    start_time: str            # "18:00"
    end_time: str | None = None
    modality: str | None = None
    coach: str | None = None
    active: bool = True


# ============================================================
# HELPERS (formatação das respostas do chat)
# ============================================================
def format_plans(plans) -> str:
    if not plans:
        return "Ainda não temos planos cadastrados."
    lines = ["💳 *Planos disponíveis:*"]
    for p in plans:
        desc = f" — {p.description}" if p.description else ""
        lines.append(f"- {p.name}: R$ {float(p.price):.2f}/mês{desc}")
    return "\n".join(lines)


def format_schedules(schedules) -> str:
    if not schedules:
        return "Ainda não temos horários cadastrados."
    lines = ["⏰ *Horários disponíveis:*"]
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
        # reserva vai ser o próximo fluxo
        return {"reply": "✅ Perfeito! Em breve vamos ativar reservas. Por enquanto posso te mostrar horários e planos 🙂"}

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


# ============================================================
# HEALTHCHECK
# ============================================================
@app.get("/")
def health():
    return {"status": "ok"}


# ============================================================
# ADMIN - PLANS
# ============================================================
@app.post("/admin/plans")
def admin_create_plan(req: PlanCreateRequest, db: Session = Depends(get_db)):
    plan = crud.create_plan(db, req.name, req.description, req.price, req.active)
    return {"id": plan.id}


@app.get("/admin/plans")
def admin_get_plans(db: Session = Depends(get_db)):
    plans = crud.list_plans(db)
    return plans


# ============================================================
# ADMIN - SCHEDULES
# ============================================================
@app.post("/admin/schedules")
def admin_create_schedule(req: ScheduleCreateRequest, db: Session = Depends(get_db)):
    try:
        weekday = Weekday(req.weekday.strip().lower())
    except Exception:
        return {"error": "weekday inválido. Use: mon/tue/wed/thu/fri/sat/sun"}

    start_t = datetime.strptime(req.start_time, "%H:%M").time()
    end_t = datetime.strptime(req.end_time, "%H:%M").time() if req.end_time else None

    schedule = ClassSchedule(
        weekday=weekday,
        start_time=start_t,
        end_time=end_t,
        modality=req.modality,
        coach=req.coach,
        active=req.active,
    )

    created = crud.create_schedule(db, schedule)
    return {"id": created.id}


@app.get("/admin/schedules")
def admin_get_schedules(db: Session = Depends(get_db)):
    return crud.list_schedules(db)


# ============================================================
# CHAT
# ============================================================
@app.post("/chat")
def chat(req: ChatRequest, db: Session = Depends(get_db)):
    intent = classify_intent(req.message)
    return {"intent": intent, **handle_intent(intent, db)}
