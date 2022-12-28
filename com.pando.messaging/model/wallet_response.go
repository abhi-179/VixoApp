package model

type Createwallet struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Data    struct {
		Total   int64
		Data    []TransferDto `json:"data,omitempty"`
		Jsonrpc string        `json:"jsonrpc,omitempty"`
		ID      int64         `json:"id,omitempty"`
		Block   struct {
			Proposer  string `json:"Proposer,omitempty"`
			Timestamp int64  `json:"Timestamp,omitempty"`
		} `json:"block,omitempty"`
		Error struct {
			Code int64 `json:"code,omitempty"`
		} `json:"error,omitempty"`
		Result struct {
			Unlocked bool   `json:"unlocked,omitempty"`
			Address  string `json:"address,omitempty"`
			Sequence string `json:"sequence,omitempty"`
			Hash     string `json:"hash,omitempty"`
			Coins    struct {
				PandoWei string `json:"PandoWwei,omitempty"`
				PTXWei   string `json:"PtxWei,omitempty"`
			} `json:"coins,omitempty"`
			Reserved_funds            []string `json:"reserved_funds,omitempty"`
			Last_updated_block_height string   `json:"last_updated_block_height,omitempty"`
			Root                      string   `json:"root,omitempty"`
			Code                      string   `json:"code,omitempty"`
		} `json:"result,omitempty"`
	} `json:"data,omitempty"`
}

type GetBalanceReq struct {
	Jsonrpc string  `json:"jsonrpc,omitempty"`
	Method  string  `json:"method,omitempty"`
	ID      int     `json:"id,omitempty"`
	Param   []Param `json:"params,omitempty"`
}

type Param struct {
	Address  string `json:"address,omitempty"`
	ChainId  string `json:"chain_id,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	PandoWei string `json:"PandoWei,omitempty"`
	PTXWei   string `json:"PtxWei,omitempty"`
	Fee      string `json:"fee,omitempty"`
	Async    bool   `json:"async,omitempty"`
}

type TransferDto struct {
	Data      INPUT
	Hash      string `json:"hash,omitempty"`
	Status    string `json:"status,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}
type INPUT struct {
	Input  []Input  `json:"inputs,omitempty"`
	Output []Output `json:"outputs,omitempty"`
}
type Input struct {
	Address string `json:"address,omitempty"`
}
type Output struct {
	Address string `json:"address,omitempty"`
	Coins   struct {
		PandoWei string `json:"PandoWei,omitempty"`
		PTXWei   string `json:"PTXWei,omitempty"`
	} `json:"coins,omitempty"`
}

//  {
//     "success": true,
//     "message": "data",
//     "data": {
//         "jsonrpc": "2.0",
//         "id": 1,
//         "result": {
//             "hash": "0x924bc98a7dabdceda10f38b007ab1e3c22623c4c330f22a903367a76a64be1c9",
//             "block": {
//                 "ChainID": "pandonet",
//                 "Epoch": 385732,
//                 "Height": 375238,
//                 "Parent": "0x8459f2de4c2d1cc5061ebab7937cdd4562ff4faa16268f83329b4534f2af7deb",
//                 "HCC": {
//                     "Votes": [
//                         {
//                             "Block": "0x8459f2de4c2d1cc5061ebab7937cdd4562ff4faa16268f83329b4534f2af7deb",
//                             "Height": 375237,
//                             "Epoch": 385731,
//                             "ID": "0x29849cd55c86c6ccd3ba2a6d593daba4cf18a211",
//                             "Signature": "0x2c73614ca3516824ac99e6b61cb7b1cd4524dbfbc03b83e142cfc1207c884bb74b64697298c7ff48239e0b5c0322b1ab523905bead8e5f457aedeaeadc0e610b00"
//                         },
//                         {
//                             "Block": "0x8459f2de4c2d1cc5061ebab7937cdd4562ff4faa16268f83329b4534f2af7deb",
//                             "Height": 375237,
//                             "Epoch": 385732,
//                             "ID": "0x793bdada3144372d6e0856de0632c0a002bc6144",
//                             "Signature": "0xf29f4755f85d50438d23a0d678c72cc077ecdff80a443bb7613379ead369db0c7c0457e0415bb1dec9337f59da32334a3cece51e89c1ca9385f2e02ad2666cbf00"
//                         },
//                         {
//                             "Block": "0x8459f2de4c2d1cc5061ebab7937cdd4562ff4faa16268f83329b4534f2af7deb",
//                             "Height": 375237,
//                             "Epoch": 385731,
//                             "ID": "0x98fd878cd2267577ea6ac47bcb5ff4dd97d2f9e5",
//                             "Signature": "0xf742149fd58d4c1adc48606aaeb044b9561e8e7a7e620161ff19af569d72a4ff7aa5e524408e302e75ace76aa394512b8fa20302d8bd07f0002240b71600c87b01"
//                         },
//                         {
//                             "Block": "0x8459f2de4c2d1cc5061ebab7937cdd4562ff4faa16268f83329b4534f2af7deb",
//                             "Height": 375237,
//                             "Epoch": 385731,
//                             "ID": "0xc02fe0c1dd42451bcf717006b4c92883fc59858f",
//                             "Signature": "0xf8cf00c26b71292582094c59ce39d9d537d7e46290b3fd37e0c58b2630bfa5ef6c2854c2fb93a4a5f567adc449db6d239be731e83dd5e36b6572bec38f675b8e01"
//                         }
//                     ],
//                     "BlockHash": "0x8459f2de4c2d1cc5061ebab7937cdd4562ff4faa16268f83329b4534f2af7deb"
//                 },
//                 "GuardianVotes": null,
//                 "TxHash": "0x2f9cca03900689591e780c001e213436f9d51b67149ddf1c43f52581c887d0b0",
//                 "StateHash": "0xf1196cf4b4332684e897e82f70b58cd950221c437a06f19b8b8cc6e4061d7e16",
//                 "Timestamp": 1630668449,
//                 "Proposer": "0x29849cd55c86c6ccd3ba2a6d593daba4cf18a211",
//                 "Signature": "0xcd197b96599a57b94c2ac57ae0b98940a7366d4b5e77b9401607ff44c7da43bc4ebf6411ccb4698058ac133e1b3d59df4e78e15fcbdc75146fae182b9d78322001"
//             }
//         }
//     }
// }
