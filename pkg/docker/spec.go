package docker

type AuthEncoder interface {
	// Encode computes the base64 encoded docker config JSON string used to
	// template kubernetes pull secrets. The actual data structure of the
	// encoded result looks like the example below.
	//
	//     {
	//         "auths": {
	//             "<registry>": {
	//                 "username": "<username>",
	//                 "password": "<password>",
	//                 "auth": "<auth>"
	//             }
	//         }
	//     }
	//
	Encode() (string, error)
}
