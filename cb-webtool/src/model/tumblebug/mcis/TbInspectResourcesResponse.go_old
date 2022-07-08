package mcis

type TbInspectResourcesResponse struct {
	ResourcesOnCsp       []string                `json:"resourcesOnCsp"` // interface type으로 3가지 를 모두 받도록 되어있어 일단은 string을 받게 함. Test하여 보완할 것.
	ResourcesOnSpider    []ResourceOnCspOrSpider `json:"resourcesOnSpider"`
	ResourcesOnTumblebug []ResourceOnTumblebug   `json:"resourcesOnTumblebug"`
}

// Request
// {
// 	"connectionName": "m-aws-is-conn",
// 	"type": "vNet"
//   }

// Response
// {
//     "resourcesOnCsp": [
//         {
//             "cspNativeId": "vpc-062b15b9d1b67c4e0",
//             "id": "m-aws-is-vnet"
//         },
//         {
//             "cspNativeId": "vpc-02f2c205cb4e6f7ec",
//             "id": ""
//         },
//         {
//             "cspNativeId": "vpc-0071c41d0678975d8",
//             "id": ""
//         },
//         {
//             "cspNativeId": "vpc-05b33cfa62ce025fd",
//             "id": ""
//         },
//         {
//             "cspNativeId": "vpc-0472a630e5ea441ca",
//             "id": ""
//         }
//     ],
//     "resourcesOnSpider": [
//         {
//             "cspNativeId": "vpc-062b15b9d1b67c4e0",
//             "id": "m-aws-is-vnet"
//         }
//     ],
//     "resourcesOnTumblebug": [
//         {
//             "cspNativeId": "vpc-062b15b9d1b67c4e0",
//             "id": "m-aws-is-vnet",
//             "nsId": "noman",
//             "objectKey": "/ns/noman/resources/vNet/m-aws-is-vnet",
//             "type": "vNet"
//         }
//     ]
// }
