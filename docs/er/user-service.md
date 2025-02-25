erDiagram
    USER {
        int id PK "Уникальный идентификатор пользователя"
        string username "Имя пользователя"
        string email "Электронная почта"
        string password_hash "Хэш пароля"
        string role "Роль пользователя"
        datetime created_at "Дата регистрации"
    }

    SESSION {
        int id PK "Уникальный идентификатор сессии"
        int user_id FK "Связь с пользователем"
        string token "JWT токен"
        datetime created_at "Дата создания сессии"
        datetime expires_at "Дата истечения"
    }

    USER_PROFILE {
        int id PK "Уникальный идентификатор профиля"
        int user_id FK "Связь с пользователем"
        string first_name "Имя"
        string last_name "Фамилия"
        string bio "Биография"
        string avatar_url "Ссылка на аватарку"
    }

    USER ||--o{ SESSION : "имеет"
    USER ||--o{ USER_PROFILE : "имеет"
