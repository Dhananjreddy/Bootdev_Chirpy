import argon2 from "argon2";
import { JwtPayload } from "jsonwebtoken";
import jwt from "jsonwebtoken";
import { BadRequestError, UnauthorizedError } from "./api/errors.js";
import { Request } from "express";
import crypto from "crypto";

const TOKEN_ISSUER = "chirpy";
type payload = Pick<JwtPayload, "iss" | "sub" | "iat" | "exp">;

export async function hashPassword(password: string): Promise<string> {
    const hashedPassword = await argon2.hash(password);
    return Promise.resolve(hashedPassword);
}

export async function verifyPassword(hashedPassword: string, password: string): Promise<boolean> {
  try {
    return await argon2.verify(hashedPassword, password);
  } catch {
    return false;
  }
}

export function makeJWT(userID: string, expiresIn: number, secret: string): string{
    
    const payload: payload = {
        iss: TOKEN_ISSUER,
        sub: userID,
        iat: Math.floor(Date.now() / 1000),
        exp: Math.floor(Date.now() / 1000) + expiresIn,
    };  

    return jwt.sign(payload, secret);
}

export function verifyJWT(token: string, secret: string) {
    let decoded: payload;
    try {
        decoded = jwt.verify(token, secret) as payload;
    } catch (e) {
       throw new UnauthorizedError("Invalid token");
    }

    if (decoded.iss !== TOKEN_ISSUER) {
        throw new UnauthorizedError("Invalid token issuer");
    }
    if (!decoded.sub){
        throw new UnauthorizedError("Invalid token subject");
    }
    return decoded.sub;
}

export function getBearerToken(req: Request) {
  const authHeader = req.get("Authorization");
  if (!authHeader) {
    throw new UnauthorizedError("Malformed authorization header");
  }

  return extractBearerToken(authHeader);
}

export function extractBearerToken(header: string) {
  const splitAuth = header.split(" ");
  if (splitAuth.length < 2 || splitAuth[0] !== "Bearer") {
    throw new BadRequestError("Malformed authorization header");
  }
  return splitAuth[1];
}

export function makeRefreshToken(): string {
    return crypto.randomBytes(32).toString("hex");  
}

export function getAPIKey(req: Request) {
  const Authorization = req.get("Authorization");
  const apiKey = Authorization?.split(" ")[1];
  if (!apiKey) {
    throw new UnauthorizedError("API key is missing");
  }
  return apiKey;
}