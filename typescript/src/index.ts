import express from "express";

import { handlerReadiness } from "./api/readiness.js";
import { errorHandler, middlewareLogResponses, middlewareMetricsInc } from "./api/middleware.js";
import { handlerCountHits } from "./admin/metrics.js";
import { handlerReset } from "./admin/reset.js";
import { handlerValidateChirp } from "./api/validate.js";
import { handlerCreateUser, handlerUpdateUser } from "./api/createUser.js";
import { handlerCreateChirp, handlerDeleteChirp } from "./api/createChirp.js";
import { handlerGetAllChirps, handlerGetChirpsByChirpId } from "./api/getChirps.js";
import { handlerLogin, handlerRefresh, handlerRevoke } from "./api/login.js";
import { handlerRedChirpy } from "./api/membership.js";

const app = express();
const PORT = 8080;

app.use("/app", middlewareMetricsInc, express.static("./src/app"));
app.use(middlewareLogResponses);
app.use(express.json());

app.get("/api/healthz", handlerReadiness);
app.get("/admin/metrics", handlerCountHits);
app.post("/admin/reset", handlerReset);
app.post("/api/validate_chirp", handlerValidateChirp);
app.post("/api/users", handlerCreateUser);
app.put("/api/users", handlerUpdateUser);
app.post("/api/chirps", handlerCreateChirp)
app.get("/api/chirps", handlerGetAllChirps);
app.get("/api/chirps/:chirpID", handlerGetChirpsByChirpId);
app.delete("/api/chirps/:chirpID", handlerDeleteChirp);
app.post("/api/login", handlerLogin);
app.post("/api/refresh", handlerRefresh);
app.post("/api/revoke", handlerRevoke);
app.post("/api/polka/webhooks", handlerRedChirpy)
app.listen(PORT, () => {
  console.log(`Server is running at http://localhost:${PORT}`);
});

app.use(errorHandler)