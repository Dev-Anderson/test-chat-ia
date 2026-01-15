from sqlalchemy.orm import Session
from models import Plan, ClassSchedule, Weekday

# Planos 
def list_plans(db: Session):
    return db.query(Plan).filter(Plan.active == True).order_by(Plan.price.asc()).all()

def create_plan(db: Session, name: str, description: str | None, price: float, active: bool = True):
    p = Plan(name=name, description=description, price=price, active=active)
    db.add(p)
    db.commit()
    db.refresh(p)
    return p

def list_schedules(db: Session):
    return (
        db.query(ClassSchedule)
        .filter(ClassSchedule.active == True)
        .order_by(ClassSchedule.weekday.asc(), ClassSchedule.start_time.asc())
        .all()
    )

def list_schedules_by_day(db: Session, weekday: Weekday):
    return (
        db.query(ClassSchedule)
        .filter(ClassSchedule.active == True, ClassSchedule.weekday == weekday)
        .order_by(ClassSchedule.start_time.asc())
        .all()
    )

def create_schedule(db: Session, schedule: ClassSchedule):
    db.add(schedule)
    db.commit()
    db.refresh(schedule)
    return schedule