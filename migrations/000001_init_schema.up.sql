-- Создание таблицы owners
CREATE TABLE owners (
                        uuid UUID PRIMARY KEY, -- Используем тип UUID для уникального идентификатора
                        name TEXT NOT NULL, -- Имя владельца
                        created_at TIMESTAMP DEFAULT now() -- Время создания записи
);

-- Добавление индекса для поля name в таблице owners
CREATE INDEX idx_owners_name ON owners(name);

-- Создание таблицы movies
CREATE TABLE movies (
                        uuid UUID PRIMARY KEY, -- Используем тип UUID для уникального идентификатора
                        owner_id UUID REFERENCES owners(uuid) ON DELETE CASCADE, -- Связь с таблицей owners
                        title TEXT NOT NULL, -- Название фильма
                        author TEXT NOT NULL, -- Автор фильма
                        description TEXT, -- Описание фильма
                        year INT CHECK (year > 0), -- Год выпуска фильма (проверка на положительное значение)
                        created_at TIMESTAMP DEFAULT now() -- Время создания записи
);