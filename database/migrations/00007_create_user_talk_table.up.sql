CREATE TABLE IF NOT EXISTS USER_TALK (
    user_id CHAR(36) NOT NULL,
    talk_id INT NOT NULL,
    PRIMARY KEY (user_id, talk_id)
)
