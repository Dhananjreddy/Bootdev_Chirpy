import { describe, it, expect, beforeAll } from "vitest";
import {
  hashPassword,
  verifyPassword,
  makeJWT,
  verifyJWT,
  extractBearerToken,
} from "./auth";
import { BadRequestError, UnauthorizedError } from "./api/errors.js";

describe("Password Hashing", () => {
  const password1 = "correctPassword123!";
  const password2 = "anotherPassword456!";
  let hash1: string;
  let hash2: string;

  beforeAll(async () => {
    hash1 = await hashPassword(password1);
    hash2 = await hashPassword(password2);
  });

  it("should return true for the correct password", async () => {
    const result = await verifyPassword(hash1, password1);
    expect(result).toBe(true);
  });

  it("should return false for an incorrect password", async () => {
    const result = await verifyPassword(hash1, "wrongPassword");
    expect(result).toBe(false);
  });

  it("should return false when password doesn't match a different hash", async () => {
    const result = await verifyPassword(hash2, password1);
    expect(result).toBe(false);
  });

  it("should return false for an empty password", async () => {
    const result = await verifyPassword(hash1, "");
    expect(result).toBe(false);
  });

  it("should return false for an invalid hash", async () => {
    const result = await verifyPassword("invalidhash", password1);
    expect(result).toBe(false);
  });
});

describe("JWT Functions", () => {
  const secret = "secret";
  const wrongSecret = "wrong_secret";
  const userID = "some-unique-user-id";
  let validToken: string;

  beforeAll(() => {
    validToken = makeJWT(userID, 3600, secret);
  });

  it("should validate a valid token", () => {
    const result = verifyJWT(validToken, secret);
    expect(result).toBe(userID);
  });

  it("should throw an error for an invalid token string", () => {
    expect(() => verifyJWT("invalid.token.string", secret)).toThrow(
      UnauthorizedError
    );
  });

  it("should throw an error when the token is signed with a wrong secret", () => {
    expect(() => verifyJWT(validToken, wrongSecret)).toThrow(
      UnauthorizedError
    );
  });
});

describe("extractBearerToken", () => {
  it("should extract the token from a valid header", () => {
    const token = "mySecretToken";
    const header = `Bearer ${token}`;
    expect(extractBearerToken(header)).toBe(token);
  });

  it("should extract the token even if there are extra parts", () => {
    const token = "mySecretToken";
    const header = `Bearer ${token} extra-data`;
    expect(extractBearerToken(header)).toBe(token);
  });

  it("should throw a BadRequestError if the header does not contain at least two parts", () => {
    const header = "Bearer";
    expect(() => extractBearerToken(header)).toThrow(BadRequestError);
  });

  it('should throw a BadRequestError if the header does not start with "Bearer"', () => {
    const header = "Basic mySecretToken";
    expect(() => extractBearerToken(header)).toThrow(BadRequestError);
  });

  it("should throw a BadRequestError if the header is an empty string", () => {
    const header = "";
    expect(() => extractBearerToken(header)).toThrow(BadRequestError);
  });
});
