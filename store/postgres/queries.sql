-- name: create-messsages-table
CREATE TABLE messages (
	timestamp INTEGER NOT NULL,
	priority INTEGER,
	sender VARCHAR NOT NULL,
	ip VARCHAR NOT NULL,
	message VARCHAR NOT NULL,
	queued BOOLEAN NOT NULL,
	PRIMARY KEY (timestamp, sender)
);

-- name: insert-message
INSERT INTO messages (timestamp, priority, sender, ip, message, queued) VALUES
    ($1, $2, $3, $4, $5, FALSE);