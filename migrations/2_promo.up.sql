CREATE TABLE IF NOT EXISTS promos (
                                      id SERIAL PRIMARY KEY,
                                      company_id INT REFERENCES companies(id) ON DELETE CASCADE,
                                      description VARCHAR(300) NOT NULL,
                                      image_url VARCHAR(350) NOT NULL,
                                      mode VARCHAR(20) NOT NULL,
                                      promo_common VARCHAR(30), -- Общее промокод
                                      promo_unique TEXT[], -- Список уникальных промокодов
                                      target JSONB NOT NULL, -- Целевая аудитория, хранимая в формате JSONB
                                      max_count INTEGER NOT NULL CHECK (max_count > 0),
                                      active_from DATE NOT NULL,
                                      active_until DATE NOT NULL,
                                      like_count INTEGER,
                                      used_count INTEGER
);
