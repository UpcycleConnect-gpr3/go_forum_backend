CREATE TABLE IF NOT EXISTS USER_MESSAGE (
    user_id CHAR(36) NOT NULL,
    message_id INT NOT NULL,
    PRIMARY KEY (user_id, message_id)
)
