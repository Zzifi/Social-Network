erDiagram
    USERS {
        int id PK "Уникальный идентификатор пользователя"
        string username "Имя пользователя"
        string email "Электронная почта"
        string password_hash "Хэш пароля"
        string role "Роль пользователя"
        datetime created_at "Дата регистрации"
    }

    SESSIONS {
        int id PK "Уникальный идентификатор сессии"
        int user_id FK "Связь с пользователем"
        string token "JWT токен"
        datetime created_at "Дата создания сессии"
        datetime expires_at "Дата истечения"
    }

    USER_PROFILE {
        int user_id FK "Связь с пользователем"
        string phone_number "Номер телефона"
        string first_name "Имя"
        string last_name "Фамилия"
        datetime birthday "День рождения"
        string bio "Биография"
        string avatar_url "Ссылка на аватарку"
        datetime updated_at "Дата обновления профиля"
    }

    USER ||--o{ SESSION : "имеет"
    USER ||--o{ USER_PROFILE : "имеет"
