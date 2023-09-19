package main

// Cook :
// Launch Pandora / Roll20
// Payloads :
// Pandora -> voiceChannelID (Redis)
// Roll20 -> Game ID	(HTTP)

// Cook
// (voiceChannelID, gameID) -> JobID

// Subscribe to JobID
// Return stream (grpc)

// Stop :

// Stop
// Pandora
// Roll20
// Cooking Server
// Encode-box
// Uploader
// DB
// Payloads :
// Pandora -> voiceChannelID (Redis)
// Roll20 -> Game ID	(HTTP)
// Cooking Server -> [Res Pandora] (Http)
// Encode-box -> [Res Cooking Server]  + [Res Roll20] (HTTP)
// Uploader -> [Res Encode-box] (HTTP)

// Emisison de progress
