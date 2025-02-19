erDiagram
    POST_STATISTICS {
        int id PK "Уникальный идентификатор статистики"
        int post_id FK "Связь с постом"
        int views_count "Количество просмотров"
        int likes_count "Количество лайков"
        int comments_count "Количество комментариев"
        datetime updated_at "Дата последнего обновления"
    }

    USER_ACTIVITY {
        int id PK "Уникальный идентификатор записи"
        int user_id FK "Связь с пользователем"
        int post_id FK "Связь с постом"
        string action_type "Тип действия (просмотр, лайк, комментарий)"
        datetime action_time "Время действия"
    }

    DAILY_STATISTICS {
        int id PK "Уникальный идентификатор"
        date stat_date "Дата"
        int total_views "Общее количество просмотров"
        int total_likes "Общее количество лайков"
        int total_comments "Общее количество комментариев"
    }

    POST ||--o{ POST_STATISTICS : "имеет"
    USER ||--o{ USER_ACTIVITY : "порождает"
    POST ||--o{ USER_ACTIVITY : "вовлекает"
    DAILY_STATISTICS }|..|{ POST_STATISTICS : "агрегирует"