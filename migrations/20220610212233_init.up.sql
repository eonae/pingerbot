CREATE TABLE "groups" (
	"id" varchar PRIMARY KEY,
	"name" varchar
);

CREATE TABLE "members" (
	"id" varchar,
	"name" varchar,
	"group_id" varchar,
	CONSTRAINT pk_groups PRIMARY KEY (id, group_id),
	CONSTRAINT fk_groups FOREIGN KEY (group_id)
	REFERENCES "groups"(id) ON DELETE CASCADE
);
