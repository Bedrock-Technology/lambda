// var rune_api_base = 'http://host.docker.internal:8580' // beta
var rune_api_base = 'http://host.docker.internal:8570' // prod

function firstOr(x, defaultValue) {
    if (x && x.length > 0) {
        return x[0]
    }
    return defaultValue
}

function getAmountByFunc(funcName, addr, start, end) {
    var payload = {
        func_name: funcName,
        params: JSON.stringify({
            user: addr,
            start_time: start,
            end_time: end,
        })
    }

    var resp = fetch(rune_api_base + '/dsn/execsql', {
        method: 'POST',
        body: JSON.stringify(payload)
    })

    var respObj = JSON.parse(resp.body)
    var mintedAmount = 0
    if (respObj.data && respObj.data[0]) {
        mintedAmount = Number(respObj.data[0].total_amount)
    }
    return mintedAmount
}

var addr = firstOr(req.query.address, '')
var start = firstOr(req.query.from, 0)
var end = firstOr(req.query.to, 0)
var amountLimit = firstOr(req.query.amount, 0)

var brBTCAmount = getAmountByFunc('FUNGetUserMintedBrBtcAmountALLCHAIN', addr, start, end)

JSON.stringify({
    result: brBTCAmount >= amountLimit
})
