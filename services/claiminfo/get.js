var { crypto, db, vars, utils } = LambdaHelper

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
require(utils.validate_address(addr), 'invalid addr')

var data = db.select(`select * from public.claim_info where address='${addr}'`)

JSON.stringify({
    data: data,
})
