from sqlalchemy import Column, Integer, String, Numeric, Boolean, Time, Enum
from database import Base
import enum

class Weekday(str, enum.Enum):
    monday = "mon"
    tuesday = "tue"
    wednesday = "wed"
    thursday = "thu"
    friday = "fri"
    saturday = "sat"
    sunday = "sun"

class Plan(Base):
    __tablename__ = "plans"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(80), nullable=False)
    description = Column(String(255), nullable=False)
    price = Column(Numeric(10, 2), nullable=False)
    duration_in_months = Column(Integer, nullable=False)
    is_active = Column(Boolean, default=True)

class ClassSchedule(Base):
    __tablename__ = "class_schedules"

    id = Column(Integer, primary_key=True, index=True)
    weekday = Column(Enum(Weekday), nullable=False, index=True)
    start_time = Column(Time, nullable=False)
    end_time = Column(Time, nullable=False)
    modality = Column(String(80), nullable=False)
    coach = Column(String(80), nullable=False)
    active = Column(Boolean, default=True)  
    
    