var { crypto, vars, db, utils } = LambdaHelper

var tableName = vars.claim_info_table

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var bodyObj = JSON.parse(vars.req.body)

var standardAddr = utils.hex_to_address(bodyObj.address)
require(utils.strings_equal_fold(standardAddr, bodyObj.address), 'invalid addr')

var typedData = {
    primaryType: 'Deposit',
    types: {
        EIP712Domain: [
            { name: 'version', type: 'string' },
        ],
        Deposit: [
            { name: 'Exchange Name', type: 'string' },
            { name: 'Exchange UID', type: 'string' },
            { name: 'Deposit Address', type: 'string' },
        ],
    },
    domain: {
        version: '1',
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
    data: resp,
})

