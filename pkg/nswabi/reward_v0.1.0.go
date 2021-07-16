package nswabi

const RewardChequeBookABIv0_1_0 = `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_issuer",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "_erc20Addr",
				"type": "address"
			}
		],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "address",
				"name": "previousOwner",
				"type": "address"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "OwnershipTransferred",
		"type": "event"
	},
	{
		"inputs": [],
		"name": "CHEQUE_TYPEHASH",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "CheckBookIdRecord",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "EIP712DOMAIN_TYPEHASH",
		"outputs": [
			{
				"internalType": "bytes32",
				"name": "",
				"type": "bytes32"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "recipient",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "cumulativePayout",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "id",
				"type": "uint256"
			},
			{
				"internalType": "bytes",
				"name": "issuerSig",
				"type": "bytes"
			}
		],
		"name": "cashChequeBeneficiary",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "addr",
				"type": "address"
			}
		],
		"name": "configure",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "deposit",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getBalance",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "renounceOwnership",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "transferOwnership",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "withDraw",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

const RewardChequeBookBinv0_1_0 = "608060405234801561001057600080fd5b506040516112e63803806112e683398101604081905261002f916100c1565b60006100426001600160e01b036100bd16565b600080546001600160a01b0319166001600160a01b0383169081178255604051929350917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a350600280546001600160a01b039384166001600160a01b03199182161790915560018054929093169116179055610112565b3390565b600080604083850312156100d3578182fd5b82516100de816100fa565b60208401519092506100ef816100fa565b809150509250929050565b6001600160a01b038116811461010f57600080fd5b50565b6111c5806101216000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c8063715018a611610071578063715018a61461011157806375cb2672146101195780638da5cb5b1461012c578063b6b55f2514610141578063c49f91d314610154578063f2fde38b1461015c576100a9565b80630fdb1c10146100ae57806312065fe0146100b857806315c3343f146100d657806345414159146100de57806358be780e146100fe575b600080fd5b6100b661016f565b005b6100c0610243565b6040516100cd9190610deb565b60405180910390f35b6100c06102c9565b6100f16100ec366004610b79565b6102e0565b6040516100cd9190610de0565b6100b661010c366004610ba3565b610300565b6100b6610315565b6100b6610127366004610b57565b610394565b6101346103eb565b6040516100cd9190610d8f565b6100b661014f366004610c6f565b6103fa565b6100c061041b565b6100b661016a366004610b57565b610427565b6002546001600160a01b031633146101a25760405162461bcd60e51b815260040161019990610ffc565b60405180910390fd5b6001546040516370a0823160e01b81526102419133916001600160a01b03909116906370a08231906101d8903090600401610d8f565b60206040518083038186803b1580156101f057600080fd5b505afa158015610204573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102289190610c87565b6001546001600160a01b0316919063ffffffff6104dd16565b565b6001546040516370a0823160e01b81526000916001600160a01b0316906370a0823190610274903090600401610d8f565b60206040518083038186803b15801561028c57600080fd5b505afa1580156102a0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c49190610c87565b905090565b6040516102d590610d25565b604051809103902081565b600460209081526000928352604080842090915290825290205460ff1681565b61030f33858585600086610538565b50505050565b61031d610662565b6000546001600160a01b0390811691161461034a5760405162461bcd60e51b815260040161019990610fc7565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b61039c610662565b6000546001600160a01b039081169116146103c95760405162461bcd60e51b815260040161019990610fc7565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031690565b600154610418906001600160a01b031633308463ffffffff61066616565b50565b6040516102d590610cd6565b61042f610662565b6000546001600160a01b0390811691161461045c5760405162461bcd60e51b815260040161019990610fc7565b6001600160a01b0381166104825760405162461bcd60e51b815260040161019990610efd565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6105338363a9059cbb60e01b84846040516024016104fc929190610dc7565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152610687565b505050565b6002546001600160a01b0316331461058c5761055f61055930888787610716565b826107a4565b6002546001600160a01b0390811691161461058c5760405162461bcd60e51b815260040161019990611058565b6003546002546001600160a01b039081169116141561060e576001546040516340c10f1960e01b81526001600160a01b03909116906340c10f19906105d79088908890600401610dc7565b600060405180830381600087803b1580156105f157600080fd5b505af1158015610605573d6000803e3d6000fd5b5050505061062b565b60015461062b906001600160a01b0316868663ffffffff6104dd16565b50506001600160a01b03909316600090815260046020908152604080832095835294905292909220805460ff191660011790555050565b3390565b61030f846323b872dd60e01b8585856040516024016104fc93929190610da3565b60606106dc826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166107f59092919063ffffffff16565b80519091501561053357808060200190518101906106fa9190610c4f565b6105335760405162461bcd60e51b81526004016101999061109c565b6001600160a01b038316600090815260046020908152604080832084845290915281205460ff161561075a5760405162461bcd60e51b8152600401610199906110e6565b60405161076690610d25565b604051908190038120610783918790879087908790602001610df4565b6040516020818303038152906040528051906020012090505b949350505050565b6000806107b76107b2610804565b61085e565b846040516020016107c9929190610cbb565b6040516020818303038152906040528051906020012090506107eb81846108c1565b9150505b92915050565b606061079c8484600085610935565b61080c610b1f565b506040805160a081018252600a6060820190815269436865717565626f6f6b60b01b608083015281528151808301835260038152620312e360ec1b602082810191909152820152469181019190915290565b600060405161086c90610cd6565b604051809103902082600001518051906020012083602001518051906020012084604001516040516020016108a49493929190610e23565b604051602081830303815290604052805190602001209050919050565b60008151604114156108f55760208201516040830151606084015160001a6108eb868285856109f9565b93505050506107ef565b81516040141561091d5760208201516040830151610914858383610aef565b925050506107ef565b60405162461bcd60e51b815260040161019990610ec6565b606061094085610b19565b61095c5760405162461bcd60e51b815260040161019990611021565b60006060866001600160a01b031685876040516109799190610c9f565b60006040518083038185875af1925050503d80600081146109b6576040519150601f19603f3d011682016040523d82523d6000602084013e6109bb565b606091505b509150915081156109cf57915061079c9050565b8051156109df5780518082602001fd5b8360405162461bcd60e51b81526004016101999190610e5c565b60007f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0821115610a3b5760405162461bcd60e51b815260040161019990610f43565b8360ff16601b1480610a5057508360ff16601c145b610a6c5760405162461bcd60e51b815260040161019990610f85565b600060018686868660405160008152602001604052604051610a919493929190610e3e565b6020604051602081039080840390855afa158015610ab3573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610ae65760405162461bcd60e51b815260040161019990610e8f565b95945050505050565b60006001600160ff1b03821660ff83901c601b01610b0f868287856109f9565b9695505050505050565b3b151590565b60405180606001604052806060815260200160608152602001600081525090565b80356001600160a01b03811681146107ef57600080fd5b600060208284031215610b68578081fd5b610b728383610b40565b9392505050565b60008060408385031215610b8b578081fd5b610b958484610b40565b946020939093013593505050565b60008060008060808587031215610bb8578182fd5b84356001600160a01b0381168114610bce578283fd5b93506020850135925060408501359150606085013567ffffffffffffffff811115610bf7578182fd5b80860187601f820112610c08578283fd5b80359150610c1d610c1883611133565b61110c565b828152886020848401011115610c31578384fd5b610c42836020830160208501611157565b9598949750929550505050565b600060208284031215610c60578081fd5b81518015158114610b72578182fd5b600060208284031215610c80578081fd5b5035919050565b600060208284031215610c98578081fd5b5051919050565b60008251610cb1818460208701611163565b9190910192915050565b61190160f01b81526002810192909252602282015260420190565b7f454950373132446f6d61696e28737472696e67206e616d652c737472696e672081527f76657273696f6e2c75696e7432353620636861696e4964290000000000000000602082015260380190565b7f436865717565286164647265737320636865717565626f6f6b2c61646472657381527f732062656e65666963696172792c75696e743235362063756d756c61746976656020820152715061796f75742c75696e743235362069642960701b604082015260520190565b6001600160a01b0391909116815260200190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b901515815260200190565b90815260200190565b9485526001600160a01b0393841660208601529190921660408401526060830191909152608082015260a00190565b93845260208401929092526040830152606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6000602082528251806020840152610e7b816040850160208701611163565b601f01601f19169190910160400192915050565b60208082526018908201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604082015260600190565b6020808252601f908201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604082015260600190565b60208082526026908201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160408201526564647265737360d01b606082015260800190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604082015261756560f01b606082015260800190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604082015261756560f01b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252600b908201526a27b7363c9024a9a9aaa2a960a91b604082015260600190565b6020808252601d908201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604082015260600190565b60208082526024908201527f53696d706c65537761703a20696e76616c696420697373756572207369676e616040820152637475726560e01b606082015260800190565b6020808252602a908201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6040820152691bdd081cdd58d8d9595960b21b606082015260800190565b6020808252600c908201526b121055914810d2105491d15160a21b604082015260600190565b60405181810167ffffffffffffffff8111828210171561112b57600080fd5b604052919050565b600067ffffffffffffffff821115611149578081fd5b50601f01601f191660200190565b82818337506000910152565b60005b8381101561117e578181015183820152602001611166565b8381111561030f575050600091015256fea26469706673582212204e574175ecca6424720703224f5df3e34b56ffed3296bd9a06ae0c71bba4c5eb64736f6c63430006080033"
const RewardChequeBookDeployedBinv0_1_0 = "608060405234801561001057600080fd5b50600436106100a95760003560e01c8063715018a611610071578063715018a61461011157806375cb2672146101195780638da5cb5b1461012c578063b6b55f2514610141578063c49f91d314610154578063f2fde38b1461015c576100a9565b80630fdb1c10146100ae57806312065fe0146100b857806315c3343f146100d657806345414159146100de57806358be780e146100fe575b600080fd5b6100b661016f565b005b6100c0610243565b6040516100cd9190610deb565b60405180910390f35b6100c06102c9565b6100f16100ec366004610b79565b6102e0565b6040516100cd9190610de0565b6100b661010c366004610ba3565b610300565b6100b6610315565b6100b6610127366004610b57565b610394565b6101346103eb565b6040516100cd9190610d8f565b6100b661014f366004610c6f565b6103fa565b6100c061041b565b6100b661016a366004610b57565b610427565b6002546001600160a01b031633146101a25760405162461bcd60e51b815260040161019990610ffc565b60405180910390fd5b6001546040516370a0823160e01b81526102419133916001600160a01b03909116906370a08231906101d8903090600401610d8f565b60206040518083038186803b1580156101f057600080fd5b505afa158015610204573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102289190610c87565b6001546001600160a01b0316919063ffffffff6104dd16565b565b6001546040516370a0823160e01b81526000916001600160a01b0316906370a0823190610274903090600401610d8f565b60206040518083038186803b15801561028c57600080fd5b505afa1580156102a0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c49190610c87565b905090565b6040516102d590610d25565b604051809103902081565b600460209081526000928352604080842090915290825290205460ff1681565b61030f33858585600086610538565b50505050565b61031d610662565b6000546001600160a01b0390811691161461034a5760405162461bcd60e51b815260040161019990610fc7565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b61039c610662565b6000546001600160a01b039081169116146103c95760405162461bcd60e51b815260040161019990610fc7565b600380546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b031690565b600154610418906001600160a01b031633308463ffffffff61066616565b50565b6040516102d590610cd6565b61042f610662565b6000546001600160a01b0390811691161461045c5760405162461bcd60e51b815260040161019990610fc7565b6001600160a01b0381166104825760405162461bcd60e51b815260040161019990610efd565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6105338363a9059cbb60e01b84846040516024016104fc929190610dc7565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152610687565b505050565b6002546001600160a01b0316331461058c5761055f61055930888787610716565b826107a4565b6002546001600160a01b0390811691161461058c5760405162461bcd60e51b815260040161019990611058565b6003546002546001600160a01b039081169116141561060e576001546040516340c10f1960e01b81526001600160a01b03909116906340c10f19906105d79088908890600401610dc7565b600060405180830381600087803b1580156105f157600080fd5b505af1158015610605573d6000803e3d6000fd5b5050505061062b565b60015461062b906001600160a01b0316868663ffffffff6104dd16565b50506001600160a01b03909316600090815260046020908152604080832095835294905292909220805460ff191660011790555050565b3390565b61030f846323b872dd60e01b8585856040516024016104fc93929190610da3565b60606106dc826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166107f59092919063ffffffff16565b80519091501561053357808060200190518101906106fa9190610c4f565b6105335760405162461bcd60e51b81526004016101999061109c565b6001600160a01b038316600090815260046020908152604080832084845290915281205460ff161561075a5760405162461bcd60e51b8152600401610199906110e6565b60405161076690610d25565b604051908190038120610783918790879087908790602001610df4565b6040516020818303038152906040528051906020012090505b949350505050565b6000806107b76107b2610804565b61085e565b846040516020016107c9929190610cbb565b6040516020818303038152906040528051906020012090506107eb81846108c1565b9150505b92915050565b606061079c8484600085610935565b61080c610b1f565b506040805160a081018252600a6060820190815269436865717565626f6f6b60b01b608083015281528151808301835260038152620312e360ec1b602082810191909152820152469181019190915290565b600060405161086c90610cd6565b604051809103902082600001518051906020012083602001518051906020012084604001516040516020016108a49493929190610e23565b604051602081830303815290604052805190602001209050919050565b60008151604114156108f55760208201516040830151606084015160001a6108eb868285856109f9565b93505050506107ef565b81516040141561091d5760208201516040830151610914858383610aef565b925050506107ef565b60405162461bcd60e51b815260040161019990610ec6565b606061094085610b19565b61095c5760405162461bcd60e51b815260040161019990611021565b60006060866001600160a01b031685876040516109799190610c9f565b60006040518083038185875af1925050503d80600081146109b6576040519150601f19603f3d011682016040523d82523d6000602084013e6109bb565b606091505b509150915081156109cf57915061079c9050565b8051156109df5780518082602001fd5b8360405162461bcd60e51b81526004016101999190610e5c565b60007f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0821115610a3b5760405162461bcd60e51b815260040161019990610f43565b8360ff16601b1480610a5057508360ff16601c145b610a6c5760405162461bcd60e51b815260040161019990610f85565b600060018686868660405160008152602001604052604051610a919493929190610e3e565b6020604051602081039080840390855afa158015610ab3573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610ae65760405162461bcd60e51b815260040161019990610e8f565b95945050505050565b60006001600160ff1b03821660ff83901c601b01610b0f868287856109f9565b9695505050505050565b3b151590565b60405180606001604052806060815260200160608152602001600081525090565b80356001600160a01b03811681146107ef57600080fd5b600060208284031215610b68578081fd5b610b728383610b40565b9392505050565b60008060408385031215610b8b578081fd5b610b958484610b40565b946020939093013593505050565b60008060008060808587031215610bb8578182fd5b84356001600160a01b0381168114610bce578283fd5b93506020850135925060408501359150606085013567ffffffffffffffff811115610bf7578182fd5b80860187601f820112610c08578283fd5b80359150610c1d610c1883611133565b61110c565b828152886020848401011115610c31578384fd5b610c42836020830160208501611157565b9598949750929550505050565b600060208284031215610c60578081fd5b81518015158114610b72578182fd5b600060208284031215610c80578081fd5b5035919050565b600060208284031215610c98578081fd5b5051919050565b60008251610cb1818460208701611163565b9190910192915050565b61190160f01b81526002810192909252602282015260420190565b7f454950373132446f6d61696e28737472696e67206e616d652c737472696e672081527f76657273696f6e2c75696e7432353620636861696e4964290000000000000000602082015260380190565b7f436865717565286164647265737320636865717565626f6f6b2c61646472657381527f732062656e65666963696172792c75696e743235362063756d756c61746976656020820152715061796f75742c75696e743235362069642960701b604082015260520190565b6001600160a01b0391909116815260200190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b901515815260200190565b90815260200190565b9485526001600160a01b0393841660208601529190921660408401526060830191909152608082015260a00190565b93845260208401929092526040830152606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6000602082528251806020840152610e7b816040850160208701611163565b601f01601f19169190910160400192915050565b60208082526018908201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604082015260600190565b6020808252601f908201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604082015260600190565b60208082526026908201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160408201526564647265737360d01b606082015260800190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604082015261756560f01b606082015260800190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604082015261756560f01b606082015260800190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252600b908201526a27b7363c9024a9a9aaa2a960a91b604082015260600190565b6020808252601d908201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604082015260600190565b60208082526024908201527f53696d706c65537761703a20696e76616c696420697373756572207369676e616040820152637475726560e01b606082015260800190565b6020808252602a908201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6040820152691bdd081cdd58d8d9595960b21b606082015260800190565b6020808252600c908201526b121055914810d2105491d15160a21b604082015260600190565b60405181810167ffffffffffffffff8111828210171561112b57600080fd5b604052919050565b600067ffffffffffffffff821115611149578081fd5b50601f01601f191660200190565b82818337506000910152565b60005b8381101561117e578181015183820152602001611166565b8381111561030f575050600091015256fea26469706673582212204e574175ecca6424720703224f5df3e34b56ffed3296bd9a06ae0c71bba4c5eb64736f6c63430006080033"
