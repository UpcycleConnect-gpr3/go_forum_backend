CREATE TABLE IF NOT EXISTS CATEGORY_TALK (
    category_id INT NOT NULL,
    talk_id INT NOT NULL,
    PRIMARY KEY (category_id, talk_id)
);

CREATE TABLE IF NOT EXISTS USER_TALK (
    user_id CHAR(36) NOT NULL,
    talk_id INT NOT NULL,
    PRIMARY KEY (user_id, talk_id)
);

CREATE TABLE IF NOT EXISTS MESSAGE_TALK (
    message_id INT NOT NULL,
    talk_id INT NOT NULL,
    PRIMARY KEY (message_id, talk_id)
);

CREATE TABLE IF NOT EXISTS USER_MESSAGE (
    user_id CHAR(36) NOT NULL,
    message_id INT NOT NULL,
    PRIMARY KEY (user_id, message_id)
);

CREATE TABLE IF NOT EXISTS TALK_EVENT (
    talk_id INT NOT NULL,
    event_id INT NOT NULL,
    PRIMARY KEY (talk_id, event_id)
);

CREATE TABLE IF NOT EXISTS TALK_PROJECT (
    talk_id INT NOT NULL,
    project_id INT NOT NULL,
    PRIMARY KEY (talk_id, project_id)
);
