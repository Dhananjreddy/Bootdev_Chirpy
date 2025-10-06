import { db } from "../index.js";
import { users } from "../schema.js";

export async function dropUsers() { 
    await db.delete(users);;   
    return;
}