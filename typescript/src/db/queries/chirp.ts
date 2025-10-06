import { db } from "../index.js";
import { NewChirp, chirps } from "../schema.js";
import { eq } from "drizzle-orm";


export async function createChirp(chirp: NewChirp) {
  const [result] = await db
    .insert(chirps)
    .values(chirp)
    .onConflictDoNothing()
    .returning();
  return result;
}

export async function getAllChirps() {
  return await db.select().from(chirps).orderBy(chirps.createdAt);
}

export async function getChirpsByChirpId(id: string) {
  return await db.select().from(chirps).where(eq(chirps.id, id));
}

export async function getChirpsByAuthorId(authorId: string) {
  return await db.select().from(chirps).where(eq(chirps.userId, authorId)).orderBy(chirps.createdAt);
}

export async function deleteChirpById(ChirpID: string) {
  return await db.delete(chirps).where(eq(chirps.id, ChirpID)).returning();
}
