-- name: create-messsages-table
CREATE TABLE messages (
	timestamp INTEGER,
	priority INTEGER,
	sender VARCHAR,
	ip VARCHAR,
	message VARCHAR,
	PRIMARY KEY (timestamp, sender)
);

-- name: insert-message
INSERT INTO messages (timestamp, priority, sender, ip, message) VALUES
    ($1, $2, $3, $4, $5);