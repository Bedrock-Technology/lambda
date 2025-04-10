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
var standardAddr = utils.hex_to_address(addr)
require(utils.strings_equal_fold(standardAddr, addr), 'invalid addr')

var data = db.select(`select address, point::numeric(20,8) as points, retain from public.baby_airdrop where address='${standardAddr}' limit 1`)

var resp = null
if (data.length > 0) {
    if (data[0].retain > 0) {
        data[0].points = 0
    }
    resp = {address: data[0].address, points: data[0].points}
}

JSON.stringify({
    code: 200,
    data: resp,
})
