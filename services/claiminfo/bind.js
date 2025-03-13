var { crypto, vars, db, utils } = LambdaHelper

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var bodyObj = JSON.parse(vars.req.body)

var standardAddr = utils.hex_to_address(bodyObj.address)
require(utils.strings_equal_fold(standardAddr, bodyObj.address), 'invalid addr')

var addr = crypto.ecrecover(crypto.keccak256(bodyObj.msg), bodyObj.signature)
require(addr == standardAddr, 'invalid signature')

db.insert("claim_info", {
    created_at: (new Date()).toJSON(),
    updated_at: (new Date()).toJSON(),
    address: standardAddr,
    cex_type: bodyObj.cex_type,
    cex_uid: bodyObj.cex_uid,
    deposit_address: bodyObj.deposit_address,
    msg_raw: bodyObj.msg,
    signature: bodyObj.signature,
})

var data = db.select(`select * from public.claim_info where address='${addr}' limit 1`)

var resp = null
if (data.length > 0) {
    resp = data[0]
}

JSON.stringify({
    data: resp,
})

