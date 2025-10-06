ALTER TABLE "chirps" ALTER COLUMN "body" SET DATA TYPE varchar(140);--> statement-breakpoint
ALTER TABLE "users" ADD COLUMN "hashed_password" varchar(256) DEFAULT 'unset' NOT NULL;