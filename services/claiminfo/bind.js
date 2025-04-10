var { crypto, vars, db, utils } = LambdaHelper

var tableName = vars.claim_info_table

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

require(false, 'disabled')

var bodyObj = JSON.parse(vars.req.body)

var standardAddr = utils.hex_to_address(bodyObj.address)
require(utils.strings_equal_fold(standardAddr, bodyObj.address), 'invalid addr')

var standardDepositAddr = utils.hex_to_address(bodyObj.deposit_address)
require(utils.strings_equal_fold(standardDepositAddr, bodyObj.deposit_address), 'invalid addr')

require(standardAddr != standardDepositAddr, 'deposit address should not equals to signer address')

require(['Gate', 'Bitget', 'Bybit'].includes(bodyObj.cex_type), 'invalid cex_type')

var data = db.select(`select address, amount::numeric(20,8) as br, (boosted > 0) as boosted from public.br_airdrop where address='${standardAddr}' limit 1`)
require(data.length > 0 && data[0].br > 0, 'not eligible')

var typedData = {
    primaryType: 'Deposit',
    types: {
        EIP712Domain: [
            { name: 'version', type: 'string' },
            { name: 'chainId', type: 'uint256' },
        ],
        Deposit: [
            { name: 'Exchange Name', type: 'string' },
            { name: 'Exchange UID', type: 'string' },
            { name: 'Deposit Address', type: 'string' },
        ],
    },
    domain: {
        version: '1',
        chainId: bodyObj.chain_id,
    },
    message: {
        'Deposit Address': bodyObj.deposit_address,
        'Exchange Name': bodyObj.cex_type,
        'Exchange UID': bodyObj.cex_uid,
    },
}

var hash = utils.hash_typed_data(JSON.stringify(typedData))
var addr = crypto.ecrecover(hash, bodyObj.signature)
require(addr == standardAddr, 'invalid signature')

db.insert(`${tableName}`, {
    created_at: (new Date()).toJSON(),
    updated_at: (new Date()).toJSON(),
    address: standardAddr,
    cex_type: bodyObj.cex_type,
    cex_uid: bodyObj.cex_uid,
    deposit_address: bodyObj.deposit_address,
    signature: bodyObj.signature,
})

var data = db.select(`select * from ${tableName} where address='${addr}' limit 1`)

var resp = null
if (data.length > 0) {
    resp = data[0]
}

JSON.stringify({
    code: 200,
    data: resp,
})

