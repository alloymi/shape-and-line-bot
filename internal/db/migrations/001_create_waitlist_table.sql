CREATE TABLE IF NOT EXISTS waitlist (
                                        id SERIAL PRIMARY KEY,
                                        chat_id BIGINT NOT NULL,
                                        fullname TEXT NOT NULL,
                                        email TEXT NOT NULL,
                                        course TEXT NOT NULL,
                                        created_at TIMESTAMP DEFAULT NOW()
    );
