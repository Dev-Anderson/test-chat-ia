from enum import Enum

class Intent(str, Enum):
    GREETING = "greeting"
    ABOUT = "about"
    SCHEDULE = "schedule"
    PLANS = "plans"
    BOOK_CLASS = "book_class"
    UNKNOWN = "unknown"
