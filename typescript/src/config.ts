import type { MigrationConfig } from "drizzle-orm/migrator";
import { env } from "process";

type Config = {
  api: APIConfig;
  db: DBConfig;
  jwt: JWTConfig;
};

type APIConfig = {
  fileServerHits: number;
  platform: string;
  port: number;
  polkaKey: string;
};

type DBConfig = {
  url: string;
  migrationConfig: MigrationConfig;
};
type JWTConfig = {
  defaultDuration: number;
  refreshDuration: number;
  secret: string;
  issuer: string;
};

process.loadEnvFile();

function envOrThrow(key: string) {
  const value = process.env[key];
  if (!value) {
    throw new Error(`Environment variable ${key} is not set`);
  }
  return value;
}

const migrationConfig: MigrationConfig = {
  migrationsFolder: "./src/db/migrations",
};

export const config: Config = {
  api: {
    fileServerHits: 0,
    port: Number(envOrThrow("PORT")),
    platform: envOrThrow("PLATFORM"),
    polkaKey: envOrThrow("POLKA_KEY"),
  },
  db: {
    url: envOrThrow("DB_URL"),
    migrationConfig: migrationConfig,
  },
  jwt: {
    defaultDuration: 60 * 60, 
    refreshDuration: 60 * 60 * 24 * 60 * 1000,
    secret: envOrThrow("JWT_SECRET"),
    issuer: "chirpy",
  },
};

