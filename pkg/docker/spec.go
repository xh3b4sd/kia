package docker

// AuthEncoder ...
type AuthEncoder interface {
	// Encode ...
	//
	//     {
	//         "auths": {
	//             "ghcr.io": {
	//                 "username": "...",
	//                 "password": "...",
	//                 "auth": "..."
	//             }
	//         }
	//     }
	//
	Encode() (string, error)
}
