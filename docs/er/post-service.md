erDiagram
    POST {
        int id PK "Уникальный идентификатор поста"
        int user_id FK "Автор поста"
        string title "Заголовок поста"
        string content "Содержимое поста"
        datetime created_at "Дата создания"
        datetime updated_at "Дата обновления"
    }

    COMMENT {
        int id PK "Уникальный идентификатор комментария"
        int post_id FK "Связь с постом"
        int user_id FK "Автор комментария"
        int parent_id "Родительский комментарий (null, если нет)"
        string content "Содержимое комментария"
        datetime created_at "Дата создания"
    }

    LIKE {
        int id PK "Уникальный идентификатор лайка"
        int post_id FK "Связь с постом"
        int user_id FK "Пользователь, поставивший лайк"
        datetime created_at "Дата создания лайка"
    }

    USER ||--o{ COMMENT : "оставляет"
    USER ||--o{ LIKE : "ставит"
    POST ||--o{ LIKE : "получает"
