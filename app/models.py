import enum
from sqlalchemy import Column, Integer, String, Numeric, Boolean, Time, Enum

from app.database import Base

class Weekday(str, enum.Enum):
    mon = "mon"
    tue = "tue"
    wed = "wed"
    thu = "thu"
    fri = "fri"
    sat = "sat"
    sun = "sun"

class Plan(Base):
    __tablename__ = "plans"

    id = Column(Integer, primary_key=True, index=True)
    name = Column(String(80), nullable=False)
    description = Column(String(255), nullable=True)
    price = Column(Numeric(10, 2), nullable=False)
    active = Column(Boolean, default=True)

class ClassSchedule(Base):
    __tablename__ = "class_schedules"

    id = Column(Integer, primary_key=True, index=True)
    weekday = Column(Enum(Weekday), nullable=False, index=True)
    start_time = Column(Time, nullable=False)
    end_time = Column(Time, nullable=True)
    modality = Column(String(80), nullable=True)
    coach = Column(String(80), nullable=True)
    active = Column(Boolean, default=True)
