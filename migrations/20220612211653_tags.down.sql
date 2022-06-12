ALTER TABLE members DROP CONSTRAINT pk_members

ALTER TABLE members ADD CONSTRAINT pk_groups PRIMARY KEY (username, group_id)

ALTER TABLE "members" DROP COLUMN tag
