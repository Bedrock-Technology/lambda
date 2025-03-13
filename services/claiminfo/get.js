var { crypto, db, vars, utils } = LambdaHelper

var tableName = vars.claim_info_table

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

function require(x, msg) {
    if (!x) {
        throw JSON.stringify({ msg: msg })
    }
}

var addr = firstOr(vars.req.query.addr, '')
var standardAddr = utils.hex_to_address(addr)
require(utils.strings_equal_fold(standardAddr, addr), 'invalid addr')

var data = db.select(`select address, cex_type, cex_uid, deposit_address, created_at from ${tableName} where address='${standardAddr}' limit 1`)

var resp = null
if (data.length > 0) {
    resp = data[0]
}

JSON.stringify({
    code: 200,
    data: resp,
})
