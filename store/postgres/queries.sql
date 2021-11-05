-- name: create-messsages-table
CREATE TABLE IF NOT EXISTS messages (
	timestamp BIGINT NOT NULL,
	priority INTEGER,
	sender VARCHAR NOT NULL,
	ip VARCHAR NOT NULL,
	message VARCHAR NOT NULL,
	queued BOOLEAN NOT NULL,
	PRIMARY KEY (timestamp, sender)
);

-- name: upsert-message
INSERT INTO messages (timestamp, priority, sender, ip, message, queued) VALUES
    ($1, $2, $3, $4, $5, $6)
ON CONFLICT(timestamp, sender) DO UPDATE
	SET queued=$6;

-- name: get-all-unqueued-messages
SELECT * FROM messages WHERE queued=FALSE;