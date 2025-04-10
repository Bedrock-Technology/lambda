var { crypto, net, vars, db, utils, slog } = LambdaHelper

var tableName = vars.babylon_info_table

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var rpcMap = {
    '1': 'https://eth.llamarpc.com',
    '10': 'https://optimism.drpc.org',
    '56': 'https://binance.llamarpc.com',
    '196': 'https://rpc.xlayer.tech',
    '7000': 'https://zeta-chain.drpc.org',
    '34443': 'https://mainnet.mode.network',
    '42161': 'https://arbitrum.drpc.org',
    '59144': 'https://linea-rpc.publicnode.com',
    '60808': 'https://rpc.gobob.xyz',
    '534352': 'https://rpc.scroll.io',
    '169': 'https://pacific-rpc.manta.network/http',
    '810180': 'https://rpc.zklink.io',
    '223': 'https://rpc.bsquared.network',
    '4200': 'https://rpc.merlinchain.io',
    '5000': 'https://mantle-rpc.publicnode.com',
    '200901': 'https://rpc.bitlayer.org',
    '4689': 'https://babel-api.mainnet.iotex.io',
    '167000': 'https://rpc.mainnet.taiko.xyz',
    '80094': 'https://rpc.berachain.com',
    '146': 'https://sonic-rpc.publicnode.com',
    '43111': 'https://rpc.hemi.network/rpc',
}

var bodyObj = JSON.parse(vars.req.body)

slog.debug('[Bind]', 'req_body', vars.req.body)

var standardAddr = utils.hex_to_address(bodyObj.address)
require(utils.strings_equal_fold(standardAddr, bodyObj.address), 'invalid addr')

var [prefix, standardBabyAddr] = utils.bech32_address(bodyObj.babylon_address)
// require(bodyObj.babylon_address.startsWith('bbn1'), 'invalid babylon address')
require(prefix === 'bbn' && standardBabyAddr === bodyObj.babylon_address, 'invalid babylon address')

var data = db.select(`select address, point::numeric(20,8) as points, retain from public.baby_airdrop where address='${standardAddr}' limit 1`)
require(data.length > 0 && data[0].retain == 0 && data[0].points > 0, 'not eligible')

var typedData = {
    primaryType: 'Binding',
    types: {
        EIP712Domain: [
            { name: 'version', type: 'string' },
            { name: 'chainId', type: 'uint256' },
            { name: 'verifyingContract', type: 'address' },
        ],
        Binding: [
            { name: 'Babylon Address', type: 'string' },
        ],
    },
    domain: {
        version: '1',
        chainId: bodyObj.chain_id,
        verifyingContract: '0x004e9c3ef86bc1ca1f0bb5c7662861ee93350568',
    },
    message: {
        'Babylon Address': bodyObj.babylon_address,
    },
}

var hash = utils.hash_typed_data(JSON.stringify(typedData))
var addr = ''
if (bodyObj.signature.length == 65 + 2) { // signature of EOA or contract with one signer
    addr = crypto.ecrecover(hash, bodyObj.signature)
}

slog.debug('[Bind]', 'hash', hash, 'addr', addr)

//require(addr == standardAddr, 'invalid signature')
if (addr != standardAddr) {
    // eip 1271
    var rpcRaw = rpcMap[bodyObj.chain_id]

    var rpcReq = {
        jsonrpc: '2.0',
        id: 1,
        method: 'eth_call',
        params: [
            {
                to: standardAddr,
                // FIXME: max length of signature is 0xffffffff, it will be 66076419 signers
                data: `0x1626ba7e${hash.replace("0x", "")}000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000${((bodyObj.signature.length - 2) / 2).toString(16).padStart(8, '0')}${bodyObj.signature.replace("0x", "")}`,
            },
            'latest',
        ],
    }

    slog.debug('[Bind] before fetch', 'chain_id', bodyObj.chain_id, 'rpcRaw', rpcRaw, 'req', JSON.stringify(rpcReq))

    var fetchResp = net.fetch(rpcRaw, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(rpcReq),
    })

    slog.debug('[Bind] fetch result', 'body', fetchResp.body)

    var rpcRes = JSON.parse(fetchResp.body)
    require(rpcRes.result === '0x1626ba7e00000000000000000000000000000000000000000000000000000000', 'invalid signature')
}

db.insert(`${tableName}`, {
    created_at: (new Date()).toJSON(),
    updated_at: (new Date()).toJSON(),
    address: standardAddr,
    babylon_address: bodyObj.babylon_address,
    signature: bodyObj.signature,
})

var info = db.select(`select address, babylon_address, created_at from ${tableName} where address='${standardAddr}' limit 1`)

var resp = null
if (info.length > 0) {
    resp = info[0]
}

JSON.stringify({
    code: 200,
    data: resp,
})

