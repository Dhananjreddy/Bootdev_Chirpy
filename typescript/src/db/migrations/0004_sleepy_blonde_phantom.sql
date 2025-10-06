ALTER TABLE "refresh_tokens" ADD COLUMN "token" varchar(256) PRIMARY KEY NOT NULL;--> statement-breakpoint
ALTER TABLE "refresh_tokens" DROP COLUMN "id";