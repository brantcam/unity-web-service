CREATE TABLE messages (
	timestamp INTEGER,
	priority INTEGER,
	sender VARCHAR,
	ip VARCHAR,
	message VARCHAR,
	PRIMARY KEY (timestamp, sender)
);

-- name: insert-message
INSERT INTO messages (ts, p, sdr, ip, msg) VALUES
    ($1, $2, $3, $4, $5);