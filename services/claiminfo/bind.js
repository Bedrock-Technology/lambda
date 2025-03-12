var { crypto, vars, db } = LambdaHelper

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var bodyObj = JSON.parse(vars.req.body)
require(utils.validate_address(bodyObj.address), 'invalid addr')

var addr = crypto.ecrecover(crypto.keccak256(bodyObj.msg), bodyObj.signature)
require(addr == bodyObj.address, 'invalid signature')

db.insert("claim_info", {
    created_at: (new Date()).toJSON(),
    updated_at: (new Date()).toJSON(),
    address: bodyObj.address,
    cex_type: bodyObj.cex_type,
    cex_uid: bodyObj.cex_uid,
    deposit_address: bodyObj.deposit_address,
    signature: bodyObj.signature,
})

var data = db.select(`select * from public.claim_info where address='${addr}'`)

JSON.stringify({
    data: data,
})
