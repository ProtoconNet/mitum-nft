### mitum-nft

*mitum-nft* is a nft contract model based on [mitum](https://github.com/ProtoconNet/mitum).

#### Features,

* user defined contract account state policy: collection.
* collection: user defined nft collection.
* *mongodb*: as mitum does, *mongodb* is the primary storage.
* ERC-721: ERC-721
* multiple collection policy for one contract account.

#### Installation

Before you build `mitum-nft`, make sure to run `docker run`.

```sh
$ git clone https://github.com/protoconNet/mitum-nft

$ cd mitum-nft

$ go build -ldflags="-X 'main.Version=v0.0.1-tutorial'" -o ./mitum-nft ./main.go
```

#### Run

```sh
$ ./mitum-nft node init <config file>

$ ./mitum-nft node run <config file>
```

[standalong.yml](standalone.yml) is a sample of `config file`.

#### Sample Operations

There are [sample operation files](sample/) in the repository.

Just try sending operations in the order listed below.

- [create-account/GEN_ACC2.json](sample/create-account/GEN_ACC2.json)
- [create-account/GEN_ACC4.json](sample/create-account/GEN_ACC4.json)
- [create-contract-account/GEN_ACC1.json](sample/create-contract-account/GEN_ACC1.json)
- [collection-register/GEN_AAA.json](sample/collection-register/GEN_AAA.json)
- [collection-register/GEN_BBB.json](sample/collection-register/GEN_BBB.json)
- [mint/GEN_AAA1.json](sample/mint/GEN_AAA1.json)
- [mint/GEN_AAA2.json](sample/mint/GEN_AAA2.json)
- [mint/GEN_AAA3.json](sample/mint/GEN_AAA3.json)
- [transfer/GEN_ACC2_AAA1.json](sample/transfer/GEN_ACC2_AAA1.json)
- [transfer/GEN_ACC2_AAA2.json](sample/transfer/GEN_ACC2_AAA2.json)
- [delegate/ACC2_AAA_ACC4.json](sample/delegate/ACC2_AAA_ACC4.json)
- [transfer/ACC4_ACC4_AAA1_delegated.json](sample/transfer/ACC4_ACC4_AAA1_delegated.json)
- [approve/ACC2_GEN_AAA2.json](sample/approve/ACC2_GEN_AAA2.json)
- [transfer/GEN_ACC4_AAA2_approved.json](sample/transfer/GEN_ACC4_AAA2_approved.json)
- [delegate/ACC2_AAA_ACC4_cancel.json](sample/delegate/ACC2_AAA_ACC4_cancel.json)
- [sign/ACC2_AAA2_creator.json](sample/sign/ACC2_AAA2_creator.json)
- [sign/ACC2_AAA2_copyrighter.json](sample/sign/ACC2_AAA2_copyrighter.json)
- [burn/ACC4_AAA1.json](sample/burn/ACC4_AAA1.json)
- [burn/ACC4_AAA2.json](sample/burn/ACC4_AAA2.json)
- [mint/ACC2_BBB1_white.json](sample/mint/ACC2_BBB1_white.json)
- [mint/ACC2_BBB2_white.json](sample/mint/ACC2_BBB2_white.json)
- [mint/ACC2_BBB3_white.json](sample/mint/ACC2_BBB3_white.json)
- [collection-policy-updater/GEN_AAA_update.json](sample/collection-policy-updater/GEN_AAA_update.json)
- [mint/ACC4_AAA4.json](sample/mint/ACC4_AAA4.json)

All results of operations must be `true`.