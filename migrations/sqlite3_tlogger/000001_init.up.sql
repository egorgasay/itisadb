CREATE TABLE transactions (
     event_type int,
     key text,
     value text
);
CREATE TABLE info (
    server_number int
);
INSERT INTO info (server_number) VALUES (0);