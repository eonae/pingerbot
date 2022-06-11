-- With privacy enabled we won't get user ids, only usernames
-- So no need to store them in separate table
DROP TABLE IF EXISTS "members";

CREATE TABLE "members" (
	"username" varchar,
	"group_id" varchar,
	CONSTRAINT pk_groups PRIMARY KEY (username, group_id),
	CONSTRAINT fk_groups FOREIGN KEY (group_id)
	REFERENCES "groups"(id)
	ON DELETE CASCADE
);
