# Standard library
from datetime import datetime
from typing import Dict, List

# Internal modules
from app import db


class Language(db.Model):  # type: ignore
    id: str = db.Column(db.String(2), primary_key=True)
    created_at: datetime = db.Column(db.DateTime, nullable=False, default=datetime.utcnow)

    def __repr__(self) -> str:
        return f"Language(id={self.id}, created_at={self.created_at})"


class TranslatedText(db.Model):  # type: ignore
    __table_args__ = (db.UniqueConstraint("key", "language", name="unique_key_language"),)
    id: int = db.Column(db.Integer, primary_key=True)
    key: str = db.Column(db.String(255), nullable=False)
    language: str = db.Column(db.String(2), db.ForeignKey("language.id"), nullable=True)
    value: str = db.Column(db.Text, nullable=False)
    created_at: datetime = db.Column(db.DateTime, nullable=False, default=datetime.utcnow)
    last_modified_at: datetime = db.Column(
        db.DateTime, nullable=False, default=datetime.utcnow, onupdate=datetime.utcnow
    )

    def __repr__(self) -> str:
        return (
            f"TranslatedText(id={self.id}, "
            f"key={self.key}, "
            f"language={self.language}, "
            f"value={self.value}, "
            f"created_at={self.created_at}, "
            f"last_modified_at={self.last_modified_at})"
        )


class TextGroup(db.Model):  # type: ignore
    __table_args__ = (
        db.UniqueConstraint("text_key", "group_id", name="unique_text_key_group_id"),
    )
    id: int = db.Column(db.Integer, primary_key=True)
    text_key: int = db.Column(db.String(255), nullable=False)
    group_id: int = db.Column(db.String(255), db.ForeignKey("group.id"), nullable=False)
    created_at: datetime = db.Column(db.DateTime, nullable=False, default=datetime.utcnow)

    def __repr__(self) -> str:
        return (
            f"TextGroup(id={self.id}, "
            f"text_key={self.text_key}, "
            f"group_id={self.group_id} ,"
            f"created_at={self.created_at})"
        )


class Group(db.Model):  # type: ignore
    id: str = db.Column(db.String(255), primary_key=True)
    created_at: datetime = db.Column(db.DateTime, nullable=False, default=datetime.utcnow)
    texts: List[TextGroup] = db.relationship("TextGroup", backref="group", lazy=True)

    def __repr__(self) -> str:
        return f"Group(id={self.id} created_at={self.created_at})"
