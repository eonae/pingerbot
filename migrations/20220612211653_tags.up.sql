ALTER TABLE members ADD tag varchar NOT NULL;

ALTER TABLE members DROP CONSTRAINT pk_groups;

ALTER TABLE members ADD CONSTRAINT pk_members PRIMARY KEY (group_id, tag, username);
