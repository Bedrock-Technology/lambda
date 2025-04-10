var { crypto, vars, db, utils } = LambdaHelper

var tableName = vars.babylon_info_table

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var bodyObj = JSON.parse(vars.req.body)

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
var addr = crypto.ecrecover(hash, bodyObj.signature)
require(addr == standardAddr, 'invalid signature')

db.insert(`${tableName}`, {
    created_at: (new Date()).toJSON(),
    updated_at: (new Date()).toJSON(),
    address: standardAddr,
    babylon_address: bodyObj.babylon_address,
    signature: bodyObj.signature,
})

var info = db.select(`select address, babylon_address, created_at from ${tableName} where address='${addr}' limit 1`)

var resp = null
if (info.length > 0) {
    resp = info[0]
}

JSON.stringify({
    code: 200,
    data: resp,
})

