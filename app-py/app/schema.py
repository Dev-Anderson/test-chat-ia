from pydantic import BaseModel
from typing import Optional
from datetime import time
from models import Weekday

# Planos
class PlanBase(BaseModel):
    name: str
    description: Optional[str]= None
    price: float
    active: bool = True

class PlanCreate(PlanBase):
    pass

class PlanOut(PlanBase):
    id: int

    class Config:
        from_attributes = True

# Schedules 
class ScheduleBase(BaseModel):
    weekday: Weekday
    start_time: time
    end_time: Optional[time] = None
    modality: Optional[str] = None
    coach: Optional[str] = None
    active: bool = True

class ScheduleCreate(ScheduleBase):
    pass

class ScheduleOut(ScheduleBase):
    id: int

    class Config:
        from_attributes = True